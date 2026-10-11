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
		renderToast(w, r, "error", "Can't render transports page", "Something went wrong")
		return
	}
	err = ui.TransportsPage(transports, page, total, filter).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render transports page", "error", err)
		renderToast(w, r, "error", "Can't render transports page", "Something went wrong")
	}
}

func (h Handler) GetTransports(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseTransportFilters(r)

	transports, total, err := h.DB.GetTransports(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of transports", "error", err)
		renderToast(w, r, "error", "Can't get transports data", "Something went wrong")
		return
	}

	slog.Debug("retrieve transports from database", "filter", filter, "page", page,
		"total pages", total, "total transports", len(transports))
	err = ui.TransportsTable(transports, page, total, filter, true).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render transports table", "error", err)
		renderToast(w, r, "error", "Can't render transports table", "Something went wrong")
	}
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
		if val, err := strconv.ParseInt(payloadMin, 10, 32); err == nil && val > 0 {
			filter.PayloadCapacityMin.SetValue(int32(val))
		}
	}
	if payloadMax := q.Get("payload_max"); payloadMax != "" {
		if val, err := strconv.ParseInt(payloadMax, 10, 32); err == nil && val > 0 {
			filter.PayloadCapacityMax.SetValue(int32(val))
		}
	}
	if fuelMin := q.Get("fuel_min"); fuelMin != "" {
		if val, err := strconv.ParseInt(fuelMin, 10, 32); err == nil && val > 0 {
			filter.FuelConsumptionMin.SetValue(int32(val))
		}
	}
	if fuelMax := q.Get("fuel_max"); fuelMax != "" {
		if val, err := strconv.ParseInt(fuelMax, 10, 32); err == nil && val > 0 {
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

func (h Handler) CreateTransportHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't parse form", "Invalid form data")
		return
	}

	form, err := parseTransportCreateForm(r)
	if err != nil {
		slog.Debug("incorrect input data for adding transport", "data", form)
		err := ui.TransportsAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render transports add content", "error", err)
		}
		return
	}

	payload, err := strconv.ParseInt(form["payload_capacity"].Value, 10, 32)
	if err != nil {
		slog.Error("can't convert payload capacity to int", "error", err)
		renderToast(w, r, "error", "Can't create transport", "Invalid payload capacity")
		return
	}
	fuel, err := strconv.ParseInt(form["fuel_consumption"].Value, 10, 32)
	if err != nil {
		slog.Error("can't convert fuel consumption to int", "error", err)
		renderToast(w, r, "error", "Can't create transport", "Invalid fuel consumption")
		return
	}

	_, err = h.DB.CreateTransport(r.Context(), db.CreateTransportArgs{
		Model:           form["model"].Value,
		LicensePlate:    form["license_plate"].Value,
		PayloadCapacity: int32(payload),
		FuelConsumption: int32(fuel),
	})
	if err != nil {
		if errors.Is(err, db.ErrDuplicateLicense) {
			form["license_plate"] = ui.FormField{Value: form["license_plate"].Value, Err: errors.New("license plate already exists")}
			err := ui.TransportsAddContent(form).Render(r.Context(), w)
			if err != nil {
				slog.Error("can't render transports add content", "error", err)
			}
			return
		}
		slog.Error("can't create transport", "error", err)
		err := ui.Toast("error", "Can't create transport", "Something went wrong").Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render toast message", "error", err)
		}
		err = ui.TransportsAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render transports add content", "error", err)
		}
		return
	}

	slog.Debug("adding new transport", "data", form)
	renderToast(w, r, "success", "Transport created", "Transport successfully created")
	h.GetTransports(w, r)
}

func parseTransportCreateForm(r *http.Request) (form ui.Form, err error) {
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

func (h Handler) GetTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Can't get transport data", "Something went wrong")
		return
	}
	transport, err := h.DB.GetTransportByID(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve transport", "error", err, "id", id)
		renderToast(w, r, "error", "Can't get transport data", "Not found")
		return
	}
	slog.Debug("retrieve transport", "transport", transport)
	err = ui.TransportsViewSheetContent(transport, ui.Form{}).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render transports view sheet content", "error", err)
	}
}

func (h Handler) UpdateTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Error", "Incorrect transport ID")
		return
	}

	existing, err := h.DB.GetTransportByID(r.Context(), id)
	if err != nil {
		slog.Error("can't receive transport", "error", err)
		renderToast(w, r, "error", "Internal error", "Something went wrong")
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		renderToast(w, r, "error", "Bad request", "Invalid form format")
		return
	}

	form, err := parseTransportUpdateForm(r, existing)
	if err != nil {
		slog.Debug("can't update transport", "form", form, "err", err)
		err = ui.TransportsViewSheetContent(existing, form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render transports view sheet content", "error", err)
		}
		return
	}

	payload, err := strconv.ParseInt(form["payload_capacity"].Value, 10, 32)
	if err != nil {
		slog.Error("can't convert payload capacity to int", "error", err)
		renderToast(w, r, "error", "Can't update transport", "Invalid payload capacity")
		return
	}
	fuel, err := strconv.ParseInt(form["fuel_consumption"].Value, 10, 32)
	if err != nil {
		slog.Error("can't convert fuel consumption to int", "error", err)
		renderToast(w, r, "error", "Can't update transport", "Invalid fuel consumption")
		return
	}

	if err := h.DB.UpdateTransport(r.Context(), db.UpdateTransportArgs{
		TransportID:     id,
		Model:           form["model"].Value,
		LicensePlate:    form["license_plate"].Value,
		PayloadCapacity: int32(payload),
		FuelConsumption: int32(fuel),
	}); err != nil {
		if errors.Is(err, db.ErrDuplicateLicense) {
			form["license_plate"] = ui.FormField{Value: form["license_plate"].Value, Err: errors.New("license plate already exists")}
			err = ui.TransportsAddContent(form).Render(r.Context(), w)
			if err != nil {
				slog.Error("can't render transports view sheet content", "error", err)
			}
			return
		}
		slog.Error("can't create transport", "error", err)
		renderToast(w, r, "error", "Can't create transport", "Something went wrong")
		err = ui.TransportsAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render transports view sheet content", "error", err)
		}
		return
	}

	slog.Debug("update transport", "form data", form)
	renderToast(w, r, "success", "Transport updated", "Transport successfully updated")
	h.GetTransportHandler(w, r)
	h.GetTransports(w, r)
}

func parseTransportUpdateForm(r *http.Request, existing models.Transport) (form ui.Form, err error) {
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

func (h Handler) DeleteTransportHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Incorrect URL", "Can't parse id from URL path")
		return
	}

	if err := h.DB.SoftDeleteTransport(r.Context(), id); err != nil {
		slog.Error("can't delete transport", "error", err, "id", id)
		renderToast(w, r, "error", "Can't delete transport", "Something went wrong")
		return
	}

	slog.Debug("deleting transport", "transportID", id)
	renderToast(w, r, "success", "Deleted", "Transport successfully deleted")
	h.GetTransports(w, r)
}

func (h Handler) BulkDeleteTransportsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't delete transports", "Can't parse form")
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		err := ui.Toast("error", "Can't delete transports", "No transports selected").Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render toast message", "error", err)
		}
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.ParseInt(idStr, 10, 32)
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

func (h Handler) NewTransportPageHandler(w http.ResponseWriter, r *http.Request) {
	err := ui.TransportsAddContent(ui.Form{}).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render transports add content", "error", err)
	}
}

func (h Handler) EditTransportPageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Can't get transport data", "Something went wrong")
		return
	}
	transport, err := h.DB.GetTransportByID(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve transport", "error", err, "id", id)
		renderToast(w, r, "error", "Can't get transport data", "Not found")
		return
	}
	err = ui.TransportsViewSheetContent(transport, ui.Form{}).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render transports view sheet content", "error", err)
	}
}
