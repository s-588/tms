package grpc

import (
	"context"
	"fmt"
	"math"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/tms"
	pb "github.com/s-588/tms/proto"
	"github.com/shopspring/decimal"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateOrder creates a new order, automatically calculating distance and total price.
func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	// Calculate distance between nodes.
	distance, err := s.DB.CalculateDistance(ctx, req.NodeIdStart, req.NodeIdEnd)
	if err != nil {
		return nil, mapError(err)
	}

	// Fetch transport to get fuel consumption and payload capacity.
	transport, err := s.DB.GetTransportByID(ctx, req.TransportId)
	if err != nil {
		return nil, mapError(err)
	}

	// Compute total price using the TMS pricing engine.
	totalPrice, err := tms.CalculateOrderCost(ctx, s.DB, tms.CalculateOrderCostArgs{
		ClientID:        req.ClientId,
		PriceID:         req.PriceId,
		Weight:          int64(req.Weight),
		FuelConsumption: transport.FuelConsumption,
		PayloadCapacity: transport.PayloadCapacity,
		NodeStartID:     req.NodeIdStart,
		NodeEndID:       req.NodeIdEnd,
	})
	if err != nil {
		return nil, mapError(err)
	}

	arg := db.CreateOrderArg{
		ClientID:    req.ClientId,
		TransportID: req.TransportId,
		EmployeeID:  req.EmployeeId,
		Grade:       0, // default grade; can be updated later
		Distance:    distance,
		Weight:      req.Weight,
		TotalPrice:  totalPrice,
		PriceID:     req.PriceId,
		Status:      fromProtoOrderStatus(req.Status),
		NodeIDStart: req.NodeIdStart,
		NodeIDEnd:   req.NodeIdEnd,
	}

	order, err := s.DB.CreateOrder(ctx, arg)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoOrder(order), nil
}

// GetOrder retrieves an order by ID.
func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.DB.GetOrderByID(ctx, req.OrderId)
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoOrder(order), nil
}

// ListOrders returns a paginated list of orders matching the filter.
func (s *Server) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	filter := models.OrderFilter{}
	if req.ClientId != nil {
		filter.ClientID.SetValue(*req.ClientId)
	}
	if req.TransportId != nil {
		filter.TransportID.SetValue(*req.TransportId)
	}
	if req.EmployeeId != nil {
		filter.EmployeeID.SetValue(*req.EmployeeId)
	}
	if req.PriceId != nil {
		filter.PriceID.SetValue(*req.PriceId)
	}
	if req.DistanceMin != nil {
		filter.DistanceMin.SetValue(*req.DistanceMin)
	}
	if req.DistanceMax != nil {
		filter.DistanceMax.SetValue(*req.DistanceMax)
	}
	if req.WeightMin != nil {
		filter.WeightMin.SetValue(*req.WeightMin)
	}
	if req.WeightMax != nil {
		filter.WeightMax.SetValue(*req.WeightMax)
	}
	if req.TotalPriceMin != nil {
		filter.TotalPriceMin.SetValue(fromProtoDecimal(*req.TotalPriceMin))
	}
	if req.TotalPriceMax != nil {
		filter.TotalPriceMax.SetValue(fromProtoDecimal(*req.TotalPriceMax))
	}
	if req.GradeMin != nil {
		if *req.GradeMin < 0 || *req.GradeMin > math.MaxUint8 {
			return nil, mapError(fmt.Errorf("grade_min must be between 0 and 255"))
		}
		filter.GradeMin.SetValue(uint8(*req.GradeMin))
	}
	if req.GradeMax != nil {
		if *req.GradeMax < 0 || *req.GradeMax > math.MaxUint8 {
			return nil, mapError(fmt.Errorf("grade_max must be between 0 and 255"))
		}
		filter.GradeMax.SetValue(uint8(*req.GradeMax))
	}
	if req.Status != nil {
		filter.Status.SetValue(fromProtoOrderStatus(*req.Status))
	}
	if req.SortBy != "" {
		filter.SortBy.SetValue(req.SortBy)
	}
	if req.SortOrder != "" {
		filter.SortOrder.SetValue(req.SortOrder)
	}

	orders, totalPages, err := s.DB.GetOrders(ctx, req.Page, filter)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &pb.ListOrdersResponse{TotalPages: totalPages}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, toProtoOrder(o))
	}
	return resp, nil
}

// UpdateOrder modifies an existing order. It recalculates the total price if key fields change.
func (s *Server) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*emptypb.Empty, error) {
	// First fetch the existing order to preserve grade and other computed fields.
	existing, err := s.DB.GetOrderByID(ctx, req.OrderId)
	if err != nil {
		return nil, mapError(err)
	}

	// Recalculate distance and price if nodes or transport/weight have changed.
	var distance float64
	var totalPrice decimal.Decimal
	if req.NodeIdStart != existing.NodeIDStart || req.NodeIdEnd != existing.NodeIDEnd ||
		req.TransportId != existing.TransportID || req.Weight != existing.Weight {
		// Recalculate distance.
		distance, err = s.DB.CalculateDistance(ctx, req.NodeIdStart, req.NodeIdEnd)
		if err != nil {
			return nil, mapError(err)
		}
		// Fetch transport for cost calculation.
		transport, err := s.DB.GetTransportByID(ctx, req.TransportId)
		if err != nil {
			return nil, mapError(err)
		}
		totalPrice, err = tms.CalculateOrderCost(ctx, s.DB, tms.CalculateOrderCostArgs{
			ClientID:        req.ClientId,
			PriceID:         req.PriceId,
			Weight:          int64(req.Weight),
			FuelConsumption: transport.FuelConsumption,
			PayloadCapacity: transport.PayloadCapacity,
			NodeStartID:     req.NodeIdStart,
			NodeEndID:       req.NodeIdEnd,
		})
		if err != nil {
			return nil, mapError(err)
		}
	} else {
		// Keep existing values.
		distance = existing.Distance
		totalPrice = existing.TotalPrice
	}

	err = s.DB.UpdateOrder(ctx, db.UpdateOrderArgs{
		OrderID:     req.OrderId,
		ClientID:    req.ClientId,
		TransportID: req.TransportId,
		EmployeeID:  req.EmployeeId,
		Grade:       existing.Grade,
		Distance:    distance,
		Weight:      req.Weight,
		TotalPrice:  totalPrice,
		PriceID:     req.PriceId,
		Status:      fromProtoOrderStatus(req.Status),
		NodeIDStart: req.NodeIdStart,
		NodeIDEnd:   req.NodeIdEnd,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// UpdateOrderStatus changes only the status of an order.
func (s *Server) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*emptypb.Empty, error) {
	err := s.DB.UpdateOrderStatus(ctx, req.OrderId, fromProtoOrderStatus(req.Status))
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// DeleteOrder soft-deletes an order by ID.
func (s *Server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*emptypb.Empty, error) {
	err := s.DB.SoftDeleteOrder(ctx, req.OrderId)
	if err != nil {
		return nil, mapError(err)
	}
	return &emptypb.Empty{}, nil
}

// CalculateOrderCost returns the estimated cost for an order without creating it.
func (s *Server) CalculateOrderCost(ctx context.Context, req *pb.CalculateOrderCostRequest) (*pb.CalculateOrderCostResponse, error) {
	price, err := tms.CalculateOrderCost(ctx, s.DB, tms.CalculateOrderCostArgs{
		ClientID:        req.ClientId,
		PriceID:         req.PriceId,
		Weight:          int64(req.Weight),
		FuelConsumption: req.FuelConsumption,
		PayloadCapacity: req.PayloadCapacity,
		NodeStartID:     req.NodeStartId,
		NodeEndID:       req.NodeEndId,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.CalculateOrderCostResponse{TotalPrice: toProtoDecimal(price)}, nil
}
