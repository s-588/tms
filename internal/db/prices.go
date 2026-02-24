package db

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

// CreatePrice inserts a new price.
func (db DB) CreatePrice(ctx context.Context, cargoType string, weight, distance int32) (models.Price, error) {
	arg := generated.CreatePriceParams{
		CargoType: cargoType,
		Weight:    weight,
		Distance:  distance,
	}
	genPrice, err := db.queries.CreatePrice(ctx, arg)
	if err != nil {
		return models.Price{}, err
	}
	return convertGeneratedPriceToModel(genPrice), nil
}

// GetPriceByID returns a price by ID.
func (db DB) GetPriceByID(ctx context.Context, priceID int) (models.Price, error) {
	genPrice, err := db.queries.GetPrice(ctx, int32(priceID))
	if err != nil {
		return models.Price{}, err
	}
	return convertGeneratedPriceToModel(genPrice), nil
}

// GetPriceByUnique returns a price by the unique combination (cargo_type, weight, distance).
func (db DB) GetPriceByUnique(ctx context.Context, cargoType string, weight, distance int32) (models.Price, error) {
	genPrice, err := db.queries.GetPriceByUnique(ctx, generated.GetPriceByUniqueParams{
		CargoType: cargoType,
		Weight:    weight,
		Distance:  distance,
	})
	if err != nil {
		return models.Price{}, err
	}
	return convertGeneratedPriceToModel(genPrice), nil
}

// GetPrices returns a paginated list of prices matching the filter.
func (db DB) GetPrices(ctx context.Context, limit, offset int, filter models.PriceFilter) ([]models.Price, int64, error) {
	arg := generated.GetPricesParams{
		Limit:           int32(limit),
		Offset:          int32(offset),
		CargoTypeFilter: ToStringPtr(filter.CargoType),
		WeightMin:       ToInt32Ptr(filter.WeightMin),
		WeightMax:       ToInt32Ptr(filter.WeightMax),
		DistanceMin:     ToInt32Ptr(filter.DistanceMin),
		DistanceMax:     ToInt32Ptr(filter.DistanceMax),
		CreatedFrom:     ToPgTypeTimestamptz(filter.CreatedFrom),
		CreatedTo:       ToPgTypeTimestamptz(filter.CreatedTo),
		UpdatedFrom:     ToPgTypeTimestamptz(filter.UpdatedFrom),
		UpdatedTo:       ToPgTypeTimestamptz(filter.UpdatedTo),
		SortOrder:       ToStringPtr(filter.SortOrder),
		SortBy:          ToStringPtr(filter.SortBy),
	}
	rows, err := db.queries.GetPrices(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var prices []models.Price
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		prices = append(prices, convertGeneratedPriceRowToModel(row))
	}
	return prices, totalCount, nil
}

// UpdatePrice updates mutable fields of a price.
func (db DB) UpdatePrice(ctx context.Context, priceID int,
	cargoType models.Optional[string],
	weight, distance models.Optional[int32],
) error {
	arg := generated.UpdatePriceParams{
		PriceID:   int32(priceID),
		CargoType: ToStringPtr(cargoType),
		Weight:    ToInt32Ptr(weight),
		Distance:  ToInt32Ptr(distance),
	}
	return db.queries.UpdatePrice(ctx, arg)
}

// SoftDeletePrice marks a price as deleted.
func (db DB) SoftDeletePrice(ctx context.Context, priceID int) error {
	return db.queries.SoftDeletePrice(ctx, int32(priceID))
}

// HardDeletePrice permanently removes a price.
func (db DB) HardDeletePrice(ctx context.Context, priceID int) error {
	return db.queries.HardDeletePrice(ctx, int32(priceID))
}

// RestorePrice removes the soft-delete mark.
func (db DB) RestorePrice(ctx context.Context, priceID int) error {
	return db.queries.RestorePrice(ctx, int32(priceID))
}

// BulkSoftDeletePrices soft-deletes multiple prices.
func (db DB) BulkSoftDeletePrices(ctx context.Context, priceIDs []int) error {
	return db.queries.BulkSoftDeletePrices(ctx, convertIntSliceToInt32(priceIDs))
}

// BulkHardDeletePrices permanently deletes multiple prices.
func (db DB) BulkHardDeletePrices(ctx context.Context, priceIDs []int) error {
	return db.queries.BulkHardDeletePrices(ctx, convertIntSliceToInt32(priceIDs))
}

// convertGeneratedPriceToModel maps a generated.Price to models.Price.
func convertGeneratedPriceToModel(p generated.Price) models.Price {
	return models.Price{
		PriceID:   p.PriceID,
		CargoType: p.CargoType,
		Weight:    decimal.NewFromInt(int64(p.Weight)),
		Distance:  decimal.NewFromInt(int64(p.Distance)),
		CreatedAt: fromPgTimestamptz(p.CreatedAt),
		UpdatedAt: fromPgTimestamptz(p.UpdatedAt),
		DeletedAt: fromPgTimestamptz(p.DeletedAt),
	}
}

// convertGeneratedPriceRowToModel maps a generated.GetPricesRow to models.Price.
func convertGeneratedPriceRowToModel(row generated.GetPricesRow) models.Price {
	return models.Price{
		PriceID:   row.PriceID,
		CargoType: row.CargoType,
		Weight:    decimal.NewFromInt(int64(row.Weight)),
		Distance:  decimal.NewFromInt(int64(row.Distance)),
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
		UpdatedAt: fromPgTimestamptz(row.UpdatedAt),
		DeletedAt: fromPgTimestamptz(row.DeletedAt),
	}
}
