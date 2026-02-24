package db

import (
	"context"
	"fmt"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
)

// CreateNode inserts a new node.
func (db DB) CreateNode(ctx context.Context, name models.Optional[string], geom string) (models.Node, error) {
	point, err := stringToPoint(geom)
	if err != nil {
		return models.Node{}, fmt.Errorf("invalid point format: %w", err)
	}
	var n *string
	if !name.Set {
		n = nil
	}
	n = &name.Value
	arg := generated.CreateNodeParams{
		Name: n,
		Geom: point,
	}
	genNode, err := db.queries.CreateNode(ctx, arg)
	if err != nil {
		return models.Node{}, err
	}
	return convertGeneratedNodeToModel(genNode), nil
}

// GetNodeByID returns a node by its ID.
func (db DB) GetNodeByID(ctx context.Context, nodeID int) (models.Node, error) {
	genNode, err := db.queries.GetNode(ctx, int32(nodeID))
	if err != nil {
		return models.Node{}, err
	}
	return convertGeneratedNodeToModel(genNode), nil
}

// GetNodes returns a paginated list of nodes matching the filter.
func (db DB) GetNodes(ctx context.Context, limit, offset int, filter models.NodeFilter) ([]models.Node, int64, error) {
	arg := generated.GetNodesParams{
		Limit:       int32(limit),
		Offset:      int32(offset),
		NameFilter:  filter.Name.ToPtr(),
		CreatedFrom: optionalTimeToPgTimestamptz(filter.CreatedFrom),
		CreatedTo:   optionalTimeToPgTimestamptz(filter.CreatedTo),
		UpdatedFrom: optionalTimeToPgTimestamptz(filter.UpdatedFrom),
		UpdatedTo:   optionalTimeToPgTimestamptz(filter.UpdatedTo),
		SortBy:      filter.SortBy.ToPtr(),
		SortOrder:   filter.SortOrder.ToPtr(),
	}
	rows, err := db.queries.GetNodes(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var nodes []models.Node
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		nodes = append(nodes, convertGeneratedNodeRowToModel(row))
	}
	return nodes, totalCount, nil
}

// UpdateNode updates mutable fields of a node.
func (db DB) UpdateNode(ctx context.Context, nodeID int, name models.Optional[string], geom models.Optional[models.Point]) error {
	arg := generated.UpdateNodeParams{
		NodeID: int32(nodeID),
	}
	if !name.Set {
		arg.Name = &name.Value
	}
	if !geom.Set {
		arg.Geom.P.X = geom.Value.X
		arg.Geom.P.Y = geom.Value.Y
	}
	return db.queries.UpdateNode(ctx, arg)
}

// SoftDeleteNode marks a node as deleted.
func (db DB) SoftDeleteNode(ctx context.Context, nodeID int) error {
	return db.queries.SoftDeleteNode(ctx, int32(nodeID))
}

// HardDeleteNode permanently removes a node.
func (db DB) HardDeleteNode(ctx context.Context, nodeID int) error {
	return db.queries.HardDeleteNode(ctx, int32(nodeID))
}

// RestoreNode removes the soft‑delete mark.
func (db DB) RestoreNode(ctx context.Context, nodeID int) error {
	return db.queries.RestoreNode(ctx, int32(nodeID))
}

// BulkSoftDeleteNodes soft‑deletes multiple nodes.
func (db DB) BulkSoftDeleteNodes(ctx context.Context, nodeIDs []int) error {
	return db.queries.BulkSoftDeleteNodes(ctx, convertIntSliceToInt32(nodeIDs))
}

// BulkHardDeleteNodes permanently deletes multiple nodes.
func (db DB) BulkHardDeleteNodes(ctx context.Context, nodeIDs []int) error {
	return db.queries.BulkHardDeleteNodes(ctx, convertIntSliceToInt32(nodeIDs))
}

// conversion helpers
func convertGeneratedNodeToModel(n generated.Node) models.Node {
	var name string
	if n.Name != nil {
		name = *n.Name
	}
	return models.Node{
		NodeID: n.NodeID,
		Name:   name,
		Geom: models.Point{
			X: n.Geom.P.X,
			Y: n.Geom.P.Y,
		},
		CreatedAt: n.CreatedAt.Time,
		UpdatedAt: fromPgTimestamptz(n.UpdatedAt),
		DeletedAt: fromPgTimestamptz(n.DeletedAt),
	}
}

func convertGeneratedNodeRowToModel(row generated.GetNodesRow) models.Node {
	var name string
	if row.Name != nil {
		name = *row.Name
	}
	return models.Node{
		NodeID: row.NodeID,
		Name:   name,
		Geom: models.Point{
			X: row.Geom.P.X,
			Y: row.Geom.P.Y,
		},
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: fromPgTimestamptz(row.UpdatedAt),
		DeletedAt: fromPgTimestamptz(row.DeletedAt),
	}
}
