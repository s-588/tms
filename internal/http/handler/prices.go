package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/s-588/tms/cmd/models"
)

// parsePriceFilter extracts models.PriceFilter from request query.
func parsePriceFilter(r *http.Request) models.PriceFilter {
	var filter models.PriceFilter

	// CargoType – take first non-empty value from comma-split or direct param
	cargoTypes := strings.Split(r.URL.Query().Get("cargo_type"), ",")
	for _, ct := range cargoTypes {
		if ct != "" {
			filter.CargoType.SetValue(ct)
			break
		}
	}

	// Weight
	if w := r.URL.Query().Get("weight_min"); w != "" {
		if val, err := strconv.Atoi(w); err == nil && val >= 0 {
			filter.WeightMin.SetValue(int32(val))
		}
	}
	if w := r.URL.Query().Get("weight_max"); w != "" {
		if val, err := strconv.Atoi(w); err == nil && val >= 0 {
			filter.WeightMax.SetValue(int32(val))
		}
	}

	// Distance
	if d := r.URL.Query().Get("distance_min"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val >= 0 {
			filter.DistanceMin.SetValue(int32(val))
		}
	}
	if d := r.URL.Query().Get("distance_max"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val >= 0 {
			filter.DistanceMax.SetValue(int32(val))
		}
	}

	filter.CreatedFrom = parseOptionalTime(r.URL.Query().Get("created_from"))
	filter.CreatedTo = parseOptionalTime(r.URL.Query().Get("created_to"))
	filter.UpdatedFrom = parseOptionalTime(r.URL.Query().Get("updated_from"))
	filter.UpdatedTo = parseOptionalTime(r.URL.Query().Get("updated_to"))

	filter.SortBy = parseOptionalString(r.URL.Query().Get("sort"))
	filter.SortOrder = parseOptionalString(r.URL.Query().Get("order"))

	return filter
}

func (h Handler) GetPricesHandler(w http.ResponseWriter, r *http.Request) {
	// limit, offset := parsePagination(r)
	filter := parsePriceFilter(r)

	if !filter.SortBy.Set {
		filter.SortBy.SetValue("price_id")
	}
	if !filter.SortOrder.Set {
		filter.SortOrder.SetValue("desc")
	}

	// prices, total, err := h.DB.GetPrices(r.Context(), limit, offset, filter)
	// if err != nil {
	// 	slog.Error("can't retrieve prices", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }

	// ui.PricesTable(prices, limit, offset, int(total), filter).Render(r.Context(), w)
}

func (h Handler) GetPriceHandler(w http.ResponseWriter, r *http.Request) {
	// id, err := parseIDFromReq(r)
	// if err != nil {
	// 	slog.Error("can't parse id", "error", err)
	// 	responseError(w, r, http.StatusBadRequest, "incorrect price id")
	// 	return
	// }
	// price, err := h.DB.GetPriceByID(r.Context(), id)
	// if err != nil {
	// 	slog.Error("can't retrieve price", "error", err, "id", id)
	// 	responseError(w, r, http.StatusNotFound, "price not found")
	// 	return
	// }
	// ui.PriceDetail(price).Render(r.Context(), w)
}

func (h Handler) CreatePriceHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged validation logic
}

func (h Handler) DeletePriceHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged
}

func (h Handler) UpdatePriceHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged validation logic
}

func (h Handler) NewPricePageHandler(w http.ResponseWriter, r *http.Request) {
	// ui.PriceCreateForm(nil).Render(r.Context(), w)
}

func (h Handler) EditPricePageHandler(w http.ResponseWriter, r *http.Request) {
	// id, err := parseIDFromReq(r)
	// if err != nil {
	// 	slog.Error("can't parse id", "error", err)
	// 	responseError(w, r, http.StatusBadRequest, "incorrect price id")
	// 	return
	// }
	// price, err := h.DB.GetPriceByID(r.Context(), id)
	// if err != nil {
	// 	slog.Error("can't retrieve price for edit", "error", err, "id", id)
	// 	responseError(w, r, http.StatusNotFound, "price not found")
	// 	return
	// }
	// ui.PriceEditForm(price, nil).Render(r.Context(), w)
}

func (h Handler) ExportPricesHandler(w http.ResponseWriter, r *http.Request) {
	// filter := parsePriceFilter(r)
	//
	// prices, err := h.DB.GetPrices(r.Context(), filter)
	// if err != nil {
	// 	slog.Error("can't retrieve prices for export", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }
	//
	// w.Header().Set("Content-Type", "text/csv")
	// w.Header().Set("Content-Disposition", "attachment;filename=prices.csv")
	//
	// fmt.Fprintln(w, "ID,Cargo Type,Cost,Weight,Distance,Created At,Updated At,Deleted At")
	// for _, price := range prices {
	// 	deletedAt := ""
	// 	if !price.DeletedAt.IsZero() {
	// 		deletedAt = price.DeletedAt.Format(time.RFC3339)
	// 	}
	// 	fmt.Fprintf(w, "%d,%s,%s,%d,%d,%s,%s,%s\n",
	// 		price.PriceID,
	// 		escapeCSV(price.CargoType),
	// 		price.Cost,
	// 		price.Weight,
	// 		price.Distance,
	// 		price.CreatedAt.Format(time.RFC3339),
	// 		price.UpdatedAt.Format(time.RFC3339),
	// 		deletedAt,
	// 	)
	// }
}
