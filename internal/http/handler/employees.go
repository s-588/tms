package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

func parseEmployeeFilter(r *http.Request) models.EmployeeFilter {
	filter := models.EmployeeFilter{}

	// Parse string fields
	if name := r.URL.Query().Get("name"); name != "" {
		filter.Name.SetValue(name)
	}

	// Parse time fields
	if createdFrom := r.URL.Query().Get("created_from"); createdFrom != "" {
		if t, err := time.Parse(time.RFC3339, createdFrom); err == nil {
			filter.CreatedFrom.SetValue(t)
		}
	}

	if createdTo := r.URL.Query().Get("created_to"); createdTo != "" {
		if t, err := time.Parse(time.RFC3339, createdTo); err == nil {
			filter.CreatedTo.SetValue(t)
		}
	}

	if updatedFrom := r.URL.Query().Get("updated_from"); updatedFrom != "" {
		if t, err := time.Parse(time.RFC3339, updatedFrom); err == nil {
			filter.UpdatedFrom.SetValue(t)
		}
	}

	if updatedTo := r.URL.Query().Get("updated_to"); updatedTo != "" {
		if t, err := time.Parse(time.RFC3339, updatedTo); err == nil {
			filter.UpdatedTo.SetValue(t)
		}
	}

	// Parse sort fields
	if sortBy := r.URL.Query().Get("sort"); sortBy != "" {
		filter.SortBy.SetValue(sortBy)
	}

	if sortOrder := r.URL.Query().Get("order"); sortOrder != "" {
		filter.SortOrder.SetValue(sortOrder)
	}

	return filter
}

func (h Handler) GetEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	// limit, offset := parsePagination(r)

	// Parse filter using the new function
	filter := parseEmployeeFilter(r)

	// Set default sort
	if !filter.SortBy.Set {
		filter.SortBy.SetValue("employee_id")
	}
	if !filter.SortOrder.Set {
		filter.SortOrder.SetValue("desc")
	}

	// employees, total, err := h.DB.GetEmployees(r.Context(), limit, offset, filter)
	// if err != nil {
	// 	slog.Error("can't retrieve list of employees", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }

	// ui.EmployeesTable(employees, limit, offset, int(total), filter).Render(r.Context(), w)
}

func (h Handler) BulkDeleteEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responseError(w, r, http.StatusBadRequest, "invalid form data")
		return
	}

	// Get selected IDs from form
	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		responseError(w, r, http.StatusBadRequest, "no employees selected")
		return
	}

	// Convert string IDs to int
	var ids []int
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("can't parse employee id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, id)
	}

	// Delete in batches
	batchSize := 10
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]

		if err := h.DB.BulkSoftDeleteEmployees(r.Context(), batch); err != nil {
			slog.Error("can't delete employees batch", "error", err, "batch", batch)
		}
	}

	// Return updated table
	// limit, offset := parsePagination(r)

	// Use the same filter parsing function
	// filter := parseEmployeeFilter(r)

	// employees, total, err := h.DB.GetEmployees(r.Context(), limit, offset, filter)
	// if err != nil {
	// 	slog.Error("can't retrieve list of employees after delete", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }

	// ui.EmployeesTable(employees, limit, offset, int(total), filter).Render(r.Context(), w)
}

func (h Handler) ExportEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	// limit, offset := parsePagination(r)
	// // Parse filter using the new function
	// filter := parseEmployeeFilter(r)
	//
	// // Get all employees with filters
	// employees, total, err := h.DB.GetEmployees(r.Context(), limit, offset, filter)
	// if err != nil {
	// 	slog.Error("can't retrieve employees for export", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }
	//
	// // Set CSV headers
	// w.Header().Set("Content-Type", "text/csv")
	// w.Header().Set("Content-Disposition", "attachment;filename=employees.csv")
	//
	// // Write CSV header
	// fmt.Fprintln(w, "ID,Name,Created At,Updated At,Deleted At")
	//
	// // Write data
	// for _, employee := range employees {
	// 	deletedAt := ""
	// 	if !employee.DeletedAt.IsZero() {
	// 		deletedAt = employee.DeletedAt.Format(time.RFC3339)
	// 	}
	// 	fmt.Fprintf(w, "%d,%s,%s,%s,%s\n",
	// 		employee.EmployeeID,
	// 		escapeCSV(employee.Name),
	// 		employee.CreatedAt.Format(time.RFC3339),
	// 		employee.UpdatedAt.Format(time.RFC3339),
	// 		deletedAt,
	// 	)
	// }
}

func (h Handler) CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responseError(w, r, http.StatusBadRequest, "invalid form data")
		return
	}

	emp, err := parseEmployeeCreateForm(r)
	employee, err := h.DB.CreateEmployee(r.Context(), emp)
	if err != nil {
		slog.Error("can't create employee", "error", err)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	ui.CreateSuccess("Employee added", "employees", int(employee.EmployeeID)).Render(r.Context(), w)
}

func parseEmployeeCreateForm(r *http.Request) (emp models.Employee, err error) {
	if err = r.ParseForm(); err != nil {
		return
	}

	emp.Name = strings.TrimSpace(r.FormValue("name"))
	if utf8.RuneCountInString(emp.Name) < 2 {
		err = errors.New("name must be at least 2 characters")
		return
	}

	// Status
	statusStr := r.FormValue("status")
	emp.Status = models.EmployeeStatus(statusStr)
	switch emp.Status {
	case models.EmployeeStatusAvailable,
		models.EmployeeStatusAssigned,
		models.EmployeeStatusUnavailable:
	default:
		err = errors.New("invalid employee status")
		return
	}

	// Job Title
	jobTitleStr := r.FormValue("job_title")
	emp.JobTitle = models.EmployeeJobTitle(jobTitleStr)
	switch emp.JobTitle {
	case models.EmployeeJobTitleDriver,
		models.EmployeeJobTitleDispatcher,
		models.EmployeeJobTitleMechanic,
		models.EmployeeJobTitleLogisticsManager:
	default:
		err = errors.New("invalid job title")
		return
	}

	// Hire Date
	emp.HireDate, err = time.Parse("2006-01-02", r.FormValue("hire_date"))
	if err != nil {
		err = errors.New("invalid hire date")
		return
	}

	// Salary
	emp.Salary, err = decimal.NewFromString(r.FormValue("salary"))
	if err != nil || emp.Salary.IsNegative() {
		err = errors.New("invalid salary")
		return
	}

	// License Issued
	emp.LicenseIssued, err = time.Parse("2006-01-02", r.FormValue("license_issued"))
	if err != nil {
		err = errors.New("invalid license issued date")
		return
	}

	// License Expiration
	emp.LicenseExpiration, err = time.Parse("2006-01-02", r.FormValue("license_expiration"))
	if err != nil {
		err = errors.New("invalid license expiration date")
		return
	}

	if emp.LicenseExpiration.Before(emp.LicenseIssued) {
		err = errors.New("license expiration must be after issue date")
		return
	}

	return
}
func (h Handler) DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		responseError(w, r, http.StatusBadRequest, "incorrect employee id")
		return
	}

	if err := h.DB.SoftDeleteEmployee(r.Context(), id); err != nil {
		slog.Error("can't delete employee", "error", err, "id", id)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		responseError(w, r, http.StatusBadRequest, "incorrect employee id")
		return
	}

	// Fetch and return the updated employee detail view
	// updatedEmployee, err := h.DB.GetEmployeeByID(r.Context(), id)
	// if err != nil {
	// 	slog.Error("can't fetch updated employee", "error", err, "id", id)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }

	args, errs := parseEmployeeUpdateForm(r)
	if len(errs) > 0 {
		// ui.EmployeeEditForm(updatedEmployee, errs)
	}

	// Perform the update
	if err := h.DB.UpdateEmployee(r.Context(), id, args.Name,
		args.Status, args.JobTitle, args.HireDate, args.Salary, args.LicenseIssued,
		args.LicenseExpiration); err != nil {
		slog.Error("can't update employee", "error", err, "id", id)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	// ui.EmployeeDetail(updatedEmployee).Render(r.Context(), w)
}

type employeeUpdateParams struct {
	EmployeeID        models.Optional[int32]
	Name              models.Optional[string]
	Status            models.Optional[models.EmployeeStatus]
	JobTitle          models.Optional[models.EmployeeJobTitle]
	HireDate          models.Optional[time.Time]
	Salary            models.Optional[decimal.Decimal]
	LicenseIssued     models.Optional[time.Time]
	LicenseExpiration models.Optional[time.Time]
}

func parseEmployeeUpdateForm(r *http.Request) (employeeUpdateParams, map[string]string) {
	var params employeeUpdateParams
	errs := make(map[string]string)

	// Ensure form is parsed
	if err := r.ParseForm(); err != nil {
		errs["form"] = "invalid form data"
		return params, errs
	}

	// Name
	if name := r.PostForm.Get("name"); name != "" {
		if utf8.RuneCountInString(name) < 3 { // assuming at least 3 characters
			errs["name"] = "employee name must be at least 3 characters"
		} else {
			params.Name.SetValue(name)
		}
	} // empty name means no update

	// Status
	if statusStr := r.PostForm.Get("status"); statusStr != "" {
		status := models.EmployeeStatus(statusStr)
		switch status {
		case models.EmployeeStatusAvailable, models.EmployeeStatusAssigned, models.EmployeeStatusUnavailable:
			params.Status.SetValue(status)
		default:
			errs["status"] = "invalid employee status"
		}
	}

	// Job Title
	if jobStr := r.PostForm.Get("job_title"); jobStr != "" {
		job := models.EmployeeJobTitle(jobStr)
		switch job {
		case models.EmployeeJobTitleDriver, models.EmployeeJobTitleDispatcher,
			models.EmployeeJobTitleMechanic, models.EmployeeJobTitleLogisticsManager:
			params.JobTitle.SetValue(job)
		default:
			errs["job_title"] = "invalid job title"
		}
	}

	// Hire Date (expected format YYYY-MM-DD)
	if hireDateStr := r.PostForm.Get("hire_date"); hireDateStr != "" {
		hireDate, err := time.Parse("2006-01-02", hireDateStr)
		if err != nil {
			errs["hire_date"] = "invalid hire date format, use YYYY-MM-DD"
		} else {
			params.HireDate.SetValue(hireDate)
		}
	}

	// Salary (decimal)
	if salaryStr := r.PostForm.Get("salary"); salaryStr != "" {
		salary, err := decimal.NewFromString(salaryStr)
		if err != nil {
			errs["salary"] = "invalid salary format"
		} else {
			params.Salary.SetValue(salary)
		}
	}

	// License Issued
	if licenseIssuedStr := r.PostForm.Get("license_issued"); licenseIssuedStr != "" {
		licenseIssued, err := time.Parse("2006-01-02", licenseIssuedStr)
		if err != nil {
			errs["license_issued"] = "invalid license issued date format, use YYYY-MM-DD"
		} else {
			params.LicenseIssued.SetValue(licenseIssued)
		}
	}

	// License Expiration
	if licenseExpStr := r.PostForm.Get("license_expiration"); licenseExpStr != "" {
		licenseExp, err := time.Parse("2006-01-02", licenseExpStr)
		if err != nil {
			errs["license_expiration"] = "invalid license expiration date format, use YYYY-MM-DD"
		} else {
			params.LicenseExpiration.SetValue(licenseExp)
		}
	}

	return params, errs
}
func (h Handler) NewEmployeePageHandler(w http.ResponseWriter, r *http.Request) {
	// ui.EmployeeCreateForm(nil).Render(r.Context(), w)
}

func (h Handler) EditEmployeePageHandler(w http.ResponseWriter, r *http.Request) {
	// id, err := parseIDFromReq(r)
	// if err != nil {
	// 	slog.Error("can't parse id from URL path", "error", err)
	// 	responseError(w, r, http.StatusBadRequest, "incorrect employee id")
	// 	return
	// }
	//
	// employee, err := h.DB.GetEmployeeByID(r.Context(), id)
	// if err != nil {
	// 	slog.Error("can't retrieve employee for edit", "error", err, "id", id)
	// 	responseError(w, r, http.StatusNotFound, "employee not found")
	// 	return
	// }

	// ui.EmployeeEditForm(employee, nil).Render(r.Context(), w)
}
