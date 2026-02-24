package db

import (
	"context"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

func (db DB) CreateEmployee(ctx context.Context, emp models.Employee) (models.Employee, error) {
	arg := generated.CreateEmployeeParams{
		Name:              emp.Name,
		Status:            generated.EmployeeStatus(emp.Status),
		JobTitle:          string(emp.JobTitle),
		HireDate:          ToPgTypeDateFromTime(emp.HireDate),
		Salary:            emp.Salary,
		LicenseIssued:     ToPgTypeDateFromTime(emp.LicenseIssued),
		LicenseExpiration: ToPgTypeDateFromTime(emp.LicenseExpiration),
	}
	genEmployee, err := db.queries.CreateEmployee(ctx, arg)
	if err != nil {
		return models.Employee{}, err
	}
	return convertGeneratedEmployeeToModel(genEmployee), nil
}

func (db DB) GetEmployeeByID(ctx context.Context, employeeID int) (models.Employee, error) {
	genEmployee, err := db.queries.GetEmployee(ctx, int32(employeeID))
	if err != nil {
		return models.Employee{}, err
	}
	return convertGeneratedEmployeeToModel(genEmployee), nil
}

func (db DB) GetEmployees(ctx context.Context, limit, offset int, filter models.EmployeeFilter) ([]models.Employee, int64, error) {
	arg := generated.GetEmployeesParams{
		Limit:          int32(limit),
		Offset:         int32(offset),
		NameFilter:     ToStringPtr(filter.Name),
		JobTitleFilter: ToStringPtr(filter.JobTitle),
		StatusFilter:   ToNullEmployeeStatus(filter.Status),
		SalaryMin:      ToPgTypeNumeric(filter.SalaryMin),
		SalaryMax:      ToPgTypeNumeric(filter.SalaryMax),
		CreatedFrom:    ToPgTypeTimestamptz(filter.CreatedFrom),
		CreatedTo:      ToPgTypeTimestamptz(filter.CreatedTo),
		UpdatedFrom:    ToPgTypeTimestamptz(filter.UpdatedFrom),
		UpdatedTo:      ToPgTypeTimestamptz(filter.UpdatedTo),
		SortOrder:      ToStringPtr(filter.SortOrder),
		SortBy:         ToStringPtr(filter.SortBy),
	}
	rows, err := db.queries.GetEmployees(ctx, arg)
	if err != nil {
		return nil, 0, err
	}
	var employees []models.Employee
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		employees = append(employees, convertGeneratedEmployeeRowToModel(row))
	}
	return employees, totalCount, nil
}

// TODO: return updated employee
func (db DB) UpdateEmployee(ctx context.Context,
	employeeID int,
	name models.Optional[string],
	status models.Optional[models.EmployeeStatus],
	jobTitle models.Optional[models.EmployeeJobTitle],
	hireDate models.Optional[time.Time],
	salary models.Optional[decimal.Decimal],
	licenseIssued models.Optional[time.Time],
	licenseExpiration models.Optional[time.Time],
) error {
	arg := generated.UpdateEmployeeParams{
		EmployeeID: int32(employeeID),
		Name:       ToStringPtr(name),
		Status:     ToNullEmployeeStatus(status),
		JobTitle: func() *string {
			if jobTitle.Set {
				s := string(jobTitle.Value)
				return &s
			}
			return nil
		}(),
		HireDate:          ToPgTypeDate(hireDate),
		Salary:            ToPgTypeNumeric(salary),
		LicenseIssued:     ToPgTypeDate(licenseIssued),
		LicenseExpiration: ToPgTypeDate(licenseExpiration),
	}
	return db.queries.UpdateEmployee(ctx, arg)
}

func (db DB) SoftDeleteEmployee(ctx context.Context, employeeID int) error {
	return db.queries.SoftDeleteEmployee(ctx, int32(employeeID))
}

func (db DB) HardDeleteEmployee(ctx context.Context, employeeID int) error {
	return db.queries.HardDeleteEmployee(ctx, int32(employeeID))
}

func (db DB) RestoreEmployee(ctx context.Context, employeeID int) error {
	return db.queries.RestoreEmployee(ctx, int32(employeeID))
}

func (db DB) BulkSoftDeleteEmployees(ctx context.Context, employeeIDs []int) error {
	return db.queries.BulkSoftDeleteEmployees(ctx, convertIntSliceToInt32(employeeIDs))
}

func (db DB) BulkHardDeleteEmployees(ctx context.Context, employeeIDs []int) error {
	return db.queries.BulkHardDeleteEmployees(ctx, convertIntSliceToInt32(employeeIDs))
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
