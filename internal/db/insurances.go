package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

// CreateInsurance inserts a new insurance record.
func (db DB) CreateInsurance(ctx context.Context,
	transportID int32,
	insuranceDate time.Time,
	insuranceExpiration time.Time,
	payment decimal.Decimal,
	coverage decimal.Decimal,
) (models.Insurance, error) {
	arg := generated.CreateInsuranceParams{
		TransportID:         transportID,
		InsuranceDate:       pgtype.Date{Time: insuranceDate, Valid: true},
		InsuranceExpiration: pgtype.Date{Time: insuranceExpiration, Valid: true},
		Payment:             payment,
		Coverage:            coverage,
	}
	genInsurance, err := db.queries.CreateInsurance(ctx, arg)
	if err != nil {
		return models.Insurance{}, err
	}
	return convertGeneratedInsuranceToModel(genInsurance), nil
}

// GetInsuranceByID returns an insurance by its ID.
func (db DB) GetInsuranceByID(ctx context.Context, insuranceID int) (models.Insurance, error) {
	genInsurance, err := db.queries.GetInsurance(ctx, int32(insuranceID))
	if err != nil {
		return models.Insurance{}, err
	}
	return convertGeneratedInsuranceToModel(genInsurance), nil
}

// GetInsuranceByTransport returns the most recent insurance for a given transport.
func (db DB) GetInsuranceByTransport(ctx context.Context, transportID int) (models.Insurance, error) {
	genInsurance, err := db.queries.GetInsuranceByTransport(ctx, int32(transportID))
	if err != nil {
		return models.Insurance{}, err
	}
	return convertGeneratedInsuranceToModel(genInsurance), nil
}

// GetInsurances returns a paginated list of insurances matching the filter.
func (db DB) GetInsurances(ctx context.Context, limit, offset int, filter models.InsuranceFilter) ([]models.Insurance, int64, error) {
	arg := generated.GetInsurancesParams{
		Limit:                   int32(limit),
		Offset:                  int32(offset),
		TransportIDFilter:       filter.TransportID.ToPtr(),
		InsuranceDateFrom:       optionalTimeToPgDate(filter.InsuranceDateFrom),
		InsuranceDateTo:         optionalTimeToPgDate(filter.InsuranceDateTo),
		InsuranceExpirationFrom: optionalTimeToPgDate(filter.InsuranceExpirationFrom),
		InsuranceExpirationTo:   optionalTimeToPgDate(filter.InsuranceExpirationTo),
		PaymentMin:              optionalDecimalToPgNumeric(filter.PaymentMin),
		PaymentMax:              optionalDecimalToPgNumeric(filter.PaymentMax),
		CoverageMin:             optionalDecimalToPgNumeric(filter.CoverageMin),
		CoverageMax:             optionalDecimalToPgNumeric(filter.CoverageMax),
		CreatedFrom:             optionalTimeToPgTimestamptz(filter.CreatedFrom),
		CreatedTo:               optionalTimeToPgTimestamptz(filter.CreatedTo),
		UpdatedFrom:             optionalTimeToPgTimestamptz(filter.UpdatedFrom),
		UpdatedTo:               optionalTimeToPgTimestamptz(filter.UpdatedTo),
		SortBy:                  filter.SortBy.ToPtr(),
		SortOrder:               filter.SortOrder.ToPtr(),
	}
	rows, err := db.queries.GetInsurances(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var insurances []models.Insurance
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		insurances = append(insurances, convertGeneratedInsuranceRowToModel(row))
	}
	return insurances, totalCount, nil
}

// UpdateInsurance updates mutable fields of an insurance.
func (db DB) UpdateInsurance(ctx context.Context,
	insuranceID int,
	transportID *int32,
	insuranceDate *time.Time,
	insuranceExpiration *time.Time,
	payment *decimal.Decimal,
	coverage *decimal.Decimal,
) error {
	arg := generated.UpdateInsuranceParams{
		InsuranceID:         int32(insuranceID),
		TransportID:         transportID,
		InsuranceDate:       timeToPgDatePtr(insuranceDate),
		InsuranceExpiration: timeToPgDatePtr(insuranceExpiration),
		Payment:             fromDecimalPtrToPgNumeric(payment),
		Coverage:            fromDecimalPtrToPgNumeric(coverage),
	}
	return db.queries.UpdateInsurance(ctx, arg)
}

// SoftDeleteInsurance marks an insurance as deleted.
func (db DB) SoftDeleteInsurance(ctx context.Context, insuranceID int) error {
	return db.queries.SoftDeleteInsurance(ctx, int32(insuranceID))
}

// HardDeleteInsurance permanently removes an insurance.
func (db DB) HardDeleteInsurance(ctx context.Context, insuranceID int) error {
	return db.queries.HardDeleteInsurance(ctx, int32(insuranceID))
}

// RestoreInsurance removes the soft‑delete mark.
func (db DB) RestoreInsurance(ctx context.Context, insuranceID int) error {
	return db.queries.RestoreInsurance(ctx, int32(insuranceID))
}

// BulkSoftDeleteInsurances soft‑deletes multiple insurances.
func (db DB) BulkSoftDeleteInsurances(ctx context.Context, insuranceIDs []int) error {
	return db.queries.BulkSoftDeleteInsurances(ctx, convertIntSliceToInt32(insuranceIDs))
}

// BulkHardDeleteInsurances permanently deletes multiple insurances.
func (db DB) BulkHardDeleteInsurances(ctx context.Context, insuranceIDs []int) error {
	return db.queries.BulkHardDeleteInsurances(ctx, convertIntSliceToInt32(insuranceIDs))
}

// conversion helpers
func convertGeneratedInsuranceToModel(i generated.Insurance) models.Insurance {
	return models.Insurance{
		InsuranceID:         i.InsuranceID,
		TransportID:         i.TransportID,
		InsuranceDate:       i.InsuranceDate.Time,
		InsuranceExpiration: i.InsuranceExpiration.Time,
		Payment:             i.Payment,
		Coverage:            i.Coverage,
		CreatedAt:           i.CreatedAt.Time,
		UpdatedAt:           fromPgTimestamptz(i.UpdatedAt),
		DeletedAt:           fromPgTimestamptz(i.DeletedAt),
	}
}

func convertGeneratedInsuranceRowToModel(row generated.GetInsurancesRow) models.Insurance {
	return models.Insurance{
		InsuranceID:         row.InsuranceID,
		TransportID:         row.TransportID,
		InsuranceDate:       row.InsuranceDate.Time,
		InsuranceExpiration: row.InsuranceExpiration.Time,
		Payment:             row.Payment,
		Coverage:            row.Coverage,
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:           fromPgTimestamptz(row.DeletedAt),
	}
}
