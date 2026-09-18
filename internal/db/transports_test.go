package db_test

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/testutil"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

func TestDB_Transports(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testCreateTransport(t, database)
	testGetTransportByID(t, database)
	testGetTransports(t, database)
	testUpdateTransport(t, database)
	testSoftDeleteTransport(t, database)
	testHardDeleteTransport(t, database)
	testRestoreTransport(t, database)
	testBulkSoftDeleteTransports(t, database)
	testBulkHardDeleteTransports(t, database)
	testGetTransportOrders(t, database)
	testListFreeTransports(t, database)
}

func testCreateTransport(t *testing.T, database db.DB) {
	ctx := context.Background()

	t.Run("creates transport successfully", func(t *testing.T) {
		args := db.CreateTransportArgs{
			Model:           "Volvo FH",
			LicensePlate:    "TR-001",
			PayloadCapacity: 10000,
			FuelConsumption: 20,
		}
		transport, err := database.CreateTransport(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if transport.TransportID == 0 {
			t.Errorf("expected transport ID to be populated")
		}
		if transport.LicensePlate != args.LicensePlate {
			t.Errorf("expected license plate %s, got %s", args.LicensePlate, transport.LicensePlate)
		}
	})

	t.Run("error; duplicate license plate", func(t *testing.T) {
		args := db.CreateTransportArgs{
			Model:           "Scania R",
			LicensePlate:    "TR-DUP",
			PayloadCapacity: 9000,
			FuelConsumption: 19,
		}
		_, err := database.CreateTransport(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		_, err = database.CreateTransport(ctx, args)
		if !errors.Is(err, db.ErrDuplicateLicense) {
			t.Fatalf("expected %v, got %v", db.ErrDuplicateLicense, err)
		}
	})
}

func testGetTransportByID(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-GETID")

	t.Run("get transport by id", func(t *testing.T) {
		got, err := database.GetTransportByID(ctx, transport.TransportID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, transport) {
			t.Fatalf("got %v, want %v", got, transport)
		}
	})
}

func testGetTransports(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-GETA")

	t.Run("get all transports", func(t *testing.T) {
		list, total, err := database.GetTransports(ctx, 1, models.TransportFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total pages to be populated")
		}
		if !slices.ContainsFunc(list, func(v models.Transport) bool {
			return reflect.DeepEqual(v, transport)
		}) {
			t.Fatalf("slice doesn't contain wanted transport; got %v, want %v", list, transport)
		}
	})
}

func testUpdateTransport(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-UPD")

	t.Run("update transport", func(t *testing.T) {
		err := database.UpdateTransport(ctx, db.UpdateTransportArgs{
			TransportID:     transport.TransportID,
			Model:           "Scania R",
			LicensePlate:    transport.LicensePlate,
			PayloadCapacity: 12000,
			FuelConsumption: 22,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := database.GetTransportByID(ctx, transport.TransportID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Model != "Scania R" {
			t.Fatalf("expected updated model, got %s", got.Model)
		}
	})
}

func testSoftDeleteTransport(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-SOFT")

	t.Run("soft delete transport", func(t *testing.T) {
		if err := database.SoftDeleteTransport(ctx, transport.TransportID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
	})
}

func testHardDeleteTransport(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-HARD")

	t.Run("hard delete transport", func(t *testing.T) {
		if err := database.HardDeleteTransport(ctx, transport.TransportID); err != nil {
			t.Fatalf("hard delete: %v", err)
		}
		if _, err := database.GetTransportByID(ctx, transport.TransportID); err == nil {
			t.Fatalf("expected error after hard delete")
		}
	})
}

func testRestoreTransport(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-REST")

	t.Run("restore transport", func(t *testing.T) {
		if err := database.SoftDeleteTransport(ctx, transport.TransportID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
		if err := database.RestoreTransport(ctx, transport.TransportID); err != nil {
			t.Fatalf("restore: %v", err)
		}
		if _, err := database.GetTransportByID(ctx, transport.TransportID); err != nil {
			t.Fatalf("get after restore: %v", err)
		}
	})
}

func testBulkSoftDeleteTransports(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-BSOFT")

	t.Run("bulk soft delete transports", func(t *testing.T) {
		if err := database.BulkSoftDeleteTransports(ctx, []int32{transport.TransportID}); err != nil {
			t.Fatalf("bulk soft delete: %v", err)
		}
	})
}

func testBulkHardDeleteTransports(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-BHARD")

	t.Run("bulk hard delete transports", func(t *testing.T) {
		if err := database.BulkHardDeleteTransports(ctx, []int32{transport.TransportID}); err != nil {
			t.Fatalf("bulk hard delete: %v", err)
		}
		if _, err := database.GetTransportByID(ctx, transport.TransportID); err == nil {
			t.Fatalf("expected error after bulk hard delete")
		}
	})
}

func testGetTransportOrders(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-ORDERS")

	client := createTestClient(t, database, "transport-orders@test.test", "+19990000002")
	employee := createTestEmployee(t, database, "Transport Orders Employee")
	price := createTestPrice(t, database, "TR-ORD-CARGO")
	start := createTestNode(t, database, "TR-GET-ORD-START", 100, 100)
	end := createTestNode(t, database, "TR-GET-ORD-END", 200, 200)

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
		t.Fatalf("create order for transport: %v", err)
	}

	t.Run("get transport orders", func(t *testing.T) {
		orders, total, err := database.GetTransportOrders(ctx, transport.TransportID, 10, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total count to be populated")
		}
		if !slices.ContainsFunc(orders, func(o models.Order) bool {
			return o.OrderID == order.OrderID
		}) {
			t.Fatalf("orders slice doesn't contain wanted order; got %v, want %v", orders, order)
		}
	})
}

func testListFreeTransports(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-FREE")
	createTestInsurance(t, database, transport.TransportID)
	t.Run("list free transports", func(t *testing.T) {
		items, err := database.ListFreeTransports(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.ContainsFunc(items, func(item ui.ListItem) bool {
			return item.ID == transport.TransportID
		}) {
			t.Fatalf("list doesn't contain free transport; got %v, want to contain %v", items, transport)
		}
	})
}

func createTestTransport(t *testing.T, database db.DB, plate string) models.Transport {
	t.Helper()
	ctx := context.Background()

	transport, err := database.CreateTransport(ctx, db.CreateTransportArgs{
		Model:           "Test Truck",
		LicensePlate:    plate,
		PayloadCapacity: 10000,
		FuelConsumption: 20,
	})
	if err != nil {
		t.Fatalf("create transport %q: %v", plate, err)
	}
	return transport
}
