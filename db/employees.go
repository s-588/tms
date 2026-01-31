package db

import (
	"context"

	"github.com/s-588/tms/db/generated"
	"github.com/s-588/tms/models"
)

func (db DB) CreateEmployeeWrapper(ctx context.Context, name string) (models.Employee, error) {
	genEmployee, err := db.queries.CreateEmployee(ctx, name)
	if err != nil {
		return models.Employee{}, err
	}

	return convertGeneratedEmployeeToModel(genEmployee), nil
}

func (db DB) DeleteEmployeeWrapper(ctx context.Context, employeeID int) error {
	return db.queries.DeleteEmployee(ctx, int32(employeeID))
}

func (db DB) GetEmployeeByIDWrapper(ctx context.Context, employeeID int) (models.Employee, error) {
	genEmployee, err := db.queries.GetEmployeeByemployee_id(ctx, int32(employeeID))
	if err != nil {
		return models.Employee{}, err
	}

	return convertGeneratedEmployeeToModel(genEmployee), nil
}

func (db DB) GetEmployeesPaginatedWrapper(ctx context.Context, limit, offset int) ([]models.Employee, int64, error) {
	arg := generated.GetEmployeesPaginatedParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	rows, err := db.queries.GetEmployeesPaginated(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	var employees []models.Employee
	var totalCount int64
	for _, row := range rows {
		totalCount = row.TotalCount
		employee := models.Employee{
			EmployeeID: row.EmployeeID,
			Name:       row.Name,
			CreatedAt:  row.CreatedAt,
			UpdatedAt:  row.UpdatedAt,
			DeletedAt:  row.DeletedAt,
		}
		employees = append(employees, employee)
	}

	return employees, totalCount, nil
}

func (db DB) UpdateEmployeeWrapper(ctx context.Context, employeeID int, name string) error {
	arg := generated.UpdateEmployeeParams{
		EmployeeID: int32(employeeID),
		Name:       name,
	}

	return db.queries.UpdateEmployee(ctx, arg)
}

func convertGeneratedEmployeeToModel(genEmployee generated.Employee) models.Employee {
	return models.Employee{
		EmployeeID: genEmployee.EmployeeID,
		Name:       genEmployee.Name,
		CreatedAt:  genEmployee.CreatedAt,
		UpdatedAt:  genEmployee.UpdatedAt,
		DeletedAt:  genEmployee.DeletedAt,
	}
}
