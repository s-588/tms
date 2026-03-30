package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

type CreatePriceArgs struct {
	CargoType string
	Weight    decimal.Decimal
	Distance  decimal.Decimal
}

func (db DB) CreatePrice(ctx context.Context, args CreatePriceArgs) (models.Price, error) {
	arg := generated.CreatePriceParams{
		CargoType: args.CargoType,
		Weight:    args.Weight,
		Distance:  args.Distance,
	}
	genPrice, err := db.queries.CreatePrice(ctx, arg)
	if err != nil {
		return models.Price{}, parsePricesError(err)
	}
	return convertGeneratedPriceToModel(genPrice), nil
}

func (db DB) GetPriceByID(ctx context.Context, priceID int32) (models.Price, error) {
	genPrice, err := db.queries.GetPrice(ctx, priceID)
	if err != nil {
		return models.Price{}, err
	}
	return convertGeneratedPriceToModel(genPrice), nil
}

func (db DB) GetPriceByUnique(ctx context.Context, cargoType string, weight, distance decimal.Decimal) (models.Price, error) {
	arg := generated.GetPriceByUniqueParams{
		CargoType: cargoType,
		Weight:    weight,
		Distance:  distance,
	}
	genPrice, err := db.queries.GetPriceByUnique(ctx, arg)
	if err != nil {
		return models.Price{}, err
	}
	return convertGeneratedPriceToModel(genPrice), nil
}

func (db DB) GetPrices(ctx context.Context, page int32, filter models.PriceFilter) ([]models.Price, int32, error) {
	arg := generated.GetPricesParams{
		Page:            page,
		CargoTypeFilter: ToStringPtr(filter.CargoType),
		WeightMin:       ToInt32Ptr(filter.WeightMin),
		WeightMax:       ToInt32Ptr(filter.WeightMax),
		DistanceMin:     ToInt32Ptr(filter.DistanceMin),
		DistanceMax:     ToInt32Ptr(filter.DistanceMax),
		SortBy:          filter.SortBy.Value,
		SortOrder:       filter.SortOrder.Value,
	}
	rows, err := db.queries.GetPrices(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var prices []models.Price
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		prices = append(prices, convertGeneratedPriceRowToModel(row))
	}
	return prices, totalPages, nil
}

type UpdatePriceArgs struct {
	PriceID   int32
	CargoType string
	Weight    decimal.Decimal
	Distance  decimal.Decimal
}

func (db DB) UpdatePrice(ctx context.Context, args UpdatePriceArgs) error {
	arg := generated.UpdatePriceParams{
		PriceID:   args.PriceID,
		CargoType: args.CargoType,
		Weight:    args.Weight,
		Distance:  args.Distance,
	}
	return parsePricesError(db.queries.UpdatePrice(ctx, arg))
}

func (db DB) SoftDeletePrice(ctx context.Context, priceID int32) error {
	return db.queries.SoftDeletePrice(ctx, priceID)
}

func (db DB) HardDeletePrice(ctx context.Context, priceID int32) error {
	return db.queries.HardDeletePrice(ctx, priceID)
}

func (db DB) RestorePrice(ctx context.Context, priceID int32) error {
	return db.queries.RestorePrice(ctx, priceID)
}

func (db DB) BulkSoftDeletePrices(ctx context.Context, priceIDs []int32) error {
	return db.queries.BulkSoftDeletePrices(ctx, priceIDs)
}

func (db DB) BulkHardDeletePrices(ctx context.Context, priceIDs []int32) error {
	return db.queries.BulkHardDeletePrices(ctx, priceIDs)
}

// conversion helpers
func convertGeneratedPriceToModel(p generated.Price) models.Price {
	return models.Price{
		PriceID:   p.PriceID,
		CargoType: p.CargoType,
		Weight:    p.Weight,
		Distance:  p.Distance,
		CreatedAt: fromPgTimestamptz(p.CreatedAt),
		UpdatedAt: fromPgTimestamptz(p.UpdatedAt),
		DeletedAt: fromPgTimestamptz(p.DeletedAt),
	}
}

func convertGeneratedPriceRowToModel(row generated.GetPricesRow) models.Price {
	return models.Price{
		PriceID:   row.PriceID,
		CargoType: row.CargoType,
		Weight:    row.Weight,
		Distance:  row.Distance,
		CreatedAt: fromPgTimestamptz(row.CreatedAt),
		UpdatedAt: fromPgTimestamptz(row.UpdatedAt),
		DeletedAt: fromPgTimestamptz(row.DeletedAt),
	}
}
func (db DB) ListPrices(ctx context.Context) ([]ui.ListItem, error) {
	rows, err := db.queries.ListPrices(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ui.ListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, ui.ListItem{
			ID:   r.PriceID,
			Name: fmt.Sprintf("%s %s %s",r.CargoType,r.Weight,r.Distance),
		})
	}
	return items, nil
}

func parsePricesError(err error) error {
	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == "23505" {
		if pgErr.ConstraintName == "prices_cargo_type_weight_distance_key" {
			return ErrDuplicatePrice
		}
		return fmt.Errorf("unhandled error: %w", err)
	}
	return fmt.Errorf("uknown error %w", err)
}
