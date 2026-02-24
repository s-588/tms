package handler

import (
	"net/http"

	"github.com/s-588/tms/cmd/models"
)

// parseTransportFilter extracts models.TransportFilter from request query.
func parseTransportFilter(r *http.Request) models.TransportFilter {
	var filter models.TransportFilter

	filter.Model = parseOptionalString(r.URL.Query().Get("model"))
	filter.LicensePlate = parseOptionalString(r.URL.Query().Get("license_plate"))

	filter.PayloadCapacityMin = parseOptionalInt32(r.URL.Query().Get("payload_min"))
	filter.PayloadCapacityMax = parseOptionalInt32(r.URL.Query().Get("payload_max"))
	filter.FuelConsumptionMin = parseOptionalInt32(r.URL.Query().Get("fuel_consumption_min"))
	filter.FuelConsumptionMax = parseOptionalInt32(r.URL.Query().Get("fuel_consumption_max"))

	filter.CreatedFrom = parseOptionalTime(r.URL.Query().Get("created_from"))
	filter.CreatedTo = parseOptionalTime(r.URL.Query().Get("created_to"))
	filter.UpdatedFrom = parseOptionalTime(r.URL.Query().Get("updated_from"))
	filter.UpdatedTo = parseOptionalTime(r.URL.Query().Get("updated_to"))

	filter.SortBy = parseOptionalString(r.URL.Query().Get("sort"))
	filter.SortOrder = parseOptionalString(r.URL.Query().Get("order"))

	return filter
}

func (h Handler) GetTransportsHandler(w http.ResponseWriter, r *http.Request) {
	// limit, offset := parsePagination(r)
	filter := parseTransportFilter(r)

	// Default sorting
	if !filter.SortBy.Set {
		filter.SortBy.SetValue("transport_id")
	}
	if !filter.SortOrder.Set {
		filter.SortOrder.SetValue("desc")
	}

	// transports, total, err := h.DB.GetTransports(r.Context(), limit, offset, filter)
	// if err != nil {
	// 	slog.Error("can't retrieve transports", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }
	//
	// ui.TransportsTable(transports, limit, offset, int(total), filter).Render(r.Context(), w)
}

func (h Handler) GetTransportHandler(w http.ResponseWriter, r *http.Request) {
	// id, err := parseIDFromReq(r)
	// if err != nil {
	// 	slog.Error("can't parse id", "error", err)
	// 	responseError(w, r, http.StatusBadRequest, "incorrect transport id")
	// 	return
	// }
	// transport, err := h.DB.GetTransportByID(r.Context(), id)
	// if err != nil {
	// 	slog.Error("can't retrieve transport", "error", err, "id", id)
	// 	responseError(w, r, http.StatusNotFound, "transport not found")
	// 	return
	// }
	// ui.TransportDetail(transport).Render(r.Context(), w)
}

func (h Handler) CreateTransportHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged validation logic (it uses GetEmployeeByID, GetFuelByID, etc.)
	// Ensure it uses the correct wrapper methods if available.
}

func (h Handler) DeleteTransportHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged
}

func (h Handler) UpdateTransportHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged
}

func (h Handler) NewTransportPageHandler(w http.ResponseWriter, r *http.Request) {
	// ui.TransportCreateForm(nil).Render(r.Context(), w)
}

func (h Handler) EditTransportPageHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged
}
