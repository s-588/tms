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
	"github.com/shopspring/decimal"
)

// ============================================================================
// Page & Table Handlers
// ============================================================================

func (h Handler) GetPricesPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parsePriceFilters(r)

	// For initial page load, fetch with default filter (page 1, no filters)
	prices, total, err := h.DB.GetPrices(r.Context(), 1, models.PriceFilter{})
	if err != nil {
		slog.Error("can't retrieve list of prices", "error", err)
		ui.Toast("error", "Can't render prices page", "Something went wrong").Render(r.Context(), w)
		return
	}
	ui.PricesPage(prices, page, total, filter).Render(r.Context(), w)
}

func (h Handler) GetPrices(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parsePriceFilters(r)

	prices, total, err := h.DB.GetPrices(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of prices", "error", err)
		ui.Toast("error", "Can't get prices data", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("retrieve prices from database", "filter", filter, "page", page,
		"total pages", total, "total prices", len(prices))
	ui.PricesTable(prices, page, total, filter, true).Render(r.Context(), w)
}

// ============================================================================
// Filter Parsing
// ============================================================================

func parsePriceFilters(r *http.Request) models.PriceFilter {
	filter := models.PriceFilter{}
	q := r.URL.Query()

	if cargo := q.Get("cargo_type"); cargo != "" {
		filter.CargoType.SetValue(cargo)
	}
	if weightMin := q.Get("weight_min"); weightMin != "" {
		if val, err := strconv.Atoi(weightMin); err == nil && val > 0 {
			filter.WeightMin.SetValue(int32(val))
		}
	}
	if weightMax := q.Get("weight_max"); weightMax != "" {
		if val, err := strconv.Atoi(weightMax); err == nil && val > 0 {
			filter.WeightMax.SetValue(int32(val))
		}
	}
	if distMin := q.Get("distance_min"); distMin != "" {
		if val, err := strconv.Atoi(distMin); err == nil && val > 0 {
			filter.DistanceMin.SetValue(int32(val))
		}
	}
	if distMax := q.Get("distance_max"); distMax != "" {
		if val, err := strconv.Atoi(distMax); err == nil && val > 0 {
			filter.DistanceMax.SetValue(int32(val))
		}
	}
	if sortBy := q.Get("sort"); sortBy != "" {
		filter.SortBy.SetValue(sortBy)
	} else {
		filter.SortBy.SetValue("price_id")
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

func (h Handler) CreatePriceHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't parse form", "Invalid form data").Render(r.Context(), w)
		return
	}

	hasError, form := parsePriceCreateForm(r)
	if hasError != nil {
		slog.Debug("incorrect input data for adding price", "data", form)
		ui.PricesAddContent(form).Render(r.Context(), w)
		return
	}

	// Convert form values
	weight, _ := decimal.NewFromString(form["weight"].Value)
	distance, _ := decimal.NewFromString(form["distance"].Value)

	_, err := h.DB.CreatePrice(r.Context(), db.CreatePriceArgs{
		CargoType: form["cargo_type"].Value,
		Weight:    weight,
		Distance:  distance,
	})
	if err != nil {
		if errors.Is(err, db.ErrDuplicatePrice) {
			// We can attach the error to any field, or create a general error message
			// Let's attach it to cargo_type for simplicity
			form["cargo_type"] = ui.FormField{Value: form["cargo_type"].Value, Err: errors.New("price configuration already exists")}
			ui.PricesAddContent(form).Render(r.Context(), w)
			return
		}
		slog.Error("can't create price", "error", err)
		ui.Toast("error", "Can't create price", "Something went wrong").Render(r.Context(), w)
		ui.PricesAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("adding new price", "data", form)
	ui.Toast("success", "Price created", "Price configuration successfully created").Render(r.Context(), w)
	h.GetPrices(w, r) // refresh table
}

func parsePriceCreateForm(r *http.Request) (err error, form ui.Form) {
	form = make(ui.Form)

	// Cargo Type
	cargo := strings.TrimSpace(r.PostForm.Get("cargo_type"))
	form["cargo_type"] = ui.FormField{Value: cargo}
	if cargo == "" {
		err = errors.New("cargo type is required")
		form["cargo_type"] = ui.FormField{Value: cargo, Err: err}
	}

	// Weight
	weightStr := strings.TrimSpace(r.PostForm.Get("weight"))
	form["weight"] = ui.FormField{Value: weightStr}
	if weightStr == "" {
		err = errors.New("weight coefficient is required")
		form["weight"] = ui.FormField{Value: weightStr, Err: err}
	} else {
		weight, e := decimal.NewFromString(weightStr)
		if e != nil || weight.IsNegative() || weight.IsZero() {
			err = errors.New("weight coefficient must be a positive decimal")
			form["weight"] = ui.FormField{Value: weightStr, Err: err}
		}
	}

	// Distance
	distanceStr := strings.TrimSpace(r.PostForm.Get("distance"))
	form["distance"] = ui.FormField{Value: distanceStr}
	if distanceStr == "" {
		err = errors.New("distance coefficient is required")
		form["distance"] = ui.FormField{Value: distanceStr, Err: err}
	} else {
		distance, e := decimal.NewFromString(distanceStr)
		if e != nil || distance.IsNegative() || distance.IsZero() {
			err = errors.New("distance coefficient must be a positive decimal")
			form["distance"] = ui.FormField{Value: distanceStr, Err: err}
		}
	}

	return
}

// ============================================================================
// Read (single price for sheet)
// ============================================================================

func (h Handler) GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Can't get price data", "Something went wrong").Render(r.Context(), w)
		return
	}
	price, err := h.DB.GetPriceByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't retrieve price", "error", err, "id", id)
		ui.Toast("error", "Can't get price data", "Not found").Render(r.Context(), w)
		return
	}
	slog.Debug("retrieve price", "price", price)
	ui.PricesViewSheetContent(price, ui.Form{}).Render(r.Context(), w)
}

// ============================================================================
// Update
// ============================================================================

func (h Handler) UpdatePriceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Error", "Incorrect price ID").Render(r.Context(), w)
		return
	}

	existing, err := h.DB.GetPriceByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't receive price", "error", err)
		ui.Toast("error", "Internal error", "Something went wrong").Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		ui.Toast("error", "Bad request", "Invalid form format").Render(r.Context(), w)
		return
	}

	err, form := parsePriceUpdateForm(r, existing)
	if err != nil {
		slog.Debug("can't update price", "form", form, "err", err)
		ui.PricesViewSheetContent(existing, form).Render(r.Context(), w)
		return
	}

	// Convert form values
	weight, _ := decimal.NewFromString(form["weight"].Value)
	distance, _ := decimal.NewFromString(form["distance"].Value)

	if err := h.DB.UpdatePrice(r.Context(), db.UpdatePriceArgs{
		PriceID:   id,
		CargoType: form["cargo_type"].Value,
		Weight:    weight,
		Distance:  distance,
	}); err != nil {
		if errors.Is(err, db.ErrDuplicatePrice) {
			// We can attach the error to any field, or create a general error message
			// Let's attach it to cargo_type for simplicity
			form["cargo_type"] = ui.FormField{Value: form["cargo_type"].Value, Err: errors.New("price configuration already exists")}
			ui.PricesAddContent(form).Render(r.Context(), w)
			return
		}
		slog.Error("can't create price", "error", err)
		ui.Toast("error", "Can't create price", "Something went wrong").Render(r.Context(), w)
		ui.PricesAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("update price", "form data", form)
	ui.Toast("success", "Price updated", "Price configuration successfully updated").Render(r.Context(), w)
	h.GetPriceHandler(w, r) // refresh sheet
	h.GetPrices(w, r)       // refresh table
}

func parsePriceUpdateForm(r *http.Request, existing models.Price) (err error, form ui.Form) {
	form = make(ui.Form)

	// Helper to get value or default
	getValue := func(key string, defaultValue string) string {
		if val := r.PostForm.Get(key); val != "" {
			return val
		}
		return defaultValue
	}

	// Cargo Type
	cargo := getValue("cargo_type", existing.CargoType)
	form["cargo_type"] = ui.FormField{Value: cargo}
	if cargo == "" {
		err = errors.New("cargo type is required")
		form["cargo_type"] = ui.FormField{Value: cargo, Err: err}
	}

	// Weight
	weightStr := getValue("weight", existing.Weight.String())
	form["weight"] = ui.FormField{Value: weightStr}
	if weightStr == "" {
		err = errors.New("weight coefficient is required")
		form["weight"] = ui.FormField{Value: weightStr, Err: err}
	} else {
		weight, e := decimal.NewFromString(weightStr)
		if e != nil || weight.IsNegative() || weight.IsZero() {
			err = errors.New("weight coefficient must be a positive decimal")
			form["weight"] = ui.FormField{Value: weightStr, Err: err}
		}
	}

	// Distance
	distanceStr := getValue("distance", existing.Distance.String())
	form["distance"] = ui.FormField{Value: distanceStr}
	if distanceStr == "" {
		err = errors.New("distance coefficient is required")
		form["distance"] = ui.FormField{Value: distanceStr, Err: err}
	} else {
		distance, e := decimal.NewFromString(distanceStr)
		if e != nil || distance.IsNegative() || distance.IsZero() {
			err = errors.New("distance coefficient must be a positive decimal")
			form["distance"] = ui.FormField{Value: distanceStr, Err: err}
		}
	}

	return
}

// ============================================================================
// Delete
// ============================================================================

func (h Handler) DeletePriceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Incorrect URL", "Can't parse id from URL path").Render(r.Context(), w)
		return
	}

	if err := h.DB.SoftDeletePrice(r.Context(), int32(id)); err != nil {
		slog.Error("can't delete price", "error", err, "id", id)
		ui.Toast("error", "Can't delete price", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("deleting price", "priceID", id)
	ui.Toast("success", "Deleted", "Price configuration successfully deleted").Render(r.Context(), w)
	h.GetPrices(w, r)
}

// ============================================================================
// Bulk Delete
// ============================================================================

func (h Handler) BulkDeletePricesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't delete prices", "Can't parse form").Render(r.Context(), w)
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		ui.Toast("error", "Can't delete prices", "No prices selected").Render(r.Context(), w)
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("can't parse price id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, int32(id))
	}

	if err := h.DB.BulkSoftDeletePrices(r.Context(), ids); err != nil {
		slog.Error("can't delete prices batch", "error", err)
	}

	h.GetPrices(w, r)
}

// ============================================================================
// Additional Handlers (optional)
// ============================================================================

// NewPricePageHandler renders the add price form (for direct access if needed)
func (h Handler) NewPricePageHandler(w http.ResponseWriter, r *http.Request) {
	ui.PricesAddContent(ui.Form{}).Render(r.Context(), w)
}
