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
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

// ============================================================================
// Page & Table Handlers
// ============================================================================

// GetEmployeesPage renders the full employees page.
func (h Handler) GetEmployeesPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseEmployeeFilters(r)

	slog.Debug("getting get employees", "page", page, "filter", filter)
	employees, total, err := h.DB.GetEmployees(r.Context(), 1, models.EmployeeFilter{})
	if err != nil {
		slog.Error("can't retrieve list of employees", "error", err)
		renderToast(w, r, "error", "Can't render employees page", "Something went wrong")
		return
	}
	err = ui.EmployeesPage(employees, page, total, filter).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

// GetEmployees returns the employees table (for HTMX partial updates).
func (h Handler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseEmployeeFilters(r)

	employees, total, err := h.DB.GetEmployees(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of employees", "error", err)
		renderToast(w, r, "error", "Can't get employees data", "Something went wrong")
		return
	}

	slog.Debug("retrieve employees from database", "filter", filter, "page", page,
		"total pages", total, "total employees", len(employees))
	err = ui.EmployeesTable(employees, page, total, filter, true).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

func parseEmployeeFilters(r *http.Request) models.EmployeeFilter {
	filter := models.EmployeeFilter{}
	q := r.URL.Query()

	if name := q.Get("name"); name != "" && checkEmployeeNameFilter(name) == nil {
		filter.Name.SetValue(name)
	}
	if job := q.Get("job_title"); job != "" {
		j := models.EmployeeJobTitle(job)
		filter.JobTitle.SetValue(j)
	}
	if status := q.Get("status"); status != "" {
		if err := checkEmployeeStatus(status); err == nil {
			filter.Status.SetValue(models.EmployeeStatus(status))
		}
	}
	if salaryMin := q.Get("salary_min"); salaryMin != "" {
		if d, err := decimal.NewFromString(salaryMin); err == nil && d.IsPositive() {
			filter.SalaryMin.SetValue(d)
		}
	}
	if salaryMax := q.Get("salary_max"); salaryMax != "" {
		if d, err := decimal.NewFromString(salaryMax); err == nil && d.IsPositive() {
			filter.SalaryMax.SetValue(d)
		}
	}
	if q.Has("sort") {
		filter.SortBy.SetValue(q.Get("sort"))
	} else {
		filter.SortBy.SetValue("employee_id")
	}
	if q.Has("order") {
		filter.SortOrder.SetValue(q.Get("order"))
	} else {
		filter.SortOrder.SetValue("desc")
	}
	return filter
}

// ============================================================================
// Create
// ============================================================================

func (h Handler) CreateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't parse form", "Invalid form data")
		return
	}

	form, err := parseEmployeeCreateForm(r)
	if err != nil {
		slog.Debug("incorrect input data for adding employee", "data", form)
		err := ui.EmployeesAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	// Convert form values to proper types
	hireDate, _ := time.Parse("2006-01-02", form["hire_date"].Value)
	salary, _ := decimal.NewFromString(form["salary"].Value)
	licenseIssued, _ := time.Parse("2006-01-02", form["license_issued"].Value)
	licenseExpiration, _ := time.Parse("2006-01-02", form["license_expiration"].Value)

	emp := models.Employee{
		Name:              form["name"].Value,
		Status:            models.EmployeeStatus(form["status"].Value),
		JobTitle:          models.EmployeeJobTitle(form["job_title"].Value),
		HireDate:          hireDate,
		Salary:            salary,
		LicenseIssued:     licenseIssued,
		LicenseExpiration: licenseExpiration,
	}

	_, err = h.DB.CreateEmployee(r.Context(), db.CreateEmployeeArgs{
		LicenseExpiration: emp.LicenseExpiration,
		LicenseIssued:     emp.LicenseIssued,
		Salary:            emp.Salary,
		HireDate:          emp.HireDate,
		JobTitle:          emp.JobTitle,
		Status:            emp.Status,
		Name:              emp.Name,
	})
	if err != nil {
		slog.Error("can't create employee", "error", err)
		renderToast(w, r, "error", "Can't create employee", "Something went wrong")
		err = ui.EmployeesAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	slog.Debug("adding new employee", "data", form)
	renderToast(w, r, "success", "Employee created", "Employee successfully created")
	h.GetEmployees(w, r)
}

func parseEmployeeCreateForm(r *http.Request) (form ui.Form, err error) {
	form = make(ui.Form)

	name := strings.TrimSpace(r.PostForm.Get("name"))
	form["name"] = ui.FormField{Value: name}
	if err = checkEmployeeName(name); err != nil {
		form["name"] = ui.FormField{Value: name, Err: err}
	}

	status := r.PostForm.Get("status")
	form["status"] = ui.FormField{Value: status}
	if err = checkEmployeeStatus(status); err != nil {
		form["status"] = ui.FormField{Value: status, Err: err}
	}

	job := r.PostForm.Get("job_title")
	form["job_title"] = ui.FormField{Value: job}
	if err = checkEmployeeJobTitle(job); err != nil {
		form["job_title"] = ui.FormField{Value: job, Err: err}
	}

	hireDateStr := r.PostForm.Get("hire_date")
	form["hire_date"] = ui.FormField{Value: hireDateStr}
	_, errHire := time.Parse("2006-01-02", hireDateStr)
	if errHire != nil {
		err = errors.New("invalid hire date (use YYYY-MM-DD)")
		form["hire_date"] = ui.FormField{Value: hireDateStr, Err: err}
	}

	salaryStr := r.PostForm.Get("salary")
	form["salary"] = ui.FormField{Value: salaryStr}
	salary, errSalary := decimal.NewFromString(salaryStr)
	if errSalary != nil || salary.IsNegative() {
		err = errors.New("invalid salary (positive number expected)")
		form["salary"] = ui.FormField{Value: salaryStr, Err: err}
	}

	licenseIssuedStr := r.PostForm.Get("license_issued")
	form["license_issued"] = ui.FormField{Value: licenseIssuedStr}
	licenseIssued, errIssued := time.Parse("2006-01-02", licenseIssuedStr)
	if errIssued != nil {
		err = errors.New("invalid license issued date (use YYYY-MM-DD)")
		form["license_issued"] = ui.FormField{Value: licenseIssuedStr, Err: err}
	}

	licenseExpStr := r.PostForm.Get("license_expiration")
	form["license_expiration"] = ui.FormField{Value: licenseExpStr}
	licenseExp, errExp := time.Parse("2006-01-02", licenseExpStr)
	if errExp != nil {
		err = errors.New("invalid license expiration date (use YYYY-MM-DD)")
		form["license_expiration"] = ui.FormField{Value: licenseExpStr, Err: err}
	}

	// Cross‑field validation
	if err == nil && licenseExp.Before(licenseIssued) {
		err = errors.New("license expiration must be after issue date")
		form["license_expiration"] = ui.FormField{Value: licenseExpStr, Err: err}
	}
	return
}

func (h Handler) GetEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Can't get employee data", "Something went wrong")
		return
	}
	employee, err := h.DB.GetEmployeeByID(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve employee", "error", err, "id", id)
		renderToast(w, r, "error", "Can't get employee data", "Not found")
		return
	}
	slog.Debug("retrieve employee", "employee", employee)
	err = ui.EmployeesViewSheetContent(employee, ui.Form{}).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

func (h Handler) UpdateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Error", "Incorrect employee ID")
		h.GetEmployeeHandler(w, r)
		return
	}

	existing, err := h.DB.GetEmployeeByID(r.Context(), id)
	if err != nil {
		slog.Error("can't receive employee", "error", err)
		renderToast(w, r, "error", "Internal error", "Something went wrong")
		h.GetEmployeeHandler(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		renderToast(w, r, "error", "Bad request", "Invalid form format")
		h.GetEmployeeHandler(w, r)
		return
	}

	form, args, err := parseEmployeeUpdateForm(r, existing)
	if err != nil {
		slog.Debug("can't update employee", "form", form, "err", err)
		if renderErr := ui.EmployeesViewSheetContent(existing, form).Render(r.Context(), w); renderErr != nil {
			slog.Error("can't render response", "error", renderErr)
		}
		return
	}

	if err := h.DB.UpdateEmployee(r.Context(), args); err != nil {
		slog.Error("can't update employee", "error", err, "id", id)
		renderToast(w, r, "error", "Internal error", "something went wrong")
		h.GetEmployeeHandler(w, r)
		return
	}

	slog.Debug("update employee", "form data", form)
	renderToast(w, r, "success", "Employee updated", "Employee successfully updated")
	h.GetEmployeeHandler(w, r)
	h.GetEmployees(w, r)
}

//nolint:funlen // function is clear and there is no need to divide in multiple functions.
func parseEmployeeUpdateForm(r *http.Request, existing models.Employee) (form ui.Form, args db.UpdateEmployeeArgs, err error) {
	form = make(ui.Form)

	// helpers
	setStr := func(key, val string, check func(string) error) {
		form[key] = ui.FormField{Value: val}
		if e := check(val); e != nil {
			form[key] = ui.FormField{Value: val, Err: e}
			err = e
		}
	}
	setDate := func(key, val string, dest *time.Time, msg string) {
		form[key] = ui.FormField{Value: val}
		t, e := time.Parse("2006-01-02", val)
		if e != nil {
			e = errors.New(msg)
			form[key] = ui.FormField{Value: val, Err: e}
			err = e
			return
		}
		*dest = t
	}

	// Name
	name := strings.TrimSpace(r.PostForm.Get("name"))
	if name == "" {
		name = existing.Name
	}
	setStr("name", name, checkEmployeeName)
	args.Name = name

	// Status
	status := r.PostForm.Get("status")
	if status == "" {
		status = string(existing.Status)
	}
	setStr("status", status, checkEmployeeStatus)
	args.Status = models.EmployeeStatus(status)

	// Job title
	job := r.PostForm.Get("job_title")
	if job == "" {
		job = string(existing.JobTitle)
	}
	setStr("job_title", job, checkEmployeeJobTitle)
	args.JobTitle = models.EmployeeJobTitle(job)

	// Hire date
	hireDateStr := r.PostForm.Get("hire_date")
	if hireDateStr == "" {
		hireDateStr = existing.HireDate.Format("2006-01-02")
	}
	setDate("hire_date", hireDateStr, &args.HireDate, "invalid hire date (use YYYY-MM-DD)")

	// Salary
	salaryStr := r.PostForm.Get("salary")
	if salaryStr == "" {
		salaryStr = existing.Salary.String()
	}
	form["salary"] = ui.FormField{Value: salaryStr}
	salary, e := decimal.NewFromString(salaryStr)
	if e != nil || salary.IsNegative() {
		e = errors.New("invalid salary (positive number expected)")
		form["salary"] = ui.FormField{Value: salaryStr, Err: e}
		err = e
	}
	args.Salary = salary

	// License issued / expiration
	issuedStr := r.PostForm.Get("license_issued")
	if issuedStr == "" {
		issuedStr = existing.LicenseIssued.Format("2006-01-02")
	}
	setDate("license_issued", issuedStr, &args.LicenseIssued, "invalid license issued date (use YYYY-MM-DD)")

	expStr := r.PostForm.Get("license_expiration")
	if expStr == "" {
		expStr = existing.LicenseExpiration.Format("2006-01-02")
	}
	setDate("license_expiration", expStr, &args.LicenseExpiration, "invalid license expiration date (use YYYY-MM-DD)")

	// Cross-field validations (only when both dates parsed cleanly)
	if form["license_issued"].Err == nil && form["license_expiration"].Err == nil {
		if args.LicenseExpiration.Before(args.LicenseIssued) {
			e := errors.New("license expiration must be after issue date")
			form["license_expiration"] = ui.FormField{Value: expStr, Err: e}
			err = e
		}
		if time.Now().After(args.LicenseExpiration) && args.Status == models.EmployeeStatusAssigned {
			e := errors.New("employee with expired license cannot be assigned")
			form["status"] = ui.FormField{Value: status, Err: e}
			err = e
		}
	}

	return form, args, err
}

func (h Handler) DeleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Incorrect URL", "Can't parse id from URL path")
		h.GetEmployees(w, r)
		return
	}

	if err := h.DB.SoftDeleteEmployee(r.Context(), id); err != nil {
		slog.Error("can't delete employee", "error", err, "id", id)
		renderToast(w, r, "error", "Can't delete employee", "Something went wrong")
		return
	}

	slog.Debug("deleting employee", "employeeID", id)
	renderToast(w, r, "success", "Deleted", "Employee successfully deleted")
	h.GetEmployees(w, r)
}

func (h Handler) BulkDeleteEmployeesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't delete employees", "Can't parse form")
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		err := ui.Toast("error", "Can't delete employees", "No employees selected").Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			slog.Error("can't parse employee id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, int32(id))
	}

	if err := h.DB.BulkSoftDeleteEmployees(r.Context(), ids); err != nil {
		slog.Error("can't delete employees batch", "error", err)
	}

	h.GetEmployees(w, r)
}

func checkEmployeeName(name string) error {
	if utf8.RuneCountInString(name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	return nil
}

func checkEmployeeNameFilter(name string) error {
	if name == "" {
		return nil
	}
	return checkEmployeeName(name)
}

func checkEmployeeStatus(status string) error {
	switch models.EmployeeStatus(status) {
	case models.EmployeeStatusAvailable, models.EmployeeStatusAssigned, models.EmployeeStatusUnavailable:
		return nil
	default:
		return errors.New("invalid employee status")
	}
}

func checkEmployeeJobTitle(job string) error {
	switch models.EmployeeJobTitle(job) {
	case models.EmployeeJobTitleDriver, models.EmployeeJobTitleDispatcher,
		models.EmployeeJobTitleMechanic, models.EmployeeJobTitleLogisticsManager:
		return nil
	default:
		return errors.New("invalid job title")
	}
}
