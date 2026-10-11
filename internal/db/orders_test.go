package db_test

import (
	"context"
	"slices"
	"strconv"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestDB_Orders(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testOrders(t, database)
}

func testOrders(t *testing.T, database db.DB) {
	ctx := context.Background()

	order := createTestOrder(t, database)

	t.Run("get order by id", func(t *testing.T) {
		got, err := database.GetOrderByID(ctx, order.OrderID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.OrderID != order.OrderID {
			t.Fatalf("got order ID %d, want %d", got.OrderID, order.OrderID)
		}
	})

	t.Run("get all orders", func(t *testing.T) {
		list, total, err := database.GetOrders(ctx, 1, models.OrderFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total pages to be populated")
		}
		if !slices.ContainsFunc(list, func(v models.Order) bool {
			return v.OrderID == order.OrderID
		}) {
			t.Fatalf("slice doesn't contain wanted order; got %v, want %v", list, order)
		}
	})

	t.Run("update order", func(t *testing.T) {
		err := database.UpdateOrder(ctx, db.UpdateOrderArgs{
			OrderID:     order.OrderID,
			ClientID:    order.ClientID,
			TransportID: order.TransportID,
			EmployeeID:  order.EmployeeID,
			Grade:       order.Grade,
			Distance:    order.Distance,
			Weight:      order.Weight,
			TotalPrice:  decimal.NewFromInt(6000),
			PriceID:     order.PriceID,
			Status:      models.OrderStatus("in_progress"),
			NodeIDStart: order.NodeIDStart,
			NodeIDEnd:   order.NodeIDEnd,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := database.GetOrderByID(ctx, order.OrderID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Status != models.OrderStatus("in_progress") {
			t.Fatalf("expected updated status, got %v", got.Status)
		}
	})

	t.Run("update order status", func(t *testing.T) {
		err := database.UpdateOrderStatus(ctx, order.OrderID, models.OrderStatus("completed"))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		got, err := database.GetOrderByID(ctx, order.OrderID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Status != models.OrderStatus("completed") {
			t.Fatalf("expected completed status, got %v", got.Status)
		}
	})

	t.Run("soft delete and restore order", func(t *testing.T) {
		if err := database.SoftDeleteOrder(ctx, order.OrderID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
		if err := database.RestoreOrder(ctx, order.OrderID); err != nil {
			t.Fatalf("restore: %v", err)
		}

		if _, err := database.GetOrderByID(ctx, order.OrderID); err != nil {
			t.Fatalf("get after restore: %v", err)
		}
	})

	t.Run("bulk delete orders", func(t *testing.T) {
		order2 := createTestOrder(t, database)

		if err := database.BulkSoftDeleteOrders(ctx, []int32{order2.OrderID}); err != nil {
			t.Fatalf("bulk soft delete: %v", err)
		}
		if err := database.BulkHardDeleteOrders(ctx, []int32{order2.OrderID}); err != nil {
			t.Fatalf("bulk hard delete: %v", err)
		}
	})

	t.Run("hard delete order", func(t *testing.T) {
		if err := database.HardDeleteOrder(ctx, order.OrderID); err != nil {
			t.Fatalf("hard delete: %v", err)
		}
		if _, err := database.GetOrderByID(ctx, order.OrderID); err == nil {
			t.Fatalf("expected error after hard delete")
		}
	})
}

var orderTestC int64

func createTestOrder(t *testing.T, database db.DB) models.Order {
	t.Helper()
	ctx := context.Background()

	orderTestC++

	email := "order-client@test.test" + strconv.FormatInt(orderTestC, 10)
	phone := "+37529389565" + strconv.FormatInt(orderTestC, 10)
	client := createTestClient(t, database, email, phone)

	plate := "ORD" + strconv.FormatInt(orderTestC, 10)
	transport := createTestTransport(t, database, plate)

	employee := createTestEmployee(t, database, "Order Employee")
	price := createTestPrice(t, database, "ORD-CARGO"+strconv.FormatInt(orderTestC, 10))

	startID := "ORD-NODE-START" + strconv.FormatInt(orderTestC, 10)
	start := createTestNode(t, database, startID, float64(10+orderTestC), 10)

	endID := "ORD-NODE-END" + strconv.FormatInt(orderTestC, 10)
	end := createTestNode(t, database, endID, float64(20+orderTestC), 20)

	order, err := database.CreateOrder(ctx, db.CreateOrderArg{
		ClientID:    client.ClientID,
		TransportID: transport.TransportID,
		EmployeeID:  employee.EmployeeID,
		Grade:       1,
		Distance:    100,
		Weight:      1000,
		TotalPrice:  decimal.NewFromInt(5000),
		PriceID:     price.PriceID,
		Status:      models.OrderStatus("pending"),
		NodeIDStart: start.NodeID,
		NodeIDEnd:   end.NodeID,
	})
	if err != nil {
		t.Fatalf("create order: %v", err)
	}
	return order
}
