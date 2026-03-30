package db

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

type CreateOrderArg struct {
	ClientID    int32
	TransportID int32
	EmployeeID  int32
	Grade       uint8
	Distance    float64
	Weight      int32
	TotalPrice  decimal.Decimal
	PriceID     int32
	Status      models.OrderStatus
	NodeIDStart int32
	NodeIDEnd   int32
}

// CreateOrder inserts a new order.
// nodeIDStart and nodeIDEnd are optional; pass models.Optional[int32]{} to set NULL.
func (db DB) CreateOrder(ctx context.Context, args CreateOrderArg) (models.Order, error) {
	arg := generated.CreateOrderParams{
		ClientID:    args.ClientID,
		TransportID: args.TransportID,
		EmployeeID:  args.EmployeeID,
		Grade:       int16(args.Grade),
		Distance:    args.Distance,
		Weight:      args.Weight,
		TotalPrice:  args.TotalPrice,
		PriceID:     args.PriceID,
		Status:      generated.OrderStatus(args.Status),
		NodeIDStart: args.NodeIDStart,
		NodeIDEnd:   args.NodeIDEnd,
	}
	genOrder, err := db.queries.CreateOrder(ctx, arg)
	if err != nil {
		return models.Order{}, err
	}
	return db.GetOrderByID(ctx, genOrder.OrderID)
}

func (db DB) GetOrderByID(ctx context.Context, orderID int32) (models.Order, error) {
	genOrder, err := db.queries.GetOrder(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}
	return models.Order{
		OrderID:               genOrder.OrderID,
		ClientID:              genOrder.ClientID,
		TransportID:           genOrder.TransportID,
		EmployeeID:            genOrder.EmployeeID,
		Grade:                 uint8(genOrder.Grade),
		Distance:              genOrder.Distance,
		Weight:                genOrder.Weight,
		TotalPrice:            genOrder.TotalPrice,
		PriceID:               genOrder.PriceID,
		Status:                models.OrderStatus(genOrder.Status),
		CreatedAt:             fromPgTimestamptz(genOrder.CreatedAt),
		UpdatedAt:             fromPgTimestamptz(genOrder.UpdatedAt),
		DeletedAt:             fromPgTimestamptz(genOrder.DeletedAt),
		ClientName:            fromStringPtr(genOrder.ClientName),
		EmployeeName:          fromStringPtr(genOrder.EmployeeName),
		TransportLicensePlate: fromStringPtr(genOrder.TransportLicensePlate),
		PriceCargoType:        fromStringPtr(genOrder.PriceCargoType),
		NodeStartName:         fromStringPtr(genOrder.NodeStartName),
		NodeEndName:           fromStringPtr(genOrder.NodeEndName),
	}, nil
}

func (db DB) GetOrders(ctx context.Context, page int32, filter models.OrderFilter) ([]models.Order, int32, error) {
	arg := generated.GetOrdersParams{
		Page:              page,
		StatusFilter:      ToNullOrderStatus(filter.Status),
		TotalPriceMin:     optionalDecimalToPgNumeric(filter.TotalPriceMin),
		TotalPriceMax:     optionalDecimalToPgNumeric(filter.TotalPriceMax),
		DistanceMin:       ToFloat64Ptr(filter.DistanceMin),
		DistanceMax:       ToFloat64Ptr(filter.DistanceMax),
		WeightMin:         ToInt32Ptr(filter.WeightMin),
		WeightMax:         ToInt32Ptr(filter.WeightMax),
		ClientIDFilter:    ToInt32Ptr(filter.ClientID),
		TransportIDFilter: ToInt32Ptr(filter.TransportID),
		EmployeeIDFilter:  ToInt32Ptr(filter.EmployeeID),
		PriceIDFilter:     ToInt32Ptr(filter.PriceID),
		GradeMin:          ToInt16PtrFromUint8(filter.GradeMin),
		GradeMax:          ToInt16PtrFromUint8(filter.GradeMax),
		SortBy:            filter.SortBy.Value,
		SortOrder:         filter.SortOrder.Value,
	}
	rows, err := db.queries.GetOrders(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var orders []models.Order
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		orders = append(orders, models.Order{
			OrderID:               row.OrderID,
			ClientID:              row.ClientID,
			TransportID:           row.TransportID,
			EmployeeID:            row.EmployeeID,
			Grade:                 uint8(row.Grade),
			Distance:              row.Distance,
			Weight:                row.Weight,
			TotalPrice:            row.TotalPrice,
			PriceID:               row.PriceID,
			Status:                models.OrderStatus(row.Status),
			NodeIDStart:           row.NodeIDStart,
			NodeIDEnd:             row.NodeIDEnd,
			CreatedAt:             row.CreatedAt.Time,
			UpdatedAt:             row.UpdatedAt.Time,
			DeletedAt:             row.DeletedAt.Time,
			ClientName:            fromStringPtr(row.ClientName),
			EmployeeName:          fromStringPtr(row.EmployeeName),
			TransportLicensePlate: fromStringPtr(row.TransportLicensePlate),
			PriceCargoType:        fromStringPtr(row.PriceCargoType),
			NodeStartName:         fromStringPtr(row.NodeStartName),
			NodeEndName:           fromStringPtr(row.NodeEndName),
		})
	}
	return orders, totalPages, nil
}

type UpdateOrderArgs struct {
	OrderID     int32
	ClientID    int32
	TransportID int32
	EmployeeID  int32
	Grade       uint8
	Distance    float64
	Weight      int32
	TotalPrice  decimal.Decimal
	PriceID     int32
	Status      models.OrderStatus
	NodeIDStart int32
	NodeIDEnd   int32
}

// UpdateOrder performs a full update of an order.
// nodeIDStart and nodeIDEnd are now required (use 0 if not applicable – ensure 0 is not a valid FK).
func (db DB) UpdateOrder(ctx context.Context, args UpdateOrderArgs) error {
	arg := generated.UpdateOrderParams{
		OrderID:     args.OrderID,
		ClientID:    args.ClientID,
		TransportID: args.TransportID,
		EmployeeID:  args.EmployeeID,
		Grade:       int16(args.Grade),
		Distance:    args.Distance,
		Weight:      args.Weight,
		TotalPrice:  args.TotalPrice,
		PriceID:     args.PriceID,
		Status:      generated.OrderStatus(args.Status),
		NodeIDStart: args.NodeIDStart,
		NodeIDEnd:   args.NodeIDEnd,
	}
	return db.queries.UpdateOrder(ctx, arg)
}

// SoftDeleteOrder marks an order as deleted.
func (db DB) SoftDeleteOrder(ctx context.Context, orderID int32) error {
	return db.queries.SoftDeleteOrder(ctx, orderID)
}

// HardDeleteOrder permanently removes an order.
func (db DB) HardDeleteOrder(ctx context.Context, orderID int32) error {
	return db.queries.HardDeleteOrder(ctx, orderID)
}

// RestoreOrder removes the soft‑delete mark.
func (db DB) RestoreOrder(ctx context.Context, orderID int32) error {
	return db.queries.RestoreOrder(ctx, orderID)
}

// BulkSoftDeleteOrders soft‑deletes multiple orders.
func (db DB) BulkSoftDeleteOrders(ctx context.Context, orderIDs []int32) error {
	return db.queries.BulkSoftDeleteOrders(ctx, orderIDs)
}

// BulkHardDeleteOrders permanently deletes multiple orders.
func (db DB) BulkHardDeleteOrders(ctx context.Context, orderIDs []int32) error {
	return db.queries.BulkHardDeleteOrders(ctx, orderIDs)
}

// UpdateOrderStatus updates only the status of an order.
func (db DB) UpdateOrderStatus(ctx context.Context, orderID int32, status models.OrderStatus) error {
	return db.queries.UpdateOrderStatus(ctx, generated.UpdateOrderStatusParams{
		OrderID: orderID,
		Status:  generated.OrderStatus(status),
	})
}
