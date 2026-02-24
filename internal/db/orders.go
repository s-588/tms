package db

import (
	"context"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

// TODO: create struct for args
func (db DB) CreateOrder(ctx context.Context,
	clientID int32,
	transportID int32,
	employeeID int32,
	grade uint8,
	distance int32,
	weight int32,
	totalPrice decimal.Decimal,
	priceID int32,
	status models.OrderStatus,
	nodeIDStart, nodeIDEnd int32,
) (models.Order, error) {
	arg := generated.CreateOrderParams{
		ClientID:    clientID,
		TransportID: transportID,
		EmployeeID:  employeeID,
		Grade:       int16(grade),
		Distance:    distance,
		Weight:      weight,
		TotalPrice:  totalPrice,
		PriceID:     priceID,
		Status:      generated.OrderStatus(status),
		NodeIDStart: nodeIDStart,
		NodeIDEnd:   nodeIDEnd,
	}
	genOrder, err := db.queries.CreateOrder(ctx, arg)
	if err != nil {
		return models.Order{}, err
	}
	return convertGeneratedOrderToModel(genOrder), nil
}

func (db DB) GetOrderByID(ctx context.Context, orderID int) (models.Order, error) {
	genOrder, err := db.queries.GetOrder(ctx, int32(orderID))
	if err != nil {
		return models.Order{}, err
	}
	return convertGeneratedOrderToModel(genOrder), nil
}

func (db DB) GetOrders(ctx context.Context, limit, offset int, filter models.OrderFilter) ([]models.Order, int64, error) {
	arg := generated.GetOrdersParams{
		Limit:             int32(limit),
		Offset:            int32(offset),
		StatusFilter:      ToNullOrderStatus(filter.Status),
		TotalPriceMin:     ToPgTypeNumeric(filter.TotalPriceMin),
		TotalPriceMax:     ToPgTypeNumeric(filter.TotalPriceMax),
		DistanceMin:       ToInt32Ptr(filter.DistanceMin),
		DistanceMax:       ToInt32Ptr(filter.DistanceMax),
		WeightMin:         ToInt32Ptr(filter.WeightMin),
		WeightMax:         ToInt32Ptr(filter.WeightMax),
		ClientIDFilter:    ToInt32Ptr(filter.ClientID),
		TransportIDFilter: ToInt32Ptr(filter.TransportID),
		EmployeeIDFilter:  ToInt32Ptr(filter.EmployeeID),
		PriceIDFilter:     ToInt32Ptr(filter.PriceID),
		GradeMin:          ToInt16PtrFromUint8(filter.GradeMin),
		GradeMax:          ToInt16PtrFromUint8(filter.GradeMax),
		CreatedFrom:       ToPgTypeTimestamptz(filter.CreatedFrom),
		CreatedTo:         ToPgTypeTimestamptz(filter.CreatedTo),
		UpdatedFrom:       ToPgTypeTimestamptz(filter.UpdatedFrom),
		UpdatedTo:         ToPgTypeTimestamptz(filter.UpdatedTo),
		SortOrder:         ToStringPtr(filter.SortOrder),
		SortBy:            ToStringPtr(filter.SortBy),
	}
	rows, err := db.queries.GetOrders(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var orders []models.Order
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		orders = append(orders, convertGeneratedOrderRowToModel(row))
	}
	return orders, totalCount, nil
}

func (db DB) UpdateOrder(ctx context.Context,
	orderID int32,
	clientID models.Optional[int32],
	transportID models.Optional[int32],
	employeeID models.Optional[int32],
	grade models.Optional[uint8],
	distance models.Optional[int32],
	weight models.Optional[int32],
	totalPrice models.Optional[decimal.Decimal],
	priceID models.Optional[int32],
	status models.Optional[models.OrderStatus],
	nodeIDStart, nodeIDEnd models.Optional[int32],
) error {
	arg := generated.UpdateOrderParams{
		OrderID:     orderID,
		ClientID:    ToInt32Ptr(clientID),
		TransportID: ToInt32Ptr(transportID),
		EmployeeID:  ToInt32Ptr(employeeID),
		Grade:       ToInt16PtrFromUint8(grade),
		Distance:    ToInt32Ptr(distance),
		Weight:      ToInt32Ptr(weight),
		TotalPrice:  ToPgTypeNumeric(totalPrice),
		PriceID:     ToInt32Ptr(priceID),
		Status:      ToNullOrderStatus(status),
		NodeIDStart: ToInt32Ptr(nodeIDStart),
		NodeIDEnd:   ToInt32Ptr(nodeIDEnd),
	}
	return db.queries.UpdateOrder(ctx, arg)
}

func (db DB) SoftDeleteOrder(ctx context.Context, orderID int) error {
	return db.queries.SoftDeleteOrder(ctx, int32(orderID))
}

func (db DB) HardDeleteOrder(ctx context.Context, orderID int) error {
	return db.queries.HardDeleteOrder(ctx, int32(orderID))
}

func (db DB) RestoreOrder(ctx context.Context, orderID int) error {
	return db.queries.RestoreOrder(ctx, int32(orderID))
}

func (db DB) BulkSoftDeleteOrders(ctx context.Context, orderIDs []int) error {
	return db.queries.BulkSoftDeleteOrders(ctx, convertIntSliceToInt32(orderIDs))
}

func (db DB) BulkHardDeleteOrders(ctx context.Context, orderIDs []int) error {
	return db.queries.BulkHardDeleteOrders(ctx, convertIntSliceToInt32(orderIDs))
}

func (db DB) UpdateOrderStatus(ctx context.Context, orderID int, status models.OrderStatus) error {
	return db.queries.UpdateOrderStatus(ctx, generated.UpdateOrderStatusParams{
		OrderID: int32(orderID),
		Status:  generated.OrderStatus(status),
	})
}

// conversion helpers
func convertGeneratedOrderToModel(o generated.Order) models.Order {
	return models.Order{
		OrderID:     o.OrderID,
		ClientID:    o.ClientID,
		TransportID: o.TransportID,
		EmployeeID:  o.EmployeeID,
		Grade:       uint8(o.Grade),
		Distance:    o.Distance,
		Weight:      o.Weight,
		TotalPrice:  o.TotalPrice,
		PriceID:     o.PriceID,
		Status:      models.OrderStatus(o.Status),
		NodeIDStart: o.NodeIDStart,
		NodeIDEnd:   o.NodeIDEnd,
		CreatedAt:   fromPgTimestamptz(o.CreatedAt),
		UpdatedAt:   fromPgTimestamptz(o.UpdatedAt),
		DeletedAt:   fromPgTimestamptz(o.DeletedAt),
	}
}

func convertGeneratedOrderRowToModel(row generated.GetOrdersRow) models.Order {
	return models.Order{
		OrderID:     row.OrderID,
		ClientID:    row.ClientID,
		TransportID: row.TransportID,
		EmployeeID:  row.EmployeeID,
		Grade:       uint8(row.Grade),
		Distance:    row.Distance,
		Weight:      row.Weight,
		TotalPrice:  row.TotalPrice,
		PriceID:     row.PriceID,
		Status:      models.OrderStatus(row.Status),
		NodeIDStart: row.NodeIDStart,
		NodeIDEnd:   row.NodeIDEnd,
		CreatedAt:   fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:   fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:   fromPgTimestamptz(row.DeletedAt),
	}
}
