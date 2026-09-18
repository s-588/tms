package grpc

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateInspection adds a new inspection record.
func (s *Server) CreateInspection(ctx context.Context, req *pb.CreateInspectionRequest) (*pb.Inspection, error) {
	args := db.CreateInspectionArgs{
		TransportID:          req.TransportId,
		InspectionDate:       fromProtoTimestamp(req.InspectionDate),
		InspectionExpiration: fromProtoTimestamp(req.InspectionExpiration),
		Status:               fromProtoInspectionStatus(req.Status),
	}
	i, err := s.DB.CreateInspection(ctx, args)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoInspection(i), nil
}

// GetInspection retrieves an inspection by ID.
func (s *Server) GetInspection(ctx context.Context, req *pb.GetInspectionRequest) (*pb.Inspection, error) {
	i, err := s.DB.GetInspectionByID(ctx, req.InspectionId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoInspection(i), nil
}

// ListInspections returns a paginated list of inspections.
func (s *Server) ListInspections(ctx context.Context, req *pb.ListInspectionsRequest) (*pb.ListInspectionsResponse, error) {
	filter := models.InspectionFilter{}
	if req.TransportId != nil {
		filter.TransportID.SetValue(*req.TransportId)
	}
	if req.Status != nil {
		filter.Status.SetValue(fromProtoInspectionStatus(*req.Status))
	}
	if req.InspectionDateFrom != nil {
		filter.InspectionDateFrom.SetValue(fromProtoTimestamp(req.InspectionDateFrom))
	}
	if req.InspectionDateTo != nil {
		filter.InspectionDateTo.SetValue(fromProtoTimestamp(req.InspectionDateTo))
	}
	if req.InspectionExpirationFrom != nil {
		filter.InspectionExpirationFrom.SetValue(fromProtoTimestamp(req.InspectionExpirationFrom))
	}
	if req.InspectionExpirationTo != nil {
		filter.InspectionExpirationTo.SetValue(fromProtoTimestamp(req.InspectionExpirationTo))
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	inspections, totalPages, err := s.DB.GetInspections(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListInspectionsResponse{TotalPages: totalPages}
	for _, i := range inspections {
		resp.Inspections = append(resp.Inspections, toProtoInspection(i))
	}
	return resp, nil
}

// UpdateInspection modifies an existing inspection.
func (s *Server) UpdateInspection(ctx context.Context, req *pb.UpdateInspectionRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdateInspection(ctx, db.UpdateInspectionArgs{
		InspectionID:         req.InspectionId,
		TransportID:          req.TransportId,
		InspectionDate:       fromProtoTimestamp(req.InspectionDate),
		InspectionExpiration: fromProtoTimestamp(req.InspectionExpiration),
		Status:               fromProtoInspectionStatus(req.Status),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteInspection soft-deletes an inspection by ID.
func (s *Server) DeleteInspection(ctx context.Context, req *pb.DeleteInspectionRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeleteInspection(ctx, req.InspectionId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}
