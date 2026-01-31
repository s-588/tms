package db

import (
	"context"

	"github.com/s-588/tms/db/generated"
	"github.com/s-588/tms/models"
	"github.com/shopspring/decimal"
)

func (db DB) CreateOrderWrapper(ctx context.Context, distance, weight int, totalPrice, status string) (models.Order, error) {
	totalPriceDec, err := decimal.NewFromString(totalPrice)
	if err != nil {
		return models.Order{}, err
	}

	arg := generated.CreateOrderParams{
		Distance:   int32(distance),
		Weight:     int32(weight),
		TotalPrice: totalPriceDec,
		Status:     status,
	}

	genOrder, err := db.queries.CreateOrder(ctx, arg)
	if err != nil {
		return models.Order{}, err
	}

	return convertGeneratedOrderToModel(genOrder), nil
}

func (db DB) DeleteOrderWrapper(ctx context.Context, orderID int) error {
	return db.queries.DeleteOrder(ctx, int32(orderID))
}

func (db DB) DeleteOrderTransportAssignmentsWrapper(ctx context.Context, orderID int) error {
	return db.queries.DeleteOrderTransportAssignments(ctx, int32(orderID))
}

func (db DB) GetOrderByIDWrapper(ctx context.Context, orderID int) (models.Order, error) {
	genOrder, err := db.queries.GetOrderByorder_id(ctx, int32(orderID))
	if err != nil {
		return models.Order{}, err
	}

	return convertGeneratedOrderToModel(genOrder), nil
}

func (db DB) GetOrdersPaginatedWrapper(ctx context.Context, limit, offset int) ([]models.Order, int64, error) {
	arg := generated.GetOrderPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := db.queries.GetOrderPaginated(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		order := models.Order{
			OrderID:    row.OrderID,
			Distance:   row.Distance,
			Weight:     row.Weight,
			TotalPrice: row.TotalPrice.String(),
			Status:     row.Status,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
			DeletedAt:  row.DeletedAt,
		}
		orders = append(orders, order)
	}

	return orders, totalCount, nil
}

func (db DB) GetOrderTransportsWrapper(ctx context.Context, orderID, limit, offset int) ([]models.Transport, int64, error) {
	arg := generated.GetOrderTransportsParams{
		OrderID: int32(orderID),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	rows, err := db.queries.GetOrderTransports(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	var transports []models.Transport
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		transport := models.Transport{
			TransportID:     row.TransportID,
			EmployeeID:      row.EmployeeID,
			Model:           row.Model,
			LicensePlate:    row.LicensePlate,
			PayloadCapacity: row.PayloadCapacity,
			FuelID:          row.FuelID,
			FuelConsumption: row.FuelConsumption,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
			DeletedAt:       row.DeletedAt,
		}
		transports = append(transports, transport)
	}

	return transports, totalCount, nil
}

func (db DB) InsertOrderTransportAssignmentsWrapper(ctx context.Context, orderID int, transportIDs []int) error {
	arg := generated.InsertOrderTransportAssignmentsParams{
		OrderID: int32(orderID),
		Column2: convertIntSliceToInt32(transportIDs),
	}

	return db.queries.InsertOrderTransportAssignments(ctx, arg)
}

func (db DB) UpdateOrderWrapper(ctx context.Context, orderID, distance, weight int, totalPrice, status string) error {
	totalPriceDec, err := decimal.NewFromString(totalPrice)
	if err != nil {
		return err
	}

	arg := generated.UpdateOrderParams{
		OrderID:    int32(orderID),
		Distance:   int32(distance),
		Weight:     int32(weight),
		TotalPrice: totalPriceDec,
		Status:     status,
	}

	return db.queries.UpdateOrder(ctx, arg)
}
func convertGeneratedOrderToModel(genOrder generated.Order) models.Order {
	return models.Order{
		OrderID:    genOrder.OrderID,
		Distance:   genOrder.Distance,
		Weight:     genOrder.Weight,
		TotalPrice: genOrder.TotalPrice.String(),
		Status:     genOrder.Status,
		CreatedAt:  genOrder.CreatedAt,
		UpdatedAt:  genOrder.UpdatedAt,
		DeletedAt:  genOrder.DeletedAt,
	}
}
