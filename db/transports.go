package db

import (
	"context"

	"github.com/s-588/tms/db/generated"
	"github.com/s-588/tms/models"
)

func (db DB) CreateTransportWrapper(ctx context.Context, employeeID *int, model string, licensePlate *string, payloadCapacity, fuelID, fuelConsumption int) (models.Transport, error) {
	var empID *int32
	if employeeID != nil {
		temp := int32(*employeeID)
		empID = &temp
	}

	arg := generated.CreateTransportParams{
		EmployeeID:      empID,
		Model:           model,
		LicensePlate:    licensePlate,
		PayloadCapacity: int32(payloadCapacity),
		FuelID:          int32(fuelID),
		FuelConsumption: int32(fuelConsumption),
	}

	genTransport, err := db.queries.CreateTransport(ctx, arg)
	if err != nil {
		return models.Transport{}, err
	}

	return convertGeneratedTransportToModel(genTransport), nil
}

func (db DB) DeleteTransportWrapper(ctx context.Context, transportID int) error {
	return db.queries.DeleteTransport(ctx, int32(transportID))
}

func (db DB) GetTransportByIDWrapper(ctx context.Context, transportID int) (models.Transport, error) {
	genTransport, err := db.queries.GetTransportBytransport_id(ctx, int32(transportID))
	if err != nil {
		return models.Transport{}, err
	}

	return convertGeneratedTransportToModel(genTransport), nil
}

func (db DB) GetTransportsPaginatedWrapper(ctx context.Context, limit, offset int) ([]models.Transport, int64, error) {
	arg := generated.GetTransportsPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := db.queries.GetTransportsPaginated(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	var transports []models.Transport
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		transport := models.Transport{
			TransportID:     row.TransportID,
			EmployeeID:      row.EmployeeID,
			Model:           row.Model,
			LicensePlate:    row.LicensePlate,
			PayloadCapacity: row.PayloadCapacity,
			FuelID:          row.FuelID,
			FuelConsumption: row.FuelConsumption,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
			DeletedAt:       row.DeletedAt,
		}
		transports = append(transports, transport)
	}

	return transports, totalCount, nil
}

func (db DB) UpdateTransportWrapper(ctx context.Context, transportID int, employeeID *int, model string, licensePlate *string, payloadCapacity, fuelID, fuelConsumption int) error {
	var empID *int32
	if employeeID != nil {
		temp := int32(*employeeID)
		empID = &temp
	}

	arg := generated.UpdateTransportParams{
		TransportID:     int32(transportID),
		EmployeeID:      empID,
		Model:           model,
		LicensePlate:    licensePlate,
		PayloadCapacity: int32(payloadCapacity),
		FuelID:          int32(fuelID),
		FuelConsumption: int32(fuelConsumption),
	}

	return db.queries.UpdateTransport(ctx, arg)
}

func convertGeneratedTransportToModel(genTransport generated.Transport) models.Transport {
	return models.Transport{
		TransportID:     genTransport.TransportID,
		EmployeeID:      genTransport.EmployeeID,
		Model:           genTransport.Model,
		LicensePlate:    genTransport.LicensePlate,
		PayloadCapacity: genTransport.PayloadCapacity,
		FuelID:          genTransport.FuelID,
		FuelConsumption: genTransport.FuelConsumption,
		CreatedAt:       genTransport.CreatedAt,
		UpdatedAt:       genTransport.UpdatedAt,
		DeletedAt:       genTransport.DeletedAt,
	}
}
