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

func TestDB_Prices(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testCreatePrice(t, database)
	testGetPriceByID(t, database)
	testGetPriceByUnique(t, database)
	testGetPrices(t, database)
	testUpdatePrice(t, database)
	testSoftDeletePrice(t, database)
	testHardDeletePrice(t, database)
	testRestorePrice(t, database)
	testBulkSoftDeletePrices(t, database)
	testBulkHardDeletePrices(t, database)
	testListPrices(t, database)
}

func testCreatePrice(t *testing.T, database db.DB) {
	ctx := context.Background()

	t.Run("creates price successfully", func(t *testing.T) {
		args := db.CreatePriceArgs{
			CargoType: "general",
			Weight:    decimal.NewFromInt(100),
			Distance:  decimal.NewFromInt(50),
		}
		price, err := database.CreatePrice(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if price.PriceID == 0 {
			t.Errorf("expected price ID to be populated")
		}
	})

	t.Run("error; duplicate price", func(t *testing.T) {
		args := db.CreatePriceArgs{
			CargoType: "dup",
			Weight:    decimal.NewFromInt(200),
			Distance:  decimal.NewFromInt(100),
		}
		_, err := database.CreatePrice(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		_, err = database.CreatePrice(ctx, args)
		if !errors.Is(err, db.ErrDuplicatePrice) {
			t.Fatalf("expected %v, got %v", db.ErrDuplicatePrice, err)
		}
	})
}

func testGetPriceByID(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "GET-ID")

	t.Run("get price by id", func(t *testing.T) {
		got, err := database.GetPriceByID(ctx, price.PriceID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, price) {
			t.Fatalf("got %v, want %v", got, price)
		}
	})
}

func testGetPriceByUnique(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "UNIQUE")

	t.Run("get price by unique", func(t *testing.T) {
		got, err := database.GetPriceByUnique(ctx, price.CargoType, price.Weight, price.Distance)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, price) {
			t.Fatalf("got %v, want %v", got, price)
		}
	})
}

func testGetPrices(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "GET-ALL")

	t.Run("get all prices", func(t *testing.T) {
		list, total, err := database.GetPrices(ctx, 1, models.PriceFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total pages to be populated")
		}
		if !slices.ContainsFunc(list, func(v models.Price) bool {
			return reflect.DeepEqual(v, price)
		}) {
			t.Fatalf("slice doesn't contain wanted price; got %v, want %v", list, price)
		}
	})
}

func testUpdatePrice(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "UPD")

	t.Run("update price", func(t *testing.T) {
		err := database.UpdatePrice(ctx, db.UpdatePriceArgs{
			PriceID:   price.PriceID,
			CargoType: "updated",
			Weight:    decimal.NewFromInt(200),
			Distance:  decimal.NewFromInt(100),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := database.GetPriceByID(ctx, price.PriceID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.CargoType != "updated" {
			t.Fatalf("expected updated cargo type, got %s", got.CargoType)
		}
	})
}

func testSoftDeletePrice(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "SOFT")

	t.Run("soft delete price", func(t *testing.T) {
		if err := database.SoftDeletePrice(ctx, price.PriceID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
	})
}

func testHardDeletePrice(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "HARD")

	t.Run("hard delete price", func(t *testing.T) {
		if err := database.HardDeletePrice(ctx, price.PriceID); err != nil {
			t.Fatalf("hard delete: %v", err)
		}
		if _, err := database.GetPriceByID(ctx, price.PriceID); err == nil {
			t.Fatalf("expected error after hard delete")
		}
	})
}

func testRestorePrice(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "REST")

	t.Run("restore price", func(t *testing.T) {
		if err := database.SoftDeletePrice(ctx, price.PriceID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
		if err := database.RestorePrice(ctx, price.PriceID); err != nil {
			t.Fatalf("restore: %v", err)
		}
		if _, err := database.GetPriceByID(ctx, price.PriceID); err != nil {
			t.Fatalf("get after restore: %v", err)
		}
	})
}

func testBulkSoftDeletePrices(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "BULK-SOFT")

	t.Run("bulk soft delete prices", func(t *testing.T) {
		if err := database.BulkSoftDeletePrices(ctx, []int32{price.PriceID}); err != nil {
			t.Fatalf("bulk soft delete: %v", err)
		}
	})
}

func testBulkHardDeletePrices(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "BULK-HARD")

	t.Run("bulk hard delete prices", func(t *testing.T) {
		if err := database.BulkHardDeletePrices(ctx, []int32{price.PriceID}); err != nil {
			t.Fatalf("bulk hard delete: %v", err)
		}
		if _, err := database.GetPriceByID(ctx, price.PriceID); err == nil {
			t.Fatalf("expected error after bulk hard delete")
		}
	})
}

func testListPrices(t *testing.T, database db.DB) {
	ctx := context.Background()
	price := createTestPrice(t, database, "LIST")

	t.Run("list prices", func(t *testing.T) {
		items, err := database.ListPrices(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.ContainsFunc(items, func(item ui.ListItem) bool {
			return item.ID == price.PriceID
		}) {
			t.Fatalf("list doesn't contain price; got %v", items)
		}
	})
}

func createTestPrice(t *testing.T, database db.DB, cargoType string) models.Price {
	t.Helper()
	ctx := context.Background()

	price, err := database.CreatePrice(ctx, db.CreatePriceArgs{
		CargoType: cargoType,
		Weight:    decimal.NewFromInt(100),
		Distance:  decimal.NewFromInt(100),
	})
	if err != nil {
		t.Fatalf("create price %q: %v", cargoType, err)
	}
	return price
}
