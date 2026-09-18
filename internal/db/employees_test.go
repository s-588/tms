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
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

func TestDB_Employees(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testCreateEmployee(t, database)
	testGetEmployeeByID(t, database)
	testGetEmployees(t, database)
	testUpdateEmployee(t, database)
	testSoftDeleteEmployee(t, database)
	testHardDeleteEmployee(t, database)
	testRestoreEmployee(t, database)
	testBulkSoftDeleteEmployees(t, database)
	testBulkHardDeleteEmployees(t, database)
	testListFreeDrivers(t, database)
}

func testCreateEmployee(t *testing.T, database db.DB) {
	ctx := context.Background()
	now := time.Now()

	t.Run("creates employee successfully", func(t *testing.T) {
		args := db.CreateEmployeeArgs{
			Name:              "Test Driver",
			Status:            models.EmployeeStatusAvailable,
			JobTitle:          models.EmployeeJobTitleDriver,
			HireDate:          now,
			Salary:            decimal.NewFromInt(1000),
			LicenseIssued:     now,
			LicenseExpiration: now.AddDate(1, 0, 0),
		}
		employee, err := database.CreateEmployee(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if employee.EmployeeID == 0 {
			t.Errorf("expected employee ID to be populated")
		}
	})
}

func testGetEmployeeByID(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Get Employee")

	t.Run("get employee by id", func(t *testing.T) {
		got, err := database.GetEmployeeByID(ctx, employee.EmployeeID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, employee) {
			t.Fatalf("got %v, want %v", got, employee)
		}
	})
}

func testGetEmployees(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Get All Employees")

	t.Run("get all employees", func(t *testing.T) {
		list, total, err := database.GetEmployees(ctx, 1, models.EmployeeFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total pages to be populated")
		}
		if !slices.ContainsFunc(list, func(v models.Employee) bool {
			return reflect.DeepEqual(v, employee)
		}) {
			t.Fatalf("slice doesn't contain wanted employee; got %v, want %v", list, employee)
		}
	})
}

func testUpdateEmployee(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Update Employee")
	now := time.Now()

	t.Run("update employee", func(t *testing.T) {
		err := database.UpdateEmployee(ctx, db.UpdateEmployeeArgs{
			EmployeeID:        employee.EmployeeID,
			Name:              "Updated Driver",
			Status:            models.EmployeeStatusAssigned,
			JobTitle:          models.EmployeeJobTitleDriver,
			HireDate:          now,
			Salary:            decimal.NewFromInt(1200),
			LicenseIssued:     now,
			LicenseExpiration: now.AddDate(2, 0, 0),
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := database.GetEmployeeByID(ctx, employee.EmployeeID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Name != "Updated Driver" {
			t.Fatalf("expected updated name, got %s", got.Name)
		}
	})
}

func testSoftDeleteEmployee(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Soft Delete Employee")

	t.Run("soft delete employee", func(t *testing.T) {
		if err := database.SoftDeleteEmployee(ctx, employee.EmployeeID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
	})
}

func testHardDeleteEmployee(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Hard Delete Employee")

	t.Run("hard delete employee", func(t *testing.T) {
		if err := database.HardDeleteEmployee(ctx, employee.EmployeeID); err != nil {
			t.Fatalf("hard delete: %v", err)
		}
		if _, err := database.GetEmployeeByID(ctx, employee.EmployeeID); err == nil {
			t.Fatalf("expected error after hard delete")
		}
	})
}

func testRestoreEmployee(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Restore Employee")

	t.Run("restore employee", func(t *testing.T) {
		if err := database.SoftDeleteEmployee(ctx, employee.EmployeeID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
		if err := database.RestoreEmployee(ctx, employee.EmployeeID); err != nil {
			t.Fatalf("restore: %v", err)
		}
		if _, err := database.GetEmployeeByID(ctx, employee.EmployeeID); err != nil {
			t.Fatalf("get after restore: %v", err)
		}
	})
}

func testBulkSoftDeleteEmployees(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Bulk Soft Employee")

	t.Run("bulk soft delete employees", func(t *testing.T) {
		if err := database.BulkSoftDeleteEmployees(ctx, []int32{employee.EmployeeID}); err != nil {
			t.Fatalf("bulk soft delete: %v", err)
		}
	})
}

func testBulkHardDeleteEmployees(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Bulk Hard Employee")

	t.Run("bulk hard delete employees", func(t *testing.T) {
		if err := database.BulkHardDeleteEmployees(ctx, []int32{employee.EmployeeID}); err != nil {
			t.Fatalf("bulk hard delete: %v", err)
		}
		if _, err := database.GetEmployeeByID(ctx, employee.EmployeeID); err == nil {
			t.Fatalf("expected error after bulk hard delete")
		}
	})
}

func testListFreeDrivers(t *testing.T, database db.DB) {
	ctx := context.Background()
	employee := createTestEmployee(t, database, "Free Driver")

	t.Run("list free drivers", func(t *testing.T) {
		items, err := database.ListFreeDrivers(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.ContainsFunc(items, func(item ui.ListItem) bool {
			return item.ID == employee.EmployeeID
		}) {
			t.Fatalf("list doesn't contain free driver; got %v", items)
		}
	})
}
func createTestEmployee(t *testing.T, database db.DB, name string) models.Employee {
	t.Helper()
	ctx := context.Background()
	now := time.Now()

	employee, err := database.CreateEmployee(ctx, db.CreateEmployeeArgs{
		Name:              name,
		Status:            models.EmployeeStatusAvailable,
		JobTitle:          models.EmployeeJobTitleDriver,
		HireDate:          now,
		Salary:            decimal.NewFromInt(1000),
		LicenseIssued:     now,
		LicenseExpiration: now.AddDate(1, 0, 0),
	})
	if err != nil {
		t.Fatalf("create employee %q: %v", name, err)
	}
	return employee
}
