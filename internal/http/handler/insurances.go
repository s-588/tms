package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

// ============================================================================
// Page & Table Handlers
// ============================================================================

func (h Handler) GetInsurancesPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseInsuranceFilters(r)

	insurances, total, err := h.DB.GetInsurances(r.Context(), 1, models.InsuranceFilter{})
	if err != nil {
		slog.Error("can't retrieve list of insurances", "error", err)
		renderToast(w, r, "error", "Can't render insurances page", "Something went wrong")
		return
	}
	err = ui.InsurancesPage(insurances, page, total, filter).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

func (h Handler) GetInsurances(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseInsuranceFilters(r)

	insurances, total, err := h.DB.GetInsurances(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of insurances", "error", err)
		renderToast(w, r, "error", "Can't get insurances data", "Something went wrong")
		return
	}

	slog.Debug("retrieve insurances from database", "filter", filter, "page", page,
		"total pages", total, "total insurances", len(insurances))
	err = ui.InsurancesTable(insurances, page, total, filter, true).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

//nolint:funlen,gocyclo // function is clear and understandable
func parseInsuranceFilters(r *http.Request) models.InsuranceFilter {
	filter := models.InsuranceFilter{}
	q := r.URL.Query()

	if transportID := q.Get("transport_id"); transportID != "" {
		if val, err := strconv.ParseInt(transportID, 10, 32); err != nil && val > 0 {
			filter.TransportID.SetValue(int32(val))
		}
	}
	if insuranceFrom := q.Get("insurance_from"); insuranceFrom != "" {
		if t, err := time.Parse("2006-01-02", insuranceFrom); err == nil {
			filter.InsuranceDateFrom.SetValue(t)
		}
	}
	if insuranceTo := q.Get("insurance_to"); insuranceTo != "" {
		if t, err := time.Parse("2006-01-02", insuranceTo); err == nil {
			filter.InsuranceDateTo.SetValue(t)
		}
	}
	if expirationFrom := q.Get("expiration_from"); expirationFrom != "" {
		if t, err := time.Parse("2006-01-02", expirationFrom); err == nil {
			filter.InsuranceExpirationFrom.SetValue(t)
		}
	}
	if expirationTo := q.Get("expiration_to"); expirationTo != "" {
		if t, err := time.Parse("2006-01-02", expirationTo); err == nil {
			filter.InsuranceExpirationTo.SetValue(t)
		}
	}
	if paymentMin := q.Get("payment_min"); paymentMin != "" {
		if d, err := decimal.NewFromString(paymentMin); err == nil && d.IsPositive() {
			filter.PaymentMin.SetValue(d)
		}
	}
	if paymentMax := q.Get("payment_max"); paymentMax != "" {
		if d, err := decimal.NewFromString(paymentMax); err == nil && d.IsPositive() {
			filter.PaymentMax.SetValue(d)
		}
	}
	if coverageMin := q.Get("coverage_min"); coverageMin != "" {
		if d, err := decimal.NewFromString(coverageMin); err == nil && d.IsPositive() {
			filter.CoverageMin.SetValue(d)
		}
	}
	if coverageMax := q.Get("coverage_max"); coverageMax != "" {
		if d, err := decimal.NewFromString(coverageMax); err == nil && d.IsPositive() {
			filter.CoverageMax.SetValue(d)
		}
	}
	if sortBy := q.Get("sort"); sortBy != "" {
		filter.SortBy.SetValue(sortBy)
	} else {
		filter.SortBy.SetValue("insurance_id")
	}
	if sortOrder := q.Get("order"); sortOrder != "" {
		filter.SortOrder.SetValue(sortOrder)
	} else {
		filter.SortOrder.SetValue("desc")
	}
	return filter
}

func (h Handler) CreateInsuranceHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't parse form", "Invalid form data")
		return
	}

	form, err := parseInsuranceCreateForm(r)
	if err != nil {
		slog.Debug("incorrect input data for adding insurance", "data", form)
		err := ui.InsurancesAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	transportID, err := strconv.ParseInt(form["transport_id"].Value, 10, 32)
	if err != nil {
		slog.Error("can't parse transport ID", "error", err)
		return
	}
	insuranceDate, err := time.Parse("2006-01-02", form["insurance_date"].Value)
	if err != nil {
		slog.Error("can't parse insurance date", "error", err)
		return
	}
	expirationDate, err := time.Parse("2006-01-02", form["insurance_expiration"].Value)
	if err != nil {
		slog.Error("can't parse expiration date", "error", err)
		return
	}
	payment, err := decimal.NewFromString(form["payment"].Value)
	if err != nil {
		slog.Error("can't parse payment", "error", err)
		return
	}
	coverage, err := decimal.NewFromString(form["coverage"].Value)
	if err != nil {
		slog.Error("can't parse coverage", "error", err)
		return
	}

	_, err = h.DB.CreateInsurance(r.Context(), db.CreateInsuranceArgs{
		TransportID:         int32(transportID),
		InsuranceDate:       insuranceDate,
		InsuranceExpiration: expirationDate,
		Payment:             payment,
		Coverage:            coverage,
	})
	if err != nil {
		slog.Error("can't create insurance", "error", err)
		renderToast(w, r, "error", "Can't create insurance", "Something went wrong")
		err = ui.InsurancesAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	slog.Debug("adding new insurance", "data", form)
	renderToast(w, r, "success", "Insurance created", "Insurance successfully created")
	h.GetInsurances(w, r)
}

func parseInsuranceCreateForm(r *http.Request) (form ui.Form, err error) {
	form = make(ui.Form)

	// Transport ID
	transportID := strings.TrimSpace(r.PostForm.Get("transport_id"))
	form["transport_id"] = ui.FormField{Value: transportID}
	if transportID == "" {
		err = errors.New("transport ID is required")
		form["transport_id"] = ui.FormField{Value: transportID, Err: err}
	} else {
		if val, e := strconv.Atoi(transportID); e != nil || val <= 0 {
			err = errors.New("transport ID must be a positive integer")
			form["transport_id"] = ui.FormField{Value: transportID, Err: err}
		}
	}

	// Insurance Date
	insuranceDate := strings.TrimSpace(r.PostForm.Get("insurance_date"))
	form["insurance_date"] = ui.FormField{Value: insuranceDate}
	if insuranceDate == "" {
		err = errors.New("insurance date is required")
		form["insurance_date"] = ui.FormField{Value: insuranceDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", insuranceDate); e != nil {
			err = errors.New("insurance date must be in YYYY-MM-DD format")
			form["insurance_date"] = ui.FormField{Value: insuranceDate, Err: err}
		}
	}

	// Expiration Date
	expirationDate := strings.TrimSpace(r.PostForm.Get("insurance_expiration"))
	form["insurance_expiration"] = ui.FormField{Value: expirationDate}
	if expirationDate == "" {
		err = errors.New("expiration date is required")
		form["insurance_expiration"] = ui.FormField{Value: expirationDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", expirationDate); e != nil {
			err = errors.New("expiration date must be in YYYY-MM-DD format")
			form["insurance_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	// Payment
	payment := strings.TrimSpace(r.PostForm.Get("payment"))
	form["payment"] = ui.FormField{Value: payment}
	if payment == "" {
		err = errors.New("payment is required")
		form["payment"] = ui.FormField{Value: payment, Err: err}
	} else {
		if d, e := decimal.NewFromString(payment); e != nil || d.IsNegative() {
			err = errors.New("payment must be a positive decimal")
			form["payment"] = ui.FormField{Value: payment, Err: err}
		}
	}

	// Coverage
	coverage := strings.TrimSpace(r.PostForm.Get("coverage"))
	form["coverage"] = ui.FormField{Value: coverage}
	if coverage == "" {
		err = errors.New("coverage is required")
		form["coverage"] = ui.FormField{Value: coverage, Err: err}
	} else {
		if d, e := decimal.NewFromString(coverage); e != nil || d.IsNegative() {
			err = errors.New("coverage must be a positive decimal")
			form["coverage"] = ui.FormField{Value: coverage, Err: err}
		}
	}

	// Cross-field validation: expiration must be after insurance date
	if err == nil {
		insT, _ := time.Parse("2006-01-02", insuranceDate)
		expT, _ := time.Parse("2006-01-02", expirationDate)
		if expT.Before(insT) {
			err = errors.New("expiration date must be after insurance date")
			form["insurance_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	return
}

func (h Handler) GetInsuranceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Can't get insurance data", "Something went wrong")
		return
	}
	insurance, err := h.DB.GetInsuranceByID(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve insurance", "error", err, "id", id)
		renderToast(w, r, "error", "Can't get insurance data", "Not found")
		return
	}
	slog.Debug("retrieve insurance", "insurance", insurance)
	err = ui.InsurancesViewSheetContent(insurance, ui.Form{}).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

func (h Handler) UpdateInsuranceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Error", "Incorrect insurance ID")
		return
	}

	existing, err := h.DB.GetInsuranceByID(r.Context(), id)
	if err != nil {
		slog.Error("can't receive insurance", "error", err)
		renderToast(w, r, "error", "Internal error", "Something went wrong")
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		renderToast(w, r, "error", "Bad request", "Invalid form format")
		return
	}

	form, err := parseInsuranceUpdateForm(r, existing)
	if err != nil {
		slog.Debug("can't update insurance", "form", form, "err", err)
		err = ui.InsurancesViewSheetContent(existing, form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	transportID, _ := strconv.ParseInt(form["transport_id"].Value, 10, 32)
	insuranceDate, _ := time.Parse("2006-01-02", form["insurance_date"].Value)
	expirationDate, _ := time.Parse("2006-01-02", form["insurance_expiration"].Value)
	payment, _ := decimal.NewFromString(form["payment"].Value)
	coverage, _ := decimal.NewFromString(form["coverage"].Value)

	if err := h.DB.UpdateInsurance(r.Context(), db.UpdateInsuranceArgs{
		InsuranceID:         id,
		TransportID:         int32(transportID),
		InsuranceDate:       insuranceDate,
		InsuranceExpiration: expirationDate,
		Payment:             payment,
		Coverage:            coverage,
	}); err != nil {
		slog.Error("can't update insurance", "error", err, "id", id)
		renderToast(w, r, "error", "Internal error", "something went wrong")
		return
	}

	slog.Debug("update insurance", "form data", form)
	renderToast(w, r, "success", "Insurance updated", "Insurance successfully updated")
	h.GetInsuranceHandler(w, r)
	h.GetInsurances(w, r)
}

func parseInsuranceUpdateForm(r *http.Request, existing models.Insurance) (form ui.Form, err error) {
	form = make(ui.Form)

	getValue := func(key string, defaultValue string) string {
		if val := r.PostForm.Get(key); val != "" {
			return val
		}
		return defaultValue
	}

	// Transport ID
	transportID := getValue("transport_id", strconv.Itoa(int(existing.TransportID)))
	form["transport_id"] = ui.FormField{Value: transportID}
	if transportID == "" {
		err = errors.New("transport ID is required")
		form["transport_id"] = ui.FormField{Value: transportID, Err: err}
	} else {
		if val, e := strconv.Atoi(transportID); e != nil || val <= 0 {
			err = errors.New("transport ID must be a positive integer")
			form["transport_id"] = ui.FormField{Value: transportID, Err: err}
		}
	}

	// Insurance Date
	insuranceDate := getValue("insurance_date", existing.InsuranceDate.Format("2006-01-02"))
	form["insurance_date"] = ui.FormField{Value: insuranceDate}
	if insuranceDate == "" {
		err = errors.New("insurance date is required")
		form["insurance_date"] = ui.FormField{Value: insuranceDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", insuranceDate); e != nil {
			err = errors.New("insurance date must be in YYYY-MM-DD format")
			form["insurance_date"] = ui.FormField{Value: insuranceDate, Err: err}
		}
	}

	// Expiration Date
	expirationDate := getValue("insurance_expiration", existing.InsuranceExpiration.Format("2006-01-02"))
	form["insurance_expiration"] = ui.FormField{Value: expirationDate}
	if expirationDate == "" {
		err = errors.New("expiration date is required")
		form["insurance_expiration"] = ui.FormField{Value: expirationDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", expirationDate); e != nil {
			err = errors.New("expiration date must be in YYYY-MM-DD format")
			form["insurance_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	// Payment
	payment := getValue("payment", existing.Payment.String())
	form["payment"] = ui.FormField{Value: payment}
	if payment == "" {
		err = errors.New("payment is required")
		form["payment"] = ui.FormField{Value: payment, Err: err}
	} else {
		if d, e := decimal.NewFromString(payment); e != nil || d.IsNegative() {
			err = errors.New("payment must be a positive decimal")
			form["payment"] = ui.FormField{Value: payment, Err: err}
		}
	}

	// Coverage
	coverage := getValue("coverage", existing.Coverage.String())
	form["coverage"] = ui.FormField{Value: coverage}
	if coverage == "" {
		err = errors.New("coverage is required")
		form["coverage"] = ui.FormField{Value: coverage, Err: err}
	} else {
		if d, e := decimal.NewFromString(coverage); e != nil || d.IsNegative() {
			err = errors.New("coverage must be a positive decimal")
			form["coverage"] = ui.FormField{Value: coverage, Err: err}
		}
	}

	// Cross-field validation: expiration must be after insurance date
	if err == nil {
		insT, _ := time.Parse("2006-01-02", insuranceDate)
		expT, _ := time.Parse("2006-01-02", expirationDate)
		if expT.Before(insT) {
			err = errors.New("expiration date must be after insurance date")
			form["insurance_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	return
}

func (h Handler) DeleteInsuranceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Incorrect URL", "Can't parse id from URL path")
		return
	}

	if err := h.DB.SoftDeleteInsurance(r.Context(), id); err != nil {
		slog.Error("can't delete insurance", "error", err, "id", id)
		renderToast(w, r, "error", "Can't delete insurance", "Something went wrong")
		return
	}

	slog.Debug("deleting insurance", "insuranceID", id)
	renderToast(w, r, "success", "Deleted", "Insurance successfully deleted")
	h.GetInsurances(w, r)
}

func (h Handler) BulkDeleteInsurancesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't delete insurances", "Can't parse form")
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		err := ui.Toast("error", "Can't delete insurances", "No insurances selected").Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			slog.Error("can't parse insurance id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, int32(id))
	}

	if err := h.DB.BulkSoftDeleteInsurances(r.Context(), ids); err != nil {
		slog.Error("can't delete insurances batch", "error", err)
		renderToast(w, r, "error", "Can't delete insurances", "Something went wrong")
	}

	h.GetInsurances(w, r)
}
