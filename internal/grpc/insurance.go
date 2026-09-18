package grpc

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	pb "github.com/s-588/tms/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateInsurance adds a new insurance record.
func (s *Server) CreateInsurance(ctx context.Context, req *pb.CreateInsuranceRequest) (*pb.Insurance, error) {
	args := db.CreateInsuranceArgs{
		TransportID:         req.TransportId,
		InsuranceDate:       fromProtoTimestamp(req.InsuranceDate),
		InsuranceExpiration: fromProtoTimestamp(req.InsuranceExpiration),
		Payment:             fromProtoDecimal(req.Payment),
		Coverage:            fromProtoDecimal(req.Coverage),
	}
	i, err := s.DB.CreateInsurance(ctx, args)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoInsurance(i), nil
}

// GetInsurance retrieves an insurance by ID.
func (s *Server) GetInsurance(ctx context.Context, req *pb.GetInsuranceRequest) (*pb.Insurance, error) {
	i, err := s.DB.GetInsuranceByID(ctx, req.InsuranceId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoInsurance(i), nil
}

// ListInsurances returns a paginated list of insurances.
func (s *Server) ListInsurances(ctx context.Context, req *pb.ListInsurancesRequest) (*pb.ListInsurancesResponse, error) {
	filter := models.InsuranceFilter{}
	if req.TransportId != nil {
		filter.TransportID.SetValue(*req.TransportId)
	}
	if req.InsuranceDateFrom != nil {
		filter.InsuranceDateFrom.SetValue(fromProtoTimestamp(req.InsuranceDateFrom))
	}
	if req.InsuranceDateTo != nil {
		filter.InsuranceDateTo.SetValue(fromProtoTimestamp(req.InsuranceDateTo))
	}
	if req.InsuranceExpirationFrom != nil {
		filter.InsuranceExpirationFrom.SetValue(fromProtoTimestamp(req.InsuranceExpirationFrom))
	}
	if req.InsuranceExpirationTo != nil {
		filter.InsuranceExpirationTo.SetValue(fromProtoTimestamp(req.InsuranceExpirationTo))
	}
	if req.PaymentMin != nil {
		filter.PaymentMin.SetValue(fromProtoDecimal(*req.PaymentMin))
	}
	if req.PaymentMax != nil {
		filter.PaymentMax.SetValue(fromProtoDecimal(*req.PaymentMax))
	}
	if req.CoverageMin != nil {
		filter.CoverageMin.SetValue(fromProtoDecimal(*req.CoverageMin))
	}
	if req.CoverageMax != nil {
		filter.CoverageMax.SetValue(fromProtoDecimal(*req.CoverageMax))
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	insurances, totalPages, err := s.DB.GetInsurances(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListInsurancesResponse{TotalPages: totalPages}
	for _, i := range insurances {
		resp.Insurances = append(resp.Insurances, toProtoInsurance(i))
	}
	return resp, nil
}

// UpdateInsurance modifies an existing insurance.
func (s *Server) UpdateInsurance(ctx context.Context, req *pb.UpdateInsuranceRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdateInsurance(ctx, db.UpdateInsuranceArgs{
		InsuranceID:         req.InsuranceId,
		TransportID:         req.TransportId,
		InsuranceDate:       fromProtoTimestamp(req.InsuranceDate),
		InsuranceExpiration: fromProtoTimestamp(req.InsuranceExpiration),
		Payment:             fromProtoDecimal(req.Payment),
		Coverage:            fromProtoDecimal(req.Coverage),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteInsurance soft-deletes an insurance by ID.
func (s *Server) DeleteInsurance(ctx context.Context, req *pb.DeleteInsuranceRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeleteInsurance(ctx, req.InsuranceId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}
