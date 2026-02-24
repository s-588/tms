package db

import (
	"context"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
)

func (db DB) CreateTransport(ctx context.Context,
	model string,
	licensePlate string,
	payloadCapacity int32,
	fuelConsumption int32,
	inspectionPassed bool,
	inspectionDate time.Time,
) (models.Transport, error) {
	arg := generated.CreateTransportParams{
		Model:            model,
		LicensePlate:     licensePlate,
		PayloadCapacity:  payloadCapacity,
		FuelConsumption:  fuelConsumption,
		InspectionPassed: inspectionPassed,
		InspectionDate:   ToPgTypeDateFromTime(inspectionDate),
	}
	genTransport, err := db.queries.CreateTransport(ctx, arg)
	if err != nil {
		return models.Transport{}, err
	}
	return convertGeneratedTransportToModel(genTransport), nil
}

func (db DB) GetTransportByID(ctx context.Context, transportID int) (models.Transport, error) {
	genTransport, err := db.queries.GetTransport(ctx, int32(transportID))
	if err != nil {
		return models.Transport{}, err
	}
	return convertGeneratedTransportToModel(genTransport), nil
}

func (db DB) GetTransports(ctx context.Context, limit, offset int, filter models.TransportFilter) ([]models.Transport, int64, error) {
	arg := generated.GetTransportsParams{
		Limit:                  int32(limit),
		Offset:                 int32(offset),
		ModelFilter:            ToStringPtr(filter.Model),
		LicensePlateFilter:     ToStringPtr(filter.LicensePlate),
		PayloadCapacityMin:     ToInt32Ptr(filter.PayloadCapacityMin),
		PayloadCapacityMax:     ToInt32Ptr(filter.PayloadCapacityMax),
		FuelConsumptionMin:     ToInt32Ptr(filter.FuelConsumptionMin),
		FuelConsumptionMax:     ToInt32Ptr(filter.FuelConsumptionMax),
		InspectionPassedFilter: ToBoolPtr(filter.InspectionPassed),
		InspectionDateFrom:     ToPgTypeDate(filter.InspectionDateFrom),
		InspectionDateTo:       ToPgTypeDate(filter.InspectionDateTo),
		CreatedFrom:            ToPgTypeTimestamptz(filter.CreatedFrom),
		CreatedTo:              ToPgTypeTimestamptz(filter.CreatedTo),
		UpdatedFrom:            ToPgTypeTimestamptz(filter.UpdatedFrom),
		UpdatedTo:              ToPgTypeTimestamptz(filter.UpdatedTo),
		SortOrder:              ToStringPtr(filter.SortOrder),
		SortBy:                 ToStringPtr(filter.SortBy),
	}
	rows, err := db.queries.GetTransports(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var transports []models.Transport
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		transports = append(transports, convertGeneratedTransportRowToModel(row))
	}
	return transports, totalCount, nil
}

func (db DB) UpdateTransport(ctx context.Context,
	transportID int,
	model models.Optional[string],
	licensePlate models.Optional[string],
	payloadCapacity models.Optional[int32],
	fuelConsumption models.Optional[int32],
	inspectionPassed models.Optional[bool],
	inspectionDate models.Optional[time.Time],
) error {
	arg := generated.UpdateTransportParams{
		TransportID:      int32(transportID),
		Model:            ToStringPtr(model),
		LicensePlate:     ToStringPtr(licensePlate),
		PayloadCapacity:  ToInt32Ptr(payloadCapacity),
		FuelConsumption:  ToInt32Ptr(fuelConsumption),
		InspectionPassed: ToBoolPtr(inspectionPassed),
		InspectionDate:   ToPgTypeDate(inspectionDate),
	}
	return db.queries.UpdateTransport(ctx, arg)
}

func (db DB) SoftDeleteTransport(ctx context.Context, transportID int) error {
	return db.queries.SoftDeleteTransport(ctx, int32(transportID))
}

func (db DB) HardDeleteTransport(ctx context.Context, transportID int) error {
	return db.queries.HardDeleteTransport(ctx, int32(transportID))
}

func (db DB) RestoreTransport(ctx context.Context, transportID int) error {
	return db.queries.RestoreTransport(ctx, int32(transportID))
}

func (db DB) BulkSoftDeleteTransports(ctx context.Context, transportIDs []int) error {
	return db.queries.BulkSoftDeleteTransports(ctx, convertIntSliceToInt32(transportIDs))
}

func (db DB) BulkHardDeleteTransports(ctx context.Context, transportIDs []int) error {
	return db.queries.BulkHardDeleteTransports(ctx, convertIntSliceToInt32(transportIDs))
}

func (db DB) GetTransportOrders(ctx context.Context, transportID, limit, offset int) ([]models.Order, int64, error) {
	arg := generated.GetTransportOrdersParams{
		TransportID: int32(transportID),
		Limit:       int32(limit),
		Offset:      int32(offset),
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
		LicensePlate:    *t.LicensePlate,
		PayloadCapacity: t.PayloadCapacity,
		FuelConsumption: t.FuelConsumption,
		CreatedAt:       fromPgTimestamptz(t.CreatedAt),
		UpdatedAt:       fromPgTimestamptz(t.UpdatedAt),
		DeletedAt:       fromPgTimestamptz(t.DeletedAt),
	}
}

func convertGeneratedTransportRowToModel(row generated.GetTransportsRow) models.Transport {
	var p string
	if row.LicensePlate != nil {
		p = *row.LicensePlate
	}
	return models.Transport{

		TransportID:     row.TransportID,
		Model:           row.Model,
		LicensePlate:    p,
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
		NodeIDStart: row.NodeIDStart,
		NodeIDEnd:   row.NodeIDEnd,
		CreatedAt:   fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:   fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:   fromPgTimestamptz(row.DeletedAt),
	}
}
