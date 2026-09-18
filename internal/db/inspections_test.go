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
)

func TestDB_Inspections(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testCreateInspection(t, database)
	testGetInspectionByID(t, database)
	testGetInspectionsByTransport(t, database)
	testGetInspections(t, database)
	testUpdateInspection(t, database)
	testSoftDeleteInspection(t, database)
	testHardDeleteInspection(t, database)
	testRestoreInspection(t, database)
	testBulkSoftDeleteInspections(t, database)
	testBulkHardDeleteInspections(t, database)
}

func testCreateInspection(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "INSP-001")
	now := time.Now()

	t.Run("creates inspection successfully", func(t *testing.T) {
		args := db.CreateInspectionArgs{
			TransportID:          transport.TransportID,
			InspectionDate:       now,
			InspectionExpiration: now.AddDate(1, 0, 0),
			Status:               models.InspectionStatusReady,
		}
		inspection, err := database.CreateInspection(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if inspection.InspectionID == 0 {
			t.Errorf("expected inspection ID to be populated")
		}
	})
}

func testGetInspectionByID(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "GET-ID")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("get inspection by id", func(t *testing.T) {
		got, err := database.GetInspectionByID(ctx, inspection.InspectionID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, inspection) {
			t.Fatalf("got %v, want %v", got, inspection)
		}
	})
}

func testGetInspectionsByTransport(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "GET-ALL")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("get inspections by transport", func(t *testing.T) {
		list, err := database.GetInspectionsByTransport(ctx, transport.TransportID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.ContainsFunc(list, func(v models.Inspection) bool {
			return reflect.DeepEqual(v, inspection)
		}) {
			t.Fatalf("slice doesn't contain wanted inspection; got %v, want %v", list, inspection)
		}
	})
}

func testGetInspections(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-GET")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("get all inspections", func(t *testing.T) {
		list, total, err := database.GetInspections(ctx, 1, models.InspectionFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total pages to be populated")
		}
		if !slices.ContainsFunc(list, func(v models.Inspection) bool {
			return reflect.DeepEqual(v, inspection)
		}) {
			t.Fatalf("slice doesn't contain wanted inspection; got %v, want %v", list, inspection)
		}
	})
}

func testUpdateInspection(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-UPD")
	inspection := createTestInspection(t, database, transport.TransportID)
	now := time.Now()

	t.Run("update inspection", func(t *testing.T) {
		err := database.UpdateInspection(ctx, db.UpdateInspectionArgs{
			InspectionID:         inspection.InspectionID,
			TransportID:          transport.TransportID,
			InspectionDate:       now,
			InspectionExpiration: now.AddDate(2, 0, 0),
			Status:               models.InspectionStatusOverdue,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := database.GetInspectionByID(ctx, inspection.InspectionID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Status != models.InspectionStatusOverdue {
			t.Fatalf("expected updated status, got %v", got.Status)
		}
	})
}

func testSoftDeleteInspection(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-SOFT")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("soft delete inspection", func(t *testing.T) {
		if err := database.SoftDeleteInspection(ctx, inspection.InspectionID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
	})
}

func testHardDeleteInspection(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-HARD")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("hard delete inspection", func(t *testing.T) {
		if err := database.HardDeleteInspection(ctx, inspection.InspectionID); err != nil {
			t.Fatalf("hard delete: %v", err)
		}
		if _, err := database.GetInspectionByID(ctx, inspection.InspectionID); err == nil {
			t.Fatalf("expected error after hard delete")
		}
	})
}

func testRestoreInspection(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-REST")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("restore inspection", func(t *testing.T) {
		if err := database.SoftDeleteInspection(ctx, inspection.InspectionID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
		if err := database.RestoreInspection(ctx, inspection.InspectionID); err != nil {
			t.Fatalf("restore: %v", err)
		}
		if _, err := database.GetInspectionByID(ctx, inspection.InspectionID); err != nil {
			t.Fatalf("get after restore: %v", err)
		}
	})
}

func testBulkSoftDeleteInspections(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-BSOFT")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("bulk soft delete inspections", func(t *testing.T) {
		if err := database.BulkSoftDeleteInspections(ctx, []int32{inspection.InspectionID}); err != nil {
			t.Fatalf("bulk soft delete: %v", err)
		}
	})
}

func testBulkHardDeleteInspections(t *testing.T, database db.DB) {
	ctx := context.Background()
	transport := createTestTransport(t, database, "TR-BHARD")
	inspection := createTestInspection(t, database, transport.TransportID)

	t.Run("bulk hard delete inspections", func(t *testing.T) {
		if err := database.BulkHardDeleteInspections(ctx, []int32{inspection.InspectionID}); err != nil {
			t.Fatalf("bulk hard delete: %v", err)
		}
		if _, err := database.GetInspectionByID(ctx, inspection.InspectionID); err == nil {
			t.Fatalf("expected error after bulk hard delete")
		}
	})
}

func createTestInspection(t *testing.T, database db.DB, transportID int32) models.Inspection {
	t.Helper()
	ctx := context.Background()
	now := time.Now()
	inspection, err := database.CreateInspection(ctx, db.CreateInspectionArgs{
		TransportID:          transportID,
		InspectionDate:       now,
		InspectionExpiration: now.AddDate(1, 0, 0),
		Status:               models.InspectionStatusReady,
	})
	if err != nil {
		t.Fatalf("create inspection: %v", err)
	}
	return inspection
}
