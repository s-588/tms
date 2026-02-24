package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
)

// CreateInspection inserts a new inspection.
func (db DB) CreateInspection(ctx context.Context,
	transportID int32,
	inspectionDate time.Time,
	inspectionExpiration time.Time,
	status models.InspectionStatus,
) (models.Inspection, error) {
	arg := generated.CreateInspectionParams{
		TransportID:          transportID,
		InspectionDate:       pgtype.Date{Time: inspectionDate, Valid: true},
		InspectionExpiration: pgtype.Date{Time: inspectionExpiration, Valid: true},
		Status:               generated.InspectionStatus(status),
	}
	genInspection, err := db.queries.CreateInspection(ctx, arg)
	if err != nil {
		return models.Inspection{}, err
	}
	return convertGeneratedInspectionToModel(genInspection), nil
}

// GetInspectionByID returns an inspection by its ID.
func (db DB) GetInspectionByID(ctx context.Context, inspectionID int) (models.Inspection, error) {
	genInspection, err := db.queries.GetInspection(ctx, int32(inspectionID))
	if err != nil {
		return models.Inspection{}, err
	}
	return convertGeneratedInspectionToModel(genInspection), nil
}

// GetInspectionsByTransport returns all inspections for a given transport, ordered by date descending.
func (db DB) GetInspectionsByTransport(ctx context.Context, transportID int) ([]models.Inspection, error) {
	genInspections, err := db.queries.GetInspectionsByTransport(ctx, int32(transportID))
	if err != nil {
		return nil, err
	}
	var inspections []models.Inspection
	for _, gi := range genInspections {
		inspections = append(inspections, convertGeneratedInspectionToModel(gi))
	}
	return inspections, nil
}

// GetInspections returns a paginated list of inspections matching the filter.
func (db DB) GetInspections(ctx context.Context, limit, offset int, filter models.InspectionFilter) ([]models.Inspection, int64, error) {
	arg := generated.GetInspectionsParams{
		Limit:                    int32(limit),
		Offset:                   int32(offset),
		TransportIDFilter:        filter.TransportID.ToPtr(),
		StatusFilter:             inspectionStatusToNullGenerated(filter.Status),
		InspectionDateFrom:       optionalTimeToPgDate(filter.InspectionDateFrom),
		InspectionDateTo:         optionalTimeToPgDate(filter.InspectionDateTo),
		InspectionExpirationFrom: optionalTimeToPgDate(filter.InspectionExpirationFrom),
		InspectionExpirationTo:   optionalTimeToPgDate(filter.InspectionExpirationTo),
		CreatedFrom:              optionalTimeToPgTimestamptz(filter.CreatedFrom),
		CreatedTo:                optionalTimeToPgTimestamptz(filter.CreatedTo),
		UpdatedFrom:              optionalTimeToPgTimestamptz(filter.UpdatedFrom),
		UpdatedTo:                optionalTimeToPgTimestamptz(filter.UpdatedTo),
		SortBy:                   filter.SortBy.ToPtr(),
		SortOrder:                filter.SortOrder.ToPtr(),
	}
	rows, err := db.queries.GetInspections(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var inspections []models.Inspection
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		inspections = append(inspections, convertGeneratedInspectionRowToModel(row))
	}
	return inspections, totalCount, nil
}

// UpdateInspection updates mutable fields of an inspection.
func (db DB) UpdateInspection(ctx context.Context,
	inspectionID int,
	transportID *int32,
	inspectionDate *time.Time,
	inspectionExpiration *time.Time,
	status *models.InspectionStatus,
) error {
	arg := generated.UpdateInspectionParams{
		InspectionID:         int32(inspectionID),
		TransportID:          transportID,
		InspectionDate:       timeToPgDatePtr(inspectionDate),
		InspectionExpiration: timeToPgDatePtr(inspectionExpiration),
		Status:               inspectionStatusPtrToNullGenerated(status),
	}
	return db.queries.UpdateInspection(ctx, arg)
}

// SoftDeleteInspection marks an inspection as deleted.
func (db DB) SoftDeleteInspection(ctx context.Context, inspectionID int) error {
	return db.queries.SoftDeleteInspection(ctx, int32(inspectionID))
}

// HardDeleteInspection permanently removes an inspection.
func (db DB) HardDeleteInspection(ctx context.Context, inspectionID int) error {
	return db.queries.HardDeleteInspection(ctx, int32(inspectionID))
}

// RestoreInspection removes the soft‑delete mark.
func (db DB) RestoreInspection(ctx context.Context, inspectionID int) error {
	return db.queries.RestoreInspection(ctx, int32(inspectionID))
}

// BulkSoftDeleteInspections soft‑deletes multiple inspections.
func (db DB) BulkSoftDeleteInspections(ctx context.Context, inspectionIDs []int) error {
	return db.queries.BulkSoftDeleteInspections(ctx, convertIntSliceToInt32(inspectionIDs))
}

// BulkHardDeleteInspections permanently deletes multiple inspections.
func (db DB) BulkHardDeleteInspections(ctx context.Context, inspectionIDs []int) error {
	return db.queries.BulkHardDeleteInspections(ctx, convertIntSliceToInt32(inspectionIDs))
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
