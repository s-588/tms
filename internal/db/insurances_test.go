package db_test

import (
	"context"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/testutil"
	"github.com/shopspring/decimal"
)

func TestDB_Insurances(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testCreateInsurance(t, database)
	testGetInsuranceByID(t, database)
	testGetInsuranceByTransport(t, database)
	testGetInsurances(t, database)
	testUpdateInsurance(t, database)
	testSoftDeleteInsurance(t, database)
	testHardDeleteInsurance(t, database)
	testRestoreInsurance(t, database)
	testBulkSoftDeleteInsurances(t, database)
	testBulkHardDeleteInsurances(t, database)
}

func testCreateInsurance(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-001")
	now := time.Now()

	t.Run("creates insurance successfully", func(t *testing.T) {
		args := db.CreateInsuranceArgs{
			TransportID:         transport.TransportID,
			InsuranceDate:       now,
			InsuranceExpiration: now.AddDate(1, 0, 0),
			Payment:             decimal.NewFromInt(500),
			Coverage:            decimal.NewFromInt(10000),
		}
		insurance, err := database.CreateInsurance(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if insurance.InsuranceID == 0 {
			t.Errorf("expected insurance ID to be populated")
		}
	})
}

func testGetInsuranceByID(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INS-GETID")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("get insurance by id", func(t *testing.T) {
		got, err := database.GetInsuranceByID(ctx, insurance.InsuranceID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, insurance) {
			t.Fatalf("got %v, want %v", got, insurance)
		}
	})
}

func testGetInsuranceByTransport(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INS-GETTR")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("get insurance by transport", func(t *testing.T) {
		got, err := database.GetInsuranceByTransport(ctx, transport.TransportID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, insurance) {
			t.Fatalf("got %v, want %v", got, insurance)
		}
	})
}

func testGetInsurances(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INS-GALL")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("get all insurances", func(t *testing.T) {
		list, total, err := database.GetInsurances(ctx, 1, models.InsuranceFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total pages to be populated")
		}
		if !slices.ContainsFunc(list, func(v models.Insurance) bool {
			return reflect.DeepEqual(v, insurance)
		}) {
			t.Fatalf("slice doesn't contain wanted insurance; got %v, want %v", list, insurance)
		}
	})
}

func testUpdateInsurance(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-UPD")
	insurance := createTestInsurance(t, database, transport.TransportID)
	now := time.Now()

	t.Run("update insurance", func(t *testing.T) {
		err := database.UpdateInsurance(ctx, db.UpdateInsuranceArgs{
			InsuranceID:         insurance.InsuranceID,
			TransportID:         transport.TransportID,
			InsuranceDate:       now,
			InsuranceExpiration: now.AddDate(2, 0, 0),
			Payment:             decimal.NewFromInt(700),
			Coverage:            decimal.NewFromInt(12000),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := database.GetInsuranceByID(ctx, insurance.InsuranceID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !got.Payment.Equal(decimal.NewFromInt(700)) {
			t.Fatalf("expected updated payment, got %v", got.Payment)
		}
	})
}

func testSoftDeleteInsurance(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INS-SOFT")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("soft delete insurance", func(t *testing.T) {
		if err := database.SoftDeleteInsurance(ctx, insurance.InsuranceID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
	})
}

func testHardDeleteInsurance(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INS-HARD")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("hard delete insurance", func(t *testing.T) {
		if err := database.HardDeleteInsurance(ctx, insurance.InsuranceID); err != nil {
			t.Fatalf("hard delete: %v", err)
		}
		if _, err := database.GetInsuranceByID(ctx, insurance.InsuranceID); err == nil {
			t.Fatalf("expected error after hard delete")
		}
	})
}

func testRestoreInsurance(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INS-REST")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("restore insurance", func(t *testing.T) {
		if err := database.SoftDeleteInsurance(ctx, insurance.InsuranceID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
		if err := database.RestoreInsurance(ctx, insurance.InsuranceID); err != nil {
			t.Fatalf("restore: %v", err)
		}
		if _, err := database.GetInsuranceByID(ctx, insurance.InsuranceID); err != nil {
			t.Fatalf("get after restore: %v", err)
		}
	})
}

func testBulkSoftDeleteInsurances(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INSBSOFT")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("bulk soft delete insurances", func(t *testing.T) {
		if err := database.BulkSoftDeleteInsurances(ctx, []int32{insurance.InsuranceID}); err != nil {
			t.Fatalf("bulk soft delete: %v", err)
		}
	})
}

func testBulkHardDeleteInsurances(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INSBHARD")
	insurance := createTestInsurance(t, database, transport.TransportID)

	t.Run("bulk hard delete insurances", func(t *testing.T) {
		if err := database.BulkHardDeleteInsurances(ctx, []int32{insurance.InsuranceID}); err != nil {
			t.Fatalf("bulk hard delete: %v", err)
		}
		if _, err := database.GetInsuranceByID(ctx, insurance.InsuranceID); err == nil {
			t.Fatalf("expected error after bulk hard delete")
		}
	})
}

func createTestInsurance(t *testing.T, database db.DB, transportID int32) models.Insurance {
	t.Helper()
	ctx := context.Background()
	now := time.Now()
	insurance, err := database.CreateInsurance(ctx, db.CreateInsuranceArgs{
		TransportID:         transportID,
		InsuranceDate:       now,
		InsuranceExpiration: now.AddDate(1, 0, 0),
		Payment:             decimal.NewFromInt(500),
		Coverage:            decimal.NewFromInt(10000),
	})
	if err != nil {
		t.Fatalf("create insurance: %v", err)
	}
	return insurance
}
