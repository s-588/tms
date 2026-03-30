package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/s-588/tms/internal/ui"
)

type CreateNodeArgs struct {
	Name    models.Optional[string]
	Address string
	Geom    models.Point
}

func (db DB) CreateNode(ctx context.Context, args CreateNodeArgs) (models.Node, error) {
	orbPt := orb.Point{args.Geom.X, args.Geom.Y}

	geomBytes, err := wkb.Marshal(orbPt)
	if err != nil {
		return models.Node{}, fmt.Errorf("wkb marshal: %w", err)
	}

	arg := generated.CreateNodeParams{
		Name: ToStringPtr(args.Name),
		Geom: geomBytes,
	}
	genNode, err := db.queries.CreateNode(ctx, arg)
	if err != nil {
		return models.Node{}, err
	}
	return convertGeneratedNodeToModel(genNode), nil
}

func (db DB) GetNodeByID(ctx context.Context, nodeID int32) (models.Node, error) {
	genNode, err := db.queries.GetNode(ctx, nodeID)
	if err != nil {
		return models.Node{}, err
	}
	return models.Node{
		NodeID: genNode.NodeID,
		Name:   fromStringPtr(genNode.Name),
		Geom: models.Point{
			X: toFloat64(genNode.X),
			Y: toFloat64(genNode.Y),
		},
		CreatedAt: fromPgTimestamptz(genNode.CreatedAt),
		UpdatedAt: fromPgTimestamptz(genNode.UpdatedAt),
		DeletedAt: fromPgTimestamptz(genNode.DeletedAt),
	}, nil
}

func (db DB) GetNodes(ctx context.Context, page int32, filter models.NodeFilter) ([]models.Node, int32, error) {
	arg := generated.GetNodesParams{
		Page:       page,
		NameFilter: ToStringPtr(filter.Name),
		SortBy:     filter.SortBy.Value,
		SortOrder:  filter.SortOrder.Value,
	}
	rows, err := db.queries.GetNodes(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var nodes []models.Node
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		nodes = append(nodes, convertGeneratedNodeRowToModel(row))
	}
	return nodes, totalPages, nil
}

type UpdateNodeArgs struct {
	NodeID  int32
	Name    models.Optional[string]
	Geom    models.Point
	Address string
}

func (db DB) UpdateNode(ctx context.Context, args UpdateNodeArgs) error {
	var n *string
	if args.Name.Set {
		n = &args.Name.Value
	}

	orbPt := orb.Point{args.Geom.X, args.Geom.Y}
	geomBytes, err := wkb.Marshal(orbPt)
	if err != nil {
		return fmt.Errorf("wkb marshal: %w", err)
	}

	arg := generated.UpdateNodeParams{
		NodeID:  args.NodeID,
		Name:    n,
		Address: args.Address,
		Geom:    geomBytes,
	}
	return db.queries.UpdateNode(ctx, arg)
}

func (db DB) SoftDeleteNode(ctx context.Context, nodeID int32) error {
	return db.queries.SoftDeleteNode(ctx, nodeID)
}

func (db DB) HardDeleteNode(ctx context.Context, nodeID int32) error {
	return db.queries.HardDeleteNode(ctx, nodeID)
}

func (db DB) RestoreNode(ctx context.Context, nodeID int32) error {
	return db.queries.RestoreNode(ctx, nodeID)
}

func (db DB) BulkSoftDeleteNodes(ctx context.Context, nodeIDs []int32) error {
	return db.queries.BulkSoftDeleteNodes(ctx, nodeIDs)
}

func (db DB) BulkHardDeleteNodes(ctx context.Context, nodeIDs []int32) error {
	return db.queries.BulkHardDeleteNodes(ctx, nodeIDs)
}

func convertGeneratedNodeToModel(n generated.Node) models.Node {
	var orbPt orb.Point
	geom, err := wkb.Unmarshal(n.Geom)
	if err != nil {
		slog.Error("can't unmarshal node from database", "node", n)
		orbPt = orb.Point{}
	}
	if point, ok := geom.(orb.Point); ok {
		orbPt = point
	} else {
		slog.Error("can't convert node geometry to point", "node", n)
	}
	return models.Node{
		NodeID: n.NodeID,
		Name:   fromStringPtr(n.Name),
		Geom: models.Point{
			Y: orbPt.Y(),
			X: orbPt.X(),
		},
		CreatedAt: fromPgTimestamptz(n.CreatedAt),
		UpdatedAt: fromPgTimestamptz(n.UpdatedAt),
		DeletedAt: fromPgTimestamptz(n.DeletedAt),
	}
}

func convertGeneratedNodeRowToModel(row generated.GetNodesRow) models.Node {
	return models.Node{
		NodeID: row.NodeID,
		Name:   fromStringPtr(row.Name),
		Geom: models.Point{
			X: row.X,
			Y: row.Y,
		},
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
		UpdatedAt: fromPgTimestamptz(row.UpdatedAt),
		DeletedAt: fromPgTimestamptz(row.DeletedAt),
	}
}

func (db DB) ListNodes(ctx context.Context) ([]ui.ListItem, error) {
	rows, err := db.queries.ListNodes(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ui.ListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, ui.ListItem{
			ID:   r.NodeID,
			Name: r.Address,
		})
	}
	return items, nil
}

func parseNodesError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "nodes_address_key" {
			return ErrDuplicatePrice
		}
		return fmt.Errorf("unhandled error: %w", err)
	}
	return fmt.Errorf("uknown error: %w", err)
}

func (db DB) CalculateDistance(ctx context.Context, nodeAID, nodeBID int32) (float64, error) {
	distance, err := db.queries.GetDistanceBetweenNodes(ctx, generated.GetDistanceBetweenNodesParams{
		NodeID:   nodeAID,
		NodeID_2: nodeBID,
	})
	if err != nil {
		return 0, parseNodesError(err)
	}
	return float64(distance) / 1000, nil
}
