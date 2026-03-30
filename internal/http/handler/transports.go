package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
)

// ============================================================================
// Page & Table Handlers
// ============================================================================

func (h Handler) GetTransportsPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseTransportFilters(r)

	transports, total, err := h.DB.GetTransports(r.Context(), 1, models.TransportFilter{})
	if err != nil {
		slog.Error("can't retrieve list of transports", "error", err)
		ui.Toast("error", "Can't render transports page", "Something went wrong").Render(r.Context(), w)
		return
	}
	ui.TransportsPage(transports, page, total, filter).Render(r.Context(), w)
}

func (h Handler) GetTransports(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseTransportFilters(r)

	transports, total, err := h.DB.GetTransports(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of transports", "error", err)
		ui.Toast("error", "Can't get transports data", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("retrieve transports from database", "filter", filter, "page", page,
		"total pages", total, "total transports", len(transports))
	ui.TransportsTable(transports, page, total, filter, true).Render(r.Context(), w)
}

// ============================================================================
// Filter Parsing
// ============================================================================

func parseTransportFilters(r *http.Request) models.TransportFilter {
	filter := models.TransportFilter{}
	q := r.URL.Query()

	if model := q.Get("model"); model != "" {
		filter.Model.SetValue(model)
	}
	if license := q.Get("license_plate"); license != "" {
		filter.LicensePlate.SetValue(license)
	}
	if payloadMin := q.Get("payload_min"); payloadMin != "" {
		if val, err := strconv.Atoi(payloadMin); err == nil && val > 0 {
			filter.PayloadCapacityMin.SetValue(int32(val))
		}
	}
	if payloadMax := q.Get("payload_max"); payloadMax != "" {
		if val, err := strconv.Atoi(payloadMax); err == nil && val > 0 {
			filter.PayloadCapacityMax.SetValue(int32(val))
		}
	}
	if fuelMin := q.Get("fuel_min"); fuelMin != "" {
		if val, err := strconv.Atoi(fuelMin); err == nil && val > 0 {
			filter.FuelConsumptionMin.SetValue(int32(val))
		}
	}
	if fuelMax := q.Get("fuel_max"); fuelMax != "" {
		if val, err := strconv.Atoi(fuelMax); err == nil && val > 0 {
			filter.FuelConsumptionMax.SetValue(int32(val))
		}
	}
	if sortBy := q.Get("sort"); sortBy != "" {
		filter.SortBy.SetValue(sortBy)
	} else {
		filter.SortBy.SetValue("transport_id")
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

func (h Handler) CreateTransportHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't parse form", "Invalid form data").Render(r.Context(), w)
		return
	}

	hasError, form := parseTransportCreateForm(r)
	if hasError != nil {
		slog.Debug("incorrect input data for adding transport", "data", form)
		ui.TransportsAddContent(form).Render(r.Context(), w)
		return
	}

	payload, _ := strconv.Atoi(form["payload_capacity"].Value)
	fuel, _ := strconv.Atoi(form["fuel_consumption"].Value)

	_, err := h.DB.CreateTransport(r.Context(),db.CreateTransportArgs{
		Model: form["model"].Value,
		LicensePlate: form["license_plate"].Value,
		PayloadCapacity: int32(payload),
		FuelConsumption: int32(fuel),
	})
	if err != nil {
		if errors.Is(err, db.ErrDuplicateLicense) {
			form["license_plate"] = ui.FormField{Value: form["license_plate"].Value, Err: errors.New("license plate already exists")}
			ui.TransportsAddContent(form).Render(r.Context(), w)
			return
		}
		slog.Error("can't create transport", "error", err)
		ui.Toast("error", "Can't create transport", "Something went wrong").Render(r.Context(), w)
		ui.TransportsAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("adding new transport", "data", form)
	ui.Toast("success", "Transport created", "Transport successfully created").Render(r.Context(), w)
	h.GetTransports(w, r)
}

func parseTransportCreateForm(r *http.Request) (err error, form ui.Form) {
	form = make(ui.Form)

	// Model
	model := strings.TrimSpace(r.PostForm.Get("model"))
	form["model"] = ui.FormField{Value: model}
	if model == "" {
		err = errors.New("model is required")
		form["model"] = ui.FormField{Value: model, Err: err}
	}

	// License Plate (optional)
	license := strings.TrimSpace(r.PostForm.Get("license_plate"))
	form["license_plate"] = ui.FormField{Value: license}

	// Payload Capacity
	payload := strings.TrimSpace(r.PostForm.Get("payload_capacity"))
	form["payload_capacity"] = ui.FormField{Value: payload}
	if payload == "" {
		err = errors.New("payload capacity is required")
		form["payload_capacity"] = ui.FormField{Value: payload, Err: err}
	} else {
		if val, e := strconv.Atoi(payload); e != nil || val <= 0 {
			err = errors.New("payload capacity must be a positive integer")
			form["payload_capacity"] = ui.FormField{Value: payload, Err: err}
		}
	}

	// Fuel Consumption
	fuel := strings.TrimSpace(r.PostForm.Get("fuel_consumption"))
	form["fuel_consumption"] = ui.FormField{Value: fuel}
	if fuel == "" {
		err = errors.New("fuel consumption is required")
		form["fuel_consumption"] = ui.FormField{Value: fuel, Err: err}
	} else {
		if val, e := strconv.Atoi(fuel); e != nil || val <= 0 {
			err = errors.New("fuel consumption must be a positive integer")
			form["fuel_consumption"] = ui.FormField{Value: fuel, Err: err}
		}
	}

	return
}

// ============================================================================
// Read (single transport for sheet)
// ============================================================================

func (h Handler) GetTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Can't get transport data", "Something went wrong").Render(r.Context(), w)
		return
	}
	transport, err := h.DB.GetTransportByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't retrieve transport", "error", err, "id", id)
		ui.Toast("error", "Can't get transport data", "Not found").Render(r.Context(), w)
		return
	}
	slog.Debug("retrieve transport", "transport", transport)
	ui.TransportsViewSheetContent(transport, ui.Form{}).Render(r.Context(), w)
}

// ============================================================================
// Update
// ============================================================================

func (h Handler) UpdateTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Error", "Incorrect transport ID").Render(r.Context(), w)
		return
	}

	existing, err := h.DB.GetTransportByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't receive transport", "error", err)
		ui.Toast("error", "Internal error", "Something went wrong").Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		ui.Toast("error", "Bad request", "Invalid form format").Render(r.Context(), w)
		return
	}

	err, form := parseTransportUpdateForm(r, existing)
	if err != nil {
		slog.Debug("can't update transport", "form", form, "err", err)
		ui.TransportsViewSheetContent(existing, form).Render(r.Context(), w)
		return
	}

	payload, _ := strconv.Atoi(form["payload_capacity"].Value)
	fuel, _ := strconv.Atoi(form["fuel_consumption"].Value)

	if err := h.DB.UpdateTransport(r.Context(),db.UpdateTransportArgs{
		TransportID: int32(id),
		Model: form["model"].Value,
		LicensePlate: form["license_plate"].Value,
		PayloadCapacity: int32(payload),
		FuelConsumption: int32(fuel),
	}); err != nil {
		if errors.Is(err, db.ErrDuplicateLicense) {
			form["license_plate"] = ui.FormField{Value: form["license_plate"].Value, Err: errors.New("license plate already exists")}
			ui.TransportsAddContent(form).Render(r.Context(), w)
			return
		}
		slog.Error("can't create transport", "error", err)
		ui.Toast("error", "Can't create transport", "Something went wrong").Render(r.Context(), w)
		ui.TransportsAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("update transport", "form data", form)
	ui.Toast("success", "Transport updated", "Transport successfully updated").Render(r.Context(), w)
	h.GetTransportHandler(w, r)
	h.GetTransports(w, r)
}

func parseTransportUpdateForm(r *http.Request, existing models.Transport) (err error, form ui.Form) {
	form = make(ui.Form)

	getValue := func(key string, defaultValue string) string {
		if val := r.PostForm.Get(key); val != "" {
			return val
		}
		return defaultValue
	}

	// Model
	model := getValue("model", existing.Model)
	form["model"] = ui.FormField{Value: model}
	if model == "" {
		err = errors.New("model is required")
		form["model"] = ui.FormField{Value: model, Err: err}
	}

	// License Plate
	license := getValue("license_plate", existing.LicensePlate)
	form["license_plate"] = ui.FormField{Value: license}

	// Payload Capacity
	payload := getValue("payload_capacity", strconv.Itoa(int(existing.PayloadCapacity)))
	form["payload_capacity"] = ui.FormField{Value: payload}
	if payload == "" {
		err = errors.New("payload capacity is required")
		form["payload_capacity"] = ui.FormField{Value: payload, Err: err}
	} else {
		if val, e := strconv.Atoi(payload); e != nil || val <= 0 {
			err = errors.New("payload capacity must be a positive integer")
			form["payload_capacity"] = ui.FormField{Value: payload, Err: err}
		}
	}

	// Fuel Consumption
	fuel := getValue("fuel_consumption", strconv.Itoa(int(existing.FuelConsumption)))
	form["fuel_consumption"] = ui.FormField{Value: fuel}
	if fuel == "" {
		err = errors.New("fuel consumption is required")
		form["fuel_consumption"] = ui.FormField{Value: fuel, Err: err}
	} else {
		if val, e := strconv.Atoi(fuel); e != nil || val <= 0 {
			err = errors.New("fuel consumption must be a positive integer")
			form["fuel_consumption"] = ui.FormField{Value: fuel, Err: err}
		}
	}

	return
}

// ============================================================================
// Delete
// ============================================================================

func (h Handler) DeleteTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Incorrect URL", "Can't parse id from URL path").Render(r.Context(), w)
		return
	}

	if err := h.DB.SoftDeleteTransport(r.Context(), int32(id)); err != nil {
		slog.Error("can't delete transport", "error", err, "id", id)
		ui.Toast("error", "Can't delete transport", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("deleting transport", "transportID", id)
	ui.Toast("success", "Deleted", "Transport successfully deleted").Render(r.Context(), w)
	h.GetTransports(w, r)
}

// ============================================================================
// Bulk Delete
// ============================================================================

func (h Handler) BulkDeleteTransportsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't delete transports", "Can't parse form").Render(r.Context(), w)
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		ui.Toast("error", "Can't delete transports", "No transports selected").Render(r.Context(), w)
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("can't parse transport id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, int32(id))
	}

	if err := h.DB.BulkSoftDeleteTransports(r.Context(), ids); err != nil {
		slog.Error("can't delete transports batch", "error", err)
	}

	h.GetTransports(w, r)
}

// ============================================================================
// Additional Handlers (optional)
// ============================================================================

func (h Handler) NewTransportPageHandler(w http.ResponseWriter, r *http.Request) {
	ui.TransportsAddContent(ui.Form{}).Render(r.Context(), w)
}

func (h Handler) EditTransportPageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Can't get transport data", "Something went wrong").Render(r.Context(), w)
		return
	}
	transport, err := h.DB.GetTransportByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't retrieve transport", "error", err, "id", id)
		ui.Toast("error", "Can't get transport data", "Not found").Render(r.Context(), w)
		return
	}
	ui.TransportsViewSheetContent(transport, ui.Form{}).Render(r.Context(), w)
}
