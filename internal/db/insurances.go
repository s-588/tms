package db

import (
	"context"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

type CreateInsuranceArgs struct {
	TransportID         int32
	InsuranceDate       time.Time
	InsuranceExpiration time.Time
	Payment             decimal.Decimal
	Coverage            decimal.Decimal
}

func (db DB) CreateInsurance(ctx context.Context, args CreateInsuranceArgs) (models.Insurance, error) {
	arg := generated.CreateInsuranceParams{
		TransportID:         args.TransportID,
		InsuranceDate:       ToPgTypeDateFromTime(args.InsuranceDate),
		InsuranceExpiration: ToPgTypeDateFromTime(args.InsuranceExpiration),
		Payment:             args.Payment,
		Coverage:            args.Coverage,
	}
	genInsurance, err := db.queries.CreateInsurance(ctx, arg)
	if err != nil {
		return models.Insurance{}, err
	}
	return convertGeneratedInsuranceToModel(genInsurance), nil
}

func (db DB) GetInsuranceByID(ctx context.Context, insuranceID int32) (models.Insurance, error) {
	genInsurance, err := db.queries.GetInsurance(ctx, insuranceID)
	if err != nil {
		return models.Insurance{}, err
	}
	return convertGeneratedInsuranceToModel(genInsurance), nil
}

func (db DB) GetInsuranceByTransport(ctx context.Context, transportID int32) (models.Insurance, error) {
	genInsurance, err := db.queries.GetInsuranceByTransport(ctx, transportID)
	if err != nil {
		return models.Insurance{}, err
	}
	return convertGeneratedInsuranceToModel(genInsurance), nil
}

func (db DB) GetInsurances(ctx context.Context, page int32, filter models.InsuranceFilter) ([]models.Insurance, int32, error) {
	arg := generated.GetInsurancesParams{
		Page:                    page,
		TransportIDFilter:       ToInt32Ptr(filter.TransportID),
		InsuranceDateFrom:       optionalTimeToPgDate(filter.InsuranceDateFrom),
		InsuranceDateTo:         optionalTimeToPgDate(filter.InsuranceDateTo),
		InsuranceExpirationFrom: optionalTimeToPgDate(filter.InsuranceExpirationFrom),
		InsuranceExpirationTo:   optionalTimeToPgDate(filter.InsuranceExpirationTo),
		PaymentMin:              optionalDecimalToPgNumeric(filter.PaymentMin),
		PaymentMax:              optionalDecimalToPgNumeric(filter.PaymentMax),
		CoverageMin:             optionalDecimalToPgNumeric(filter.CoverageMin),
		CoverageMax:             optionalDecimalToPgNumeric(filter.CoverageMax),
		SortBy:                  filter.SortBy.Value,
		SortOrder:               filter.SortOrder.Value,
	}
	rows, err := db.queries.GetInsurances(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var insurances []models.Insurance
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		insurances = append(insurances, convertGeneratedInsuranceRowToModel(row))
	}
	return insurances, totalPages, nil
}

type UpdateInsuranceArgs struct {
	InsuranceID         int32
	TransportID         int32
	InsuranceDate       time.Time
	InsuranceExpiration time.Time
	Payment             decimal.Decimal
	Coverage            decimal.Decimal
}

func (db DB) UpdateInsurance(ctx context.Context, args UpdateInsuranceArgs) error {
	arg := generated.UpdateInsuranceParams{
		InsuranceID:         args.InsuranceID,
		TransportID:         args.TransportID,
		InsuranceDate:       ToPgTypeDateFromTime(args.InsuranceDate),
		InsuranceExpiration: ToPgTypeDateFromTime(args.InsuranceExpiration),
		Payment:             args.Payment,
		Coverage:            args.Coverage,
	}
	return db.queries.UpdateInsurance(ctx, arg)
}

func (db DB) SoftDeleteInsurance(ctx context.Context, insuranceID int32) error {
	return db.queries.SoftDeleteInsurance(ctx, insuranceID)
}

func (db DB) HardDeleteInsurance(ctx context.Context, insuranceID int32) error {
	return db.queries.HardDeleteInsurance(ctx, insuranceID)
}

func (db DB) RestoreInsurance(ctx context.Context, insuranceID int32) error {
	return db.queries.RestoreInsurance(ctx, insuranceID)
}

func (db DB) BulkSoftDeleteInsurances(ctx context.Context, insuranceIDs []int32) error {
	return db.queries.BulkSoftDeleteInsurances(ctx, insuranceIDs)
}

func (db DB) BulkHardDeleteInsurances(ctx context.Context, insuranceIDs []int32) error {
	return db.queries.BulkHardDeleteInsurances(ctx, insuranceIDs)
}

// conversion helpers
func convertGeneratedInsuranceToModel(i generated.Insurance) models.Insurance {
	return models.Insurance{
		InsuranceID:         i.InsuranceID,
		TransportID:         i.TransportID,
		InsuranceDate:       fromPgDate(i.InsuranceDate),
		InsuranceExpiration: fromPgDate(i.InsuranceExpiration),
		Payment:             i.Payment,
		Coverage:            i.Coverage,
		CreatedAt:           fromPgTimestamptz(i.CreatedAt),
		UpdatedAt:           fromPgTimestamptz(i.UpdatedAt),
		DeletedAt:           fromPgTimestamptz(i.DeletedAt),
	}
}

func convertGeneratedInsuranceRowToModel(row generated.GetInsurancesRow) models.Insurance {
	return models.Insurance{
		InsuranceID:         row.InsuranceID,
		TransportID:         row.TransportID,
		InsuranceDate:       fromPgDate(row.InsuranceDate),
		InsuranceExpiration: fromPgDate(row.InsuranceExpiration),
		Payment:             row.Payment,
		Coverage:            row.Coverage,
		CreatedAt:           fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:           fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:           fromPgTimestamptz(row.DeletedAt),
	}
}
