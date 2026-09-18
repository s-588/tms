package grpc

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreatePrice adds a new price configuration.
func (s *Server) CreatePrice(ctx context.Context, req *pb.CreatePriceRequest) (*pb.Price, error) {
	args := db.CreatePriceArgs{
		CargoType: req.CargoType,
		Weight:    fromProtoDecimal(req.Weight),
		Distance:  fromProtoDecimal(req.Distance),
	}
	p, err := s.DB.CreatePrice(ctx, args)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoPrice(p), nil
}

// GetPrice retrieves a price configuration by ID.
func (s *Server) GetPrice(ctx context.Context, req *pb.GetPriceRequest) (*pb.Price, error) {
	p, err := s.DB.GetPriceByID(ctx, req.PriceId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoPrice(p), nil
}

// ListPrices returns a paginated list of prices matching the filter.
func (s *Server) ListPrices(ctx context.Context, req *pb.ListPricesRequest) (*pb.ListPricesResponse, error) {
	filter := models.PriceFilter{}
	if req.CargoType != nil {
		filter.CargoType.SetValue(*req.CargoType)
	}
	if req.WeightMin != nil {
		filter.WeightMin.SetValue(*req.WeightMin)
	}
	if req.WeightMax != nil {
		filter.WeightMax.SetValue(*req.WeightMax)
	}
	if req.DistanceMin != nil {
		filter.DistanceMin.SetValue(*req.DistanceMin)
	}
	if req.DistanceMax != nil {
		filter.DistanceMax.SetValue(*req.DistanceMax)
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	prices, totalPages, err := s.DB.GetPrices(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListPricesResponse{TotalPages: totalPages}
	for _, p := range prices {
		resp.Prices = append(resp.Prices, toProtoPrice(p))
	}
	return resp, nil
}

// UpdatePrice modifies an existing price configuration.
func (s *Server) UpdatePrice(ctx context.Context, req *pb.UpdatePriceRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdatePrice(ctx, db.UpdatePriceArgs{
		PriceID:   req.PriceId,
		CargoType: req.CargoType,
		Weight:    fromProtoDecimal(req.Weight),
		Distance:  fromProtoDecimal(req.Distance),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeletePrice soft-deletes a price configuration by ID.
func (s *Server) DeletePrice(ctx context.Context, req *pb.DeletePriceRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeletePrice(ctx, req.PriceId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}
