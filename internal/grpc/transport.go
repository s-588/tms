package grpc

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateTransport adds a new transport vehicle.
func (s *Server) CreateTransport(ctx context.Context, req *pb.CreateTransportRequest) (*pb.Transport, error) {
	args := db.CreateTransportArgs{
		Model:           req.Model,
		LicensePlate:    req.LicensePlate,
		PayloadCapacity: req.PayloadCapacity,
		FuelConsumption: req.FuelConsumption,
	}
	t, err := s.DB.CreateTransport(ctx, args)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoTransport(t), nil
}

// GetTransport retrieves a transport by ID.
func (s *Server) GetTransport(ctx context.Context, req *pb.GetTransportRequest) (*pb.Transport, error) {
	t, err := s.DB.GetTransportByID(ctx, req.TransportId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoTransport(t), nil
}

// ListTransports returns a paginated list of transports matching the filter.
func (s *Server) ListTransports(ctx context.Context, req *pb.ListTransportsRequest) (*pb.ListTransportsResponse, error) {
	filter := models.TransportFilter{}
	if req.Model != nil {
		filter.Model.SetValue(*req.Model)
	}
	if req.LicensePlate != nil {
		filter.LicensePlate.SetValue(*req.LicensePlate)
	}
	if req.PayloadCapacityMin != nil {
		filter.PayloadCapacityMin.SetValue(*req.PayloadCapacityMin)
	}
	if req.PayloadCapacityMax != nil {
		filter.PayloadCapacityMax.SetValue(*req.PayloadCapacityMax)
	}
	if req.FuelConsumptionMin != nil {
		filter.FuelConsumptionMin.SetValue(*req.FuelConsumptionMin)
	}
	if req.FuelConsumptionMax != nil {
		filter.FuelConsumptionMax.SetValue(*req.FuelConsumptionMax)
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	transports, totalPages, err := s.DB.GetTransports(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListTransportsResponse{TotalPages: totalPages}
	for _, t := range transports {
		resp.Transports = append(resp.Transports, toProtoTransport(t))
	}
	return resp, nil
}

// UpdateTransport modifies an existing transport.
func (s *Server) UpdateTransport(ctx context.Context, req *pb.UpdateTransportRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdateTransport(ctx, db.UpdateTransportArgs{
		TransportID:     req.TransportId,
		Model:           req.Model,
		LicensePlate:    req.LicensePlate,
		PayloadCapacity: req.PayloadCapacity,
		FuelConsumption: req.FuelConsumption,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteTransport soft-deletes a transport by ID.
func (s *Server) DeleteTransport(ctx context.Context, req *pb.DeleteTransportRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeleteTransport(ctx, req.TransportId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}
