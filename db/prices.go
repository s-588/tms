package db

import (
	"context"

	"github.com/s-588/tms/db/generated"
	"github.com/s-588/tms/models"
	"github.com/shopspring/decimal"
)

func (db DB) CreatePriceWrapper(ctx context.Context, cargoType, cost string, weight, distance int) (models.Price, error) {
	costDec, err := decimal.NewFromString(cost)
	if err != nil {
		return models.Price{}, err
	}

	arg := generated.CreatePriceParams{
		CargoType: cargoType,
		Cost:      costDec,
		Weight:    int32(weight),
		Distance:  int32(distance),
	}

	genPrice, err := db.queries.CreatePrice(ctx, arg)
	if err != nil {
		return models.Price{}, err
	}

	return convertGeneratedPriceToModel(genPrice), nil
}

func (db DB) DeletePriceWrapper(ctx context.Context, priceID int) error {
	return db.queries.DeletePrice(ctx, int32(priceID))
}

func (db DB) GetPriceByIDWrapper(ctx context.Context, priceID int) (models.Price, error) {
	genPrice, err := db.queries.GetPriceByprice_id(ctx, int32(priceID))
	if err != nil {
		return models.Price{}, err
	}

	return convertGeneratedPriceToModel(genPrice), nil
}

func (db DB) GetPricesPaginatedWrapper(ctx context.Context, limit, offset int) ([]models.Price, int64, error) {
	arg := generated.GetPricesPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := db.queries.GetPricesPaginated(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	var prices []models.Price
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		price := models.Price{
			PriceID:   row.PriceID,
			CargoType: row.CargoType,
			Cost:      row.Cost.String(),
			Weight:    row.Weight,
			Distance:  row.Distance,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			DeletedAt: row.DeletedAt,
		}
		prices = append(prices, price)
	}

	return prices, totalCount, nil
}

func (db DB) UpdatePriceWrapper(ctx context.Context, priceID int, cargoType, cost string, weight, distance int) error {
	costDec, err := decimal.NewFromString(cost)
	if err != nil {
		return err
	}

	arg := generated.UpdatePriceParams{
		PriceID:   int32(priceID),
		CargoType: cargoType,
		Cost:      costDec,
		Weight:    int32(weight),
		Distance:  int32(distance),
	}

	return db.queries.UpdatePrice(ctx, arg)
}

func convertGeneratedPriceToModel(genPrice generated.Price) models.Price {
	return models.Price{
		PriceID:   genPrice.PriceID,
		CargoType: genPrice.CargoType,
		Cost:      genPrice.Cost.String(),
		Weight:    genPrice.Weight,
		Distance:  genPrice.Distance,
		CreatedAt: genPrice.CreatedAt,
		UpdatedAt: genPrice.UpdatedAt,
		DeletedAt: genPrice.DeletedAt,
	}
}
