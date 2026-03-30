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
)

// ============================================================================
// Page & Table Handlers
// ============================================================================

func (h Handler) GetInspectionsPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseInspectionFilters(r)

	inspections, total, err := h.DB.GetInspections(r.Context(), 1, models.InspectionFilter{})
	if err != nil {
		slog.Error("can't retrieve list of inspections", "error", err)
		ui.Toast("error", "Can't render inspections page", "Something went wrong").Render(r.Context(), w)
		return
	}
	ui.InspectionsPage(inspections, page, total, filter).Render(r.Context(), w)
}

func (h Handler) GetInspections(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseInspectionFilters(r)

	inspections, total, err := h.DB.GetInspections(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of inspections", "error", err)
		ui.Toast("error", "Can't get inspections data", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("retrieve inspections from database", "filter", filter, "page", page,
		"total pages", total, "total inspections", len(inspections))
	ui.InspectionsTable(inspections, page, total, filter, true).Render(r.Context(), w)
}

// ============================================================================
// Filter Parsing
// ============================================================================

func parseInspectionFilters(r *http.Request) models.InspectionFilter {
	filter := models.InspectionFilter{}
	q := r.URL.Query()

	if transportID := q.Get("transport_id"); transportID != "" {
		if val, err := strconv.Atoi(transportID); err == nil && val > 0 {
			filter.TransportID.SetValue(int32(val))
		}
	}
	if status := q.Get("status"); status != "" {
		if err := checkInspectionStatus(status); err == nil {
			filter.Status.SetValue(models.InspectionStatus(status))
		}
	}
	if inspectionFrom := q.Get("inspection_from"); inspectionFrom != "" {
		if t, err := time.Parse("2006-01-02", inspectionFrom); err == nil {
			filter.InspectionDateFrom.SetValue(t)
		}
	}
	if inspectionTo := q.Get("inspection_to"); inspectionTo != "" {
		if t, err := time.Parse("2006-01-02", inspectionTo); err == nil {
			filter.InspectionDateTo.SetValue(t)
		}
	}
	if expirationFrom := q.Get("expiration_from"); expirationFrom != "" {
		if t, err := time.Parse("2006-01-02", expirationFrom); err == nil {
			filter.InspectionExpirationFrom.SetValue(t)
		}
	}
	if expirationTo := q.Get("expiration_to"); expirationTo != "" {
		if t, err := time.Parse("2006-01-02", expirationTo); err == nil {
			filter.InspectionExpirationTo.SetValue(t)
		}
	}
	if sortBy := q.Get("sort"); sortBy != "" {
		filter.SortBy.SetValue(sortBy)
	} else {
		filter.SortBy.SetValue("inspection_id")
	}
	if sortOrder := q.Get("order"); sortOrder != "" {
		filter.SortOrder.SetValue(sortOrder)
	} else {
		filter.SortOrder.SetValue("desc")
	}
	return filter
}

// ============================================================================
// Create
// ============================================================================

func (h Handler) CreateInspectionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't parse form", "Invalid form data").Render(r.Context(), w)
		return
	}

	hasError, form := parseInspectionCreateForm(r)
	if hasError != nil {
		slog.Debug("incorrect input data for adding inspection", "data", form)
		ui.InspectionsAddContent(form).Render(r.Context(), w)
		return
	}

	transportID, _ := strconv.ParseInt(form["transport_id"].Value,10,32)
	inspectionDate, _ := time.Parse("2006-01-02", form["inspection_date"].Value)
	expirationDate, _ := time.Parse("2006-01-02", form["inspection_expiration"].Value)

	_, err := h.DB.CreateInspection(r.Context(),db.CreateInspectionArgs{
		TransportID:          int32(transportID),
		InspectionDate:       inspectionDate,
		InspectionExpiration: expirationDate,
		Status:               models.InspectionStatus(form["status"].Value)})
	if err != nil {
		slog.Error("can't create inspection", "error", err)
		ui.Toast("error", "Can't create inspection", "Something went wrong").Render(r.Context(), w)
		ui.InspectionsAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("adding new inspection", "data", form)
	ui.Toast("success", "Inspection created", "Inspection successfully created").Render(r.Context(), w)
	h.GetInspections(w, r)
}

func parseInspectionCreateForm(r *http.Request) (err error, form ui.Form) {
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

	// Status
	status := strings.TrimSpace(r.PostForm.Get("status"))
	form["status"] = ui.FormField{Value: status}
	if status == "" {
		err = errors.New("status is required")
		form["status"] = ui.FormField{Value: status, Err: err}
	} else {
		if e := checkInspectionStatus(status); e != nil {
			err = e
			form["status"] = ui.FormField{Value: status, Err: err}
		}
	}

	// Inspection Date
	inspectionDate := strings.TrimSpace(r.PostForm.Get("inspection_date"))
	form["inspection_date"] = ui.FormField{Value: inspectionDate}
	if inspectionDate == "" {
		err = errors.New("inspection date is required")
		form["inspection_date"] = ui.FormField{Value: inspectionDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", inspectionDate); e != nil {
			err = errors.New("inspection date must be in YYYY-MM-DD format")
			form["inspection_date"] = ui.FormField{Value: inspectionDate, Err: err}
		}
	}

	// Expiration Date
	expirationDate := strings.TrimSpace(r.PostForm.Get("inspection_expiration"))
	form["inspection_expiration"] = ui.FormField{Value: expirationDate}
	if expirationDate == "" {
		err = errors.New("expiration date is required")
		form["inspection_expiration"] = ui.FormField{Value: expirationDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", expirationDate); e != nil {
			err = errors.New("expiration date must be in YYYY-MM-DD format")
			form["inspection_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	// Cross-field validation: expiration must be after inspection
	if err == nil {
		inspT, _ := time.Parse("2006-01-02", inspectionDate)
		expT, _ := time.Parse("2006-01-02", expirationDate)
		if expT.Before(inspT) {
			err = errors.New("expiration date must be after inspection date")
			form["inspection_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	return
}

// ============================================================================
// Read (single inspection for sheet)
// ============================================================================

func (h Handler) GetInspectionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Can't get inspection data", "Something went wrong").Render(r.Context(), w)
		return
	}
	inspection, err := h.DB.GetInspectionByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't retrieve inspection", "error", err, "id", id)
		ui.Toast("error", "Can't get inspection data", "Not found").Render(r.Context(), w)
		return
	}
	slog.Debug("retrieve inspection", "inspection", inspection)
	ui.InspectionsViewSheetContent(inspection, ui.Form{}).Render(r.Context(), w)
}

// ============================================================================
// Update
// ============================================================================

func (h Handler) UpdateInspectionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Error", "Incorrect inspection ID").Render(r.Context(), w)
		return
	}

	existing, err := h.DB.GetInspectionByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't receive inspection", "error", err)
		ui.Toast("error", "Internal error", "Something went wrong").Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		ui.Toast("error", "Bad request", "Invalid form format").Render(r.Context(), w)
		return
	}

	err, form := parseInspectionUpdateForm(r, existing)
	if err != nil {
		slog.Debug("can't update inspection", "form", form, "err", err)
		ui.InspectionsViewSheetContent(existing, form).Render(r.Context(), w)
		return
	}

	transportID, _ := strconv.ParseInt(form["transport_id"].Value,10,32)
	inspectionDate, _ := time.Parse("2006-01-02", form["inspection_date"].Value)
	expirationDate, _ := time.Parse("2006-01-02", form["inspection_expiration"].Value)

	if err := h.DB.UpdateInspection(r.Context(), db.UpdateInspectionArgs{
		InspectionID:         int32(id),
		TransportID:          int32(transportID),
		InspectionDate:       inspectionDate,
		InspectionExpiration: expirationDate,
		Status:               models.InspectionStatus(form["status"].Value),
	}); err != nil {
		slog.Error("can't update inspection", "error", err, "id", id)
		ui.Toast("error", "Internal error", "something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("update inspection", "form data", form)
	ui.Toast("success", "Inspection updated", "Inspection successfully updated").Render(r.Context(), w)
	h.GetInspectionHandler(w, r)
	h.GetInspections(w, r)
}

func parseInspectionUpdateForm(r *http.Request, existing models.Inspection) (err error, form ui.Form) {
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

	// Status
	status := getValue("status", string(existing.Status))
	form["status"] = ui.FormField{Value: status}
	if status == "" {
		err = errors.New("status is required")
		form["status"] = ui.FormField{Value: status, Err: err}
	} else {
		if e := checkInspectionStatus(status); e != nil {
			err = e
			form["status"] = ui.FormField{Value: status, Err: err}
		}
	}

	// Inspection Date
	inspectionDate := getValue("inspection_date", existing.InspectionDate.Format("2006-01-02"))
	form["inspection_date"] = ui.FormField{Value: inspectionDate}
	if inspectionDate == "" {
		err = errors.New("inspection date is required")
		form["inspection_date"] = ui.FormField{Value: inspectionDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", inspectionDate); e != nil {
			err = errors.New("inspection date must be in YYYY-MM-DD format")
			form["inspection_date"] = ui.FormField{Value: inspectionDate, Err: err}
		}
	}

	// Expiration Date
	expirationDate := getValue("inspection_expiration", existing.InspectionExpiration.Format("2006-01-02"))
	form["inspection_expiration"] = ui.FormField{Value: expirationDate}
	if expirationDate == "" {
		err = errors.New("expiration date is required")
		form["inspection_expiration"] = ui.FormField{Value: expirationDate, Err: err}
	} else {
		if _, e := time.Parse("2006-01-02", expirationDate); e != nil {
			err = errors.New("expiration date must be in YYYY-MM-DD format")
			form["inspection_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	// Cross-field validation: expiration must be after inspection
	if err == nil {
		inspT, _ := time.Parse("2006-01-02", inspectionDate)
		expT, _ := time.Parse("2006-01-02", expirationDate)
		if expT.Before(inspT) {
			err = errors.New("expiration date must be after inspection date")
			form["inspection_expiration"] = ui.FormField{Value: expirationDate, Err: err}
		}
	}

	return
}

// ============================================================================
// Delete
// ============================================================================

func (h Handler) DeleteInspectionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Incorrect URL", "Can't parse id from URL path").Render(r.Context(), w)
		return
	}

	if err := h.DB.SoftDeleteInspection(r.Context(), int32(id)); err != nil {
		slog.Error("can't delete inspection", "error", err, "id", id)
		ui.Toast("error", "Can't delete inspection", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("deleting inspection", "inspectionID", id)
	ui.Toast("success", "Deleted", "Inspection successfully deleted").Render(r.Context(), w)
	h.GetInspections(w, r)
}

// ============================================================================
// Bulk Delete
// ============================================================================

func (h Handler) BulkDeleteInspectionsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't delete inspections", "Can't parse form").Render(r.Context(), w)
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		ui.Toast("error", "Can't delete inspections", "No inspections selected").Render(r.Context(), w)
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("can't parse inspection id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, int32(id))
	}

	if err := h.DB.BulkSoftDeleteInspections(r.Context(), ids); err != nil {
		slog.Error("can't delete inspections batch", "error", err)
	}

	h.GetInspections(w, r)
}

// ============================================================================
// Validation Helpers
// ============================================================================

func checkInspectionStatus(status string) error {
	switch models.InspectionStatus(status) {
	case models.InspectionStatusReady, models.InspectionStatusRepair, models.InspectionStatusOverdue:
		return nil
	default:
		return errors.New("invalid inspection status (must be ready, repair, or overdue)")
	}
}