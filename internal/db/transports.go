package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/s-588/tms/internal/ui"
)

type CreateTransportArgs struct {
	Model           string
	LicensePlate    string
	PayloadCapacity int32
	FuelConsumption int32
}

func (db DB) CreateTransport(ctx context.Context, args CreateTransportArgs) (models.Transport, error) {
	arg := generated.CreateTransportParams{
		Model:           args.Model,
		LicensePlate:    args.LicensePlate,
		PayloadCapacity: args.PayloadCapacity,
		FuelConsumption: args.FuelConsumption,
	}
	genTransport, err := db.queries.CreateTransport(ctx, arg)
	if err != nil {
		return models.Transport{}, parseTransportsError(err)
	}
	return convertGeneratedTransportToModel(genTransport), nil
}

func (db DB) GetTransportByID(ctx context.Context, transportID int32) (models.Transport, error) {
	genTransport, err := db.queries.GetTransport(ctx, transportID)
	if err != nil {
		return models.Transport{}, err
	}
	return convertGeneratedTransportToModel(genTransport), nil
}

func (db DB) GetTransports(ctx context.Context, page int32, filter models.TransportFilter) ([]models.Transport, int32, error) {
	arg := generated.GetTransportsParams{
		Page:               page,
		ModelFilter:        ToStringPtr(filter.Model),
		LicensePlateFilter: ToStringPtr(filter.LicensePlate),
		PayloadCapacityMin: ToInt32Ptr(filter.PayloadCapacityMin),
		PayloadCapacityMax: ToInt32Ptr(filter.PayloadCapacityMax),
		FuelConsumptionMin: ToInt32Ptr(filter.FuelConsumptionMin),
		FuelConsumptionMax: ToInt32Ptr(filter.FuelConsumptionMax),
		SortBy:             filter.SortBy.Value,
		SortOrder:          filter.SortOrder.Value,
	}
	rows, err := db.queries.GetTransports(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var transports []models.Transport
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		transports = append(transports, convertGeneratedTransportRowToModel(row))
	}
	return transports, totalPages, nil
}

type UpdateTransportArgs struct {
	TransportID     int32
	Model           string
	LicensePlate    string
	PayloadCapacity int32
	FuelConsumption int32
}

func (db DB) UpdateTransport(ctx context.Context, args UpdateTransportArgs) error {
	arg := generated.UpdateTransportParams{
		TransportID:     args.TransportID,
		Model:           args.Model,
		LicensePlate:    args.LicensePlate,
		PayloadCapacity: args.PayloadCapacity,
		FuelConsumption: args.FuelConsumption,
	}
	return parseTransportsError(db.queries.UpdateTransport(ctx, arg))
}

func (db DB) SoftDeleteTransport(ctx context.Context, transportID int32) error {
	return db.queries.SoftDeleteTransport(ctx, transportID)
}

func (db DB) HardDeleteTransport(ctx context.Context, transportID int32) error {
	return db.queries.HardDeleteTransport(ctx, transportID)
}

func (db DB) RestoreTransport(ctx context.Context, transportID int32) error {
	return db.queries.RestoreTransport(ctx, transportID)
}

func (db DB) BulkSoftDeleteTransports(ctx context.Context, transportIDs []int32) error {
	return db.queries.BulkSoftDeleteTransports(ctx, transportIDs)
}

func (db DB) BulkHardDeleteTransports(ctx context.Context, transportIDs []int32) error {
	return db.queries.BulkHardDeleteTransports(ctx, transportIDs)
}

func (db DB) GetTransportOrders(ctx context.Context, transportID int32, limit, offset int32) ([]models.Order, int64, error) {
	arg := generated.GetTransportOrdersParams{
		TransportID: transportID,
		Limit:       limit,
		Offset:      offset,
	}
	rows, err := db.queries.GetTransportOrders(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var orders []models.Order
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		orders = append(orders, convertGetTransportOrdersRowToModel(row))
	}
	return orders, totalCount, nil
}

// conversion helpers
func convertGeneratedTransportToModel(t generated.Transport) models.Transport {
	return models.Transport{
		TransportID:     t.TransportID,
		Model:           t.Model,
		LicensePlate:    fromStringPtr(t.LicensePlate),
		PayloadCapacity: t.PayloadCapacity,
		FuelConsumption: t.FuelConsumption,
		CreatedAt:       fromPgTimestamptz(t.CreatedAt),
		UpdatedAt:       fromPgTimestamptz(t.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(t.DeletedAt),
	}
}

func convertGeneratedTransportRowToModel(row generated.GetTransportsRow) models.Transport {
	return models.Transport{
		TransportID:     row.TransportID,
		Model:           row.Model,
		LicensePlate:    fromStringPtr(row.LicensePlate),
		PayloadCapacity: row.PayloadCapacity,
		FuelConsumption: row.FuelConsumption,
		CreatedAt:       fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:       fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(row.DeletedAt),
	}
}

func convertGetTransportOrdersRowToModel(row generated.GetTransportOrdersRow) models.Order {
	return models.Order{
		OrderID:     row.OrderID,
		ClientID:    row.ClientID,
		TransportID: row.TransportID,
		EmployeeID:  row.EmployeeID,
		Grade:       uint8(row.Grade),
		Distance:    row.Distance,
		Weight:      row.Weight,
		TotalPrice:  row.TotalPrice,
		PriceID:     row.PriceID,
		Status:      models.OrderStatus(row.Status),
		CreatedAt:   fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:   fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:   fromPgTimestamptz(row.DeletedAt),
	}
}

func (db DB) ListFreeTransports(ctx context.Context) ([]ui.ListItem, error) {
	rows, err := db.queries.ListFreeTransports(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ui.ListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, ui.ListItem{
			ID: r.TransportID,
			Name: strings.Join([]string{r.Model,
				fromStringPtr(r.LicensePlate),
				strconv.FormatInt(int64(r.PayloadCapacity), 10), "kg"}, " "),
		})
	}
	return items, nil
}

func parseTransportsError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "transports_license_plate_key" {
			return ErrDuplicateLicense
		}
		return fmt.Errorf("unhandled error: %w", err)
	}
	return fmt.Errorf("uknown error: %w", err)
}
