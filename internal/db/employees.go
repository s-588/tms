package db

import (
	"context"
	"strings"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

type CreateEmployeeArgs struct {
	Name              string
	Status            models.EmployeeStatus
	JobTitle          models.EmployeeJobTitle
	HireDate          time.Time
	Salary            decimal.Decimal
	LicenseIssued     time.Time
	LicenseExpiration time.Time
}

func (db DB) CreateEmployee(ctx context.Context, args CreateEmployeeArgs) (models.Employee, error) {
	arg := generated.CreateEmployeeParams{
		Name:              args.Name,
		Status:            generated.EmployeeStatus(args.Status),
		JobTitle:          generated.EmployeeJobTitle(args.JobTitle),
		HireDate:          ToPgTypeDateFromTime(args.HireDate),
		Salary:            args.Salary,
		LicenseIssued:     ToPgTypeDateFromTime(args.LicenseIssued),
		LicenseExpiration: ToPgTypeDateFromTime(args.LicenseExpiration),
	}
	genEmployee, err := db.queries.CreateEmployee(ctx, arg)
	if err != nil {
		return models.Employee{}, err
	}
	return convertGeneratedEmployeeToModel(genEmployee), nil
}

func (db DB) GetEmployeeByID(ctx context.Context, employeeID int32) (models.Employee, error) {
	genEmployee, err := db.queries.GetEmployee(ctx, employeeID)
	if err != nil {
		return models.Employee{}, err
	}
	return convertGeneratedEmployeeToModel(genEmployee), nil
}

func (db DB) GetEmployees(ctx context.Context, page int32, filter models.EmployeeFilter) ([]models.Employee, int32, error) {
	arg := generated.GetEmployeesParams{
		Page:           page,
		NameFilter:     ToStringPtr(filter.Name),
		JobTitleFilter: ToNullEmployeeJobTitle(filter.JobTitle),
		StatusFilter:   ToNullEmployeeStatus(filter.Status),
		SalaryMin:      ToPgTypeNumeric(filter.SalaryMin),
		SalaryMax:      ToPgTypeNumeric(filter.SalaryMax),
		SortBy:         filter.GetSortBy(),
		SortOrder:      filter.GetSortOrder(),
	}
	rows, err := db.queries.GetEmployees(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var employees []models.Employee
	var totalPages int32
	for _, row := range rows {
		totalPages = row.TotalCount
		employees = append(employees, convertGeneratedEmployeeRowToModel(row))
	}
	return employees, totalPages, nil
}

type UpdateEmployeeArgs struct {
	EmployeeID        int32
	Name              string
	Status            models.EmployeeStatus
	JobTitle          models.EmployeeJobTitle
	HireDate          time.Time
	Salary            decimal.Decimal
	LicenseIssued     time.Time
	LicenseExpiration time.Time
}

func (db DB) UpdateEmployee(ctx context.Context, args UpdateEmployeeArgs) error {
	arg := generated.UpdateEmployeeParams{
		EmployeeID:        args.EmployeeID,
		Name:              args.Name,
		Status:            generated.EmployeeStatus(args.Status),
		JobTitle:          generated.EmployeeJobTitle(args.JobTitle),
		HireDate:          ToPgTypeDateFromTime(args.HireDate),
		Salary:            args.Salary,
		LicenseIssued:     ToPgTypeDateFromTime(args.LicenseIssued),
		LicenseExpiration: ToPgTypeDateFromTime(args.LicenseExpiration),
	}
	return db.queries.UpdateEmployee(ctx, arg)
}

func (db DB) SoftDeleteEmployee(ctx context.Context, employeeID int32) error {
	return db.queries.SoftDeleteEmployee(ctx, employeeID)
}

func (db DB) HardDeleteEmployee(ctx context.Context, employeeID int32) error {
	return db.queries.HardDeleteEmployee(ctx, employeeID)
}

func (db DB) RestoreEmployee(ctx context.Context, employeeID int32) error {
	return db.queries.RestoreEmployee(ctx, employeeID)
}

func (db DB) BulkSoftDeleteEmployees(ctx context.Context, employeeIDs []int32) error {
	return db.queries.BulkSoftDeleteEmployees(ctx, employeeIDs)
}

func (db DB) BulkHardDeleteEmployees(ctx context.Context, employeeIDs []int32) error {
	return db.queries.BulkHardDeleteEmployees(ctx, employeeIDs)
}

// conversion helpers
func convertGeneratedEmployeeToModel(e generated.Employee) models.Employee {
	return models.Employee{
		EmployeeID:        e.EmployeeID,
		Name:              e.Name,
		Status:            models.EmployeeStatus(e.Status),
		JobTitle:          models.EmployeeJobTitle(e.JobTitle),
		HireDate:          fromPgDate(e.HireDate),
		Salary:            e.Salary,
		LicenseIssued:     fromPgDate(e.LicenseIssued),
		LicenseExpiration: fromPgDate(e.LicenseExpiration),
		CreatedAt:         fromPgTimestamptz(e.CreatedAt),
		UpdatedAt:         fromPgTimestamptz(e.UpdatedAt),
		DeletedAt:         fromPgTimestamptz(e.DeletedAt),
	}
}

func convertGeneratedEmployeeRowToModel(row generated.GetEmployeesRow) models.Employee {
	return models.Employee{
		EmployeeID:        row.EmployeeID,
		Name:              row.Name,
		Status:            models.EmployeeStatus(row.Status),
		JobTitle:          models.EmployeeJobTitle(row.JobTitle),
		HireDate:          fromPgDate(row.HireDate),
		Salary:            row.Salary,
		LicenseIssued:     fromPgDate(row.LicenseIssued),
		LicenseExpiration: fromPgDate(row.LicenseExpiration),
		CreatedAt:         fromPgTimestamptz(row.CreatedAt),
		UpdatedAt:         fromPgTimestamptz(row.UpdatedAt),
		DeletedAt:         fromPgTimestamptz(row.DeletedAt),
	}
}

func (db DB) ListFreeDrivers(ctx context.Context) ([]ui.ListItem, error) {
	rows, err := db.queries.ListFreeDrivers(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ui.ListItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, ui.ListItem{
			ID:   r.EmployeeID,
			Name: strings.Join([]string{r.Name, string(r.Status), string(r.JobTitle)}, " "),
		})
	}
	return items, nil
}
