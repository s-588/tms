package db

import (
	"context"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
)

type CreateInspectionArgs struct {
	TransportID          int32
	InspectionDate       time.Time
	InspectionExpiration time.Time
	Status               models.InspectionStatus
}

func (db DB) CreateInspection(ctx context.Context, args CreateInspectionArgs) (models.Inspection, error) {
	arg := generated.CreateInspectionParams{
		TransportID:          args.TransportID,
		InspectionDate:       ToPgTypeDateFromTime(args.InspectionDate),
		InspectionExpiration: ToPgTypeDateFromTime(args.InspectionExpiration),
		Status:               generated.InspectionStatus(args.Status),
	}
	genInspection, err := db.queries.CreateInspection(ctx, arg)
	if err != nil {
		return models.Inspection{}, err
	}
	return convertGeneratedInspectionToModel(genInspection), nil
}

func (db DB) GetInspectionByID(ctx context.Context, inspectionID int32) (models.Inspection, error) {
	genInspection, err := db.queries.GetInspection(ctx, inspectionID)
	if err != nil {
		return models.Inspection{}, err
	}
	return convertGeneratedInspectionToModel(genInspection), nil
}

func (db DB) GetInspectionsByTransport(ctx context.Context, transportID int32) ([]models.Inspection, error) {
	genInspections, err := db.queries.GetInspectionsByTransport(ctx, transportID)
	if err != nil {
		return nil, err
	}
	var inspections []models.Inspection
	for _, gi := range genInspections {
		inspections = append(inspections, convertGeneratedInspectionToModel(gi))
	}
	return inspections, nil
}

func (db DB) GetInspections(ctx context.Context, page int32, filter models.InspectionFilter) ([]models.Inspection, int32, error) {
	arg := generated.GetInspectionsParams{
		Page:                     page,
		TransportIDFilter:        filter.TransportID.ToPtr(),
		StatusFilter:             ToNullInspectionStatus(filter.Status),
		InspectionDateFrom:       optionalTimeToPgDate(filter.InspectionDateFrom),
		InspectionDateTo:         optionalTimeToPgDate(filter.InspectionDateTo),
		InspectionExpirationFrom: optionalTimeToPgDate(filter.InspectionExpirationFrom),
		InspectionExpirationTo:   optionalTimeToPgDate(filter.InspectionExpirationTo),
		SortBy:                   filter.GetSortBy(),
		SortOrder:                filter.GetSortOrder(),
	}
	rows, err := db.queries.GetInspections(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var inspections []models.Inspection
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		inspections = append(inspections, convertGeneratedInspectionRowToModel(row))
	}
	return inspections, totalPages, nil
}

type UpdateInspectionArgs struct {
	InspectionID         int32
	TransportID          int32
	InspectionDate       time.Time
	InspectionExpiration time.Time
	Status               models.InspectionStatus
}

func (db DB) UpdateInspection(ctx context.Context, args UpdateInspectionArgs) error {
	arg := generated.UpdateInspectionParams{
		InspectionID:         args.InspectionID,
		TransportID:          args.TransportID,
		InspectionDate:       ToPgTypeDateFromTime(args.InspectionDate),
		InspectionExpiration: ToPgTypeDateFromTime(args.InspectionExpiration),
		Status:               generated.InspectionStatus(args.Status),
	}
	return db.queries.UpdateInspection(ctx, arg)
}

func (db DB) SoftDeleteInspection(ctx context.Context, inspectionID int32) error {
	return db.queries.SoftDeleteInspection(ctx, inspectionID)
}

func (db DB) HardDeleteInspection(ctx context.Context, inspectionID int32) error {
	return db.queries.HardDeleteInspection(ctx, inspectionID)
}

func (db DB) RestoreInspection(ctx context.Context, inspectionID int32) error {
	return db.queries.RestoreInspection(ctx, inspectionID)
}

func (db DB) BulkSoftDeleteInspections(ctx context.Context, inspectionIDs []int32) error {
	return db.queries.BulkSoftDeleteInspections(ctx, inspectionIDs)
}

func (db DB) BulkHardDeleteInspections(ctx context.Context, inspectionIDs []int32) error {
	return db.queries.BulkHardDeleteInspections(ctx, inspectionIDs)
}

// conversion helpers
func convertGeneratedInspectionToModel(i generated.Inspection) models.Inspection {
	return models.Inspection{
		InspectionID:         i.InspectionID,
		TransportID:          i.TransportID,
		InspectionDate:       i.InspectionDate.Time,
		InspectionExpiration: i.InspectionExpiration.Time,
		Status:               models.InspectionStatus(i.Status),
		CreatedAt:            i.CreatedAt.Time,
		UpdatedAt:            fromPgTimestamptz(i.UpdatedAt),
		DeletedAt:            fromPgTimestamptz(i.DeletedAt),
	}
}

func convertGeneratedInspectionRowToModel(row generated.GetInspectionsRow) models.Inspection {
	return models.Inspection{
		InspectionID:         row.InspectionID,
		TransportID:          row.TransportID,
		InspectionDate:       row.InspectionDate.Time,
		InspectionExpiration: row.InspectionExpiration.Time,
		Status:               models.InspectionStatus(row.Status),
		CreatedAt:            row.CreatedAt.Time,
		UpdatedAt:            fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:            fromPgTimestamptz(row.DeletedAt),
	}
}

func inspectionStatusToNullGenerated(o models.Optional[models.InspectionStatus]) generated.NullInspectionStatus {
	if !o.Set {
		return generated.NullInspectionStatus{Valid: false}
	}
	return generated.NullInspectionStatus{
		InspectionStatus: generated.InspectionStatus(o.Value),
		Valid:            true,
	}
}

func inspectionStatusPtrToNullGenerated(p *models.InspectionStatus) generated.NullInspectionStatus {
	if p == nil {
		return generated.NullInspectionStatus{Valid: false}
	}
	return generated.NullInspectionStatus{
		InspectionStatus: generated.InspectionStatus(*p),
		Valid:            true,
	}
}
