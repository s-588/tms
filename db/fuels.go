package db

import (
	"context"

	"github.com/s-588/tms/db/generated"
	"github.com/s-588/tms/models"
	"github.com/shopspring/decimal"
)

func (db DB) CreateFuelWrapper(ctx context.Context, name string, supplier *string, price string) (models.Fuel, error) {
	priceDec, err := decimal.NewFromString(price)
	if err != nil {
		return models.Fuel{}, err
	}

	arg := generated.CreateFuelParams{
		Name:     name,
		Supplier: supplier,
		Price:    priceDec,
	}

	genFuel, err := db.queries.CreateFuel(ctx, arg)
	if err != nil {
		return models.Fuel{}, err
	}

	return convertGeneratedFuelToModel(genFuel), nil
}

func (db DB) DeleteFuelWrapper(ctx context.Context, fuelID int) error {
	return db.queries.DeleteFuel(ctx, int32(fuelID))
}

func (db DB) GetFuelByIDWrapper(ctx context.Context, fuelID int) (models.Fuel, error) {
	genFuel, err := db.queries.GetFuelByfuel_id(ctx, int32(fuelID))
	if err != nil {
		return models.Fuel{}, err
	}

	return convertGeneratedFuelToModel(genFuel), nil
}

func (db DB) GetFuelsPaginatedWrapper(ctx context.Context, limit, offset int) ([]models.Fuel, int64, error) {
	arg := generated.GetFuelsPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := db.queries.GetFuelsPaginated(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	var fuels []models.Fuel
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		fuel := models.Fuel{
			FuelID:    row.FuelID,
			Name:      row.Name,
			Supplier:  row.Supplier,
			Price:     row.Price.String(),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			DeletedAt: row.DeletedAt,
		}
		fuels = append(fuels, fuel)
	}

	return fuels, totalCount, nil
}

func (db DB) UpdateFuelWrapper(ctx context.Context, fuelID int, name string, supplier *string, price string) error {
	priceDec, err := decimal.NewFromString(price)
	if err != nil {
		return err
	}

	arg := generated.UpdateFuelParams{
		FuelID:   int32(fuelID),
		Name:     name,
		Supplier: supplier,
		Price:    priceDec,
	}

	return db.queries.UpdateFuel(ctx, arg)
}
func convertGeneratedFuelToModel(genFuel generated.Fuel) models.Fuel {
	return models.Fuel{
		FuelID:    genFuel.FuelID,
		Name:      genFuel.Name,
		Supplier:  genFuel.Supplier,
		Price:     genFuel.Price.String(),
		CreatedAt: genFuel.CreatedAt,
		UpdatedAt: genFuel.UpdatedAt,
		DeletedAt: genFuel.DeletedAt,
	}
}
