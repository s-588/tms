package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

// parseOrderFilter extracts models.OrderFilter from request query.
func parseOrderFilter(r *http.Request) models.OrderFilter {
	var filter models.OrderFilter

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

	// Total price (string)
	var tpMin, tpMax models.Optional[decimal.Decimal]
	if r.URL.Query().Has("price_min") {
		if d, err := decimal.NewFromString(r.URL.Query().Get("price_min")); err == nil {
			tpMin.SetValue(d)
		}
	}
	if r.URL.Query().Has("price_max") {
		if d, err := decimal.NewFromString(r.URL.Query().Get("price_max")); err == nil {
			tpMax.SetValue(d)
		}
	}

	// Status – take first non-empty value
	statuses := r.URL.Query()["status"]
	for _, s := range statuses {
		if s != "" {
			// TODO: add order status check
			filter.Status.SetValue(models.OrderStatus(s))
			break
		}
	}

	// Timestamps
	filter.CreatedFrom = parseOptionalTime(r.URL.Query().Get("created_from"))
	filter.CreatedTo = parseOptionalTime(r.URL.Query().Get("created_to"))
	filter.UpdatedFrom = parseOptionalTime(r.URL.Query().Get("updated_from"))
	filter.UpdatedTo = parseOptionalTime(r.URL.Query().Get("updated_to"))

	// Sorting
	filter.SortBy = parseOptionalString(r.URL.Query().Get("sort"))
	filter.SortOrder = parseOptionalString(r.URL.Query().Get("order"))

	return filter
}

func (h Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// limit, offset := parsePagination(r)
	filter := parseOrderFilter(r)

	// Default sort
	if !filter.SortBy.Set {
		filter.SortBy.SetValue("order_id")
	}
	if !filter.SortOrder.Set {
		filter.SortOrder.SetValue("desc")
	}

	// Get transports for filter dropdown (unfiltered)
	// transports, _, err := h.DB.GetTransports(r.Context(), limit, offset, models.TransportFilter{})
	// if err != nil {
	// 	slog.Error("can't retrieve transports for filter", "error", err)
	// }
	//
	// orders, total, err := h.DB.GetOrders(r.Context(), limit, offset, filter)
	// if err != nil {
	// 	slog.Error("can't retrieve orders", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }
	//
	// ui.OrdersTable(orders, limit, offset, int(total), filter, transports).Render(r.Context(), w)
}

func (h Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id", "error", err)
		responseError(w, r, http.StatusBadRequest, "incorrect order id")
		return
	}
	// order, err := h.DB.GetOrderByID(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve order", "error", err, "id", id)
		responseError(w, r, http.StatusNotFound, "order not found")
		return
	}
	// ui.OrderDetail(order).Render(r.Context(), w)
}

func (h Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and validate the form
	params, errs := ParseOrderCreateForm(r)
	if len(errs) > 0 {
		// Re‑render the creation form with errors
		// ui.OrderCreateForm(errs).Render(r.Context(), w)
		return
	}

	// Create the order using all required fields
	order, err := h.DB.CreateOrder(r.Context(),
		params.ClientID,
		params.TransportID,
		params.EmployeeID,
		params.Grade,
		params.Distance,
		params.Weight,
		params.TotalPrice,
		params.PriceID,
		params.Status,
		params.NodeIDStart,
		params.NodeIDEnd,
	)
	if err != nil {
		slog.Error("can't create order", "error", err)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	ui.CreateSuccess("Order added", "orders", int(order.OrderID)).Render(r.Context(), w)
}

func calculateTotalPrice(d, w int, s models.OrderStatus) decimal.Decimal {
	// TODO: add real logic
	return decimal.Decimal{}
}

// OrderCreateParams holds all data needed to create an Order.
// TODO: delete and replace with database layer struct
type OrderCreateParams struct {
	ClientID    int32
	TransportID int32
	EmployeeID  int32
	Grade       uint8
	Distance    int32
	Weight      int32
	PriceID     int32
	Status      models.OrderStatus
	NodeIDStart int32
	NodeIDEnd   int32
	TotalPrice  decimal.Decimal
}

// ParseOrderCreateForm reads the POST form data, validates all fields,
// calculates the total price, and returns an OrderCreateParams struct
// along with a map of field‑specific validation errors.
func ParseOrderCreateForm(r *http.Request) (OrderCreateParams, map[string]string) {
	var params OrderCreateParams
	errs := make(map[string]string)

	// Ensure form is parsed
	if err := r.ParseForm(); err != nil {
		errs["form"] = "invalid form data"
		return params, errs
	}

	// Helper to parse int32 fields
	parseInt32Field := func(fieldName string) (int32, bool) {
		valStr := r.PostForm.Get(fieldName)
		if valStr == "" {
			errs[fieldName] = fieldName + " is required"
			return 0, false
		}
		val, err := strconv.Atoi(valStr)
		if err != nil || val <= 0 {
			errs[fieldName] = fieldName + " must be a positive integer"
			return 0, false
		}
		return int32(val), true
	}

	// ClientID
	if clientID, ok := parseInt32Field("client_id"); ok {
		params.ClientID = clientID
	}

	// TransportID
	if transportID, ok := parseInt32Field("transport_id"); ok {
		params.TransportID = transportID
	}

	// EmployeeID
	if employeeID, ok := parseInt32Field("employee_id"); ok {
		params.EmployeeID = employeeID
	}

	// Grade (0-5 typically)
	gradeStr := r.PostForm.Get("grade")
	if gradeStr == "" {
		errs["grade"] = "grade is required"
	} else {
		grade, err := strconv.Atoi(gradeStr)
		if err != nil || grade < 0 || grade > 5 {
			errs["grade"] = "grade must be an integer between 0 and 5"
		} else {
			params.Grade = uint8(grade)
		}
	}

	// Distance
	if distance, ok := parseInt32Field("distance"); ok {
		params.Distance = distance
	}

	// Weight
	if weight, ok := parseInt32Field("weight"); ok {
		params.Weight = weight
	}

	// PriceID
	if priceID, ok := parseInt32Field("price_id"); ok {
		params.PriceID = priceID
	}

	// Status
	statusStr := r.PostForm.Get("status")
	if statusStr == "" {
		errs["status"] = "status is required"
	} else {
		status := models.OrderStatus(statusStr)
		switch status {
		case models.OrderStatusPending, models.OrderStatusAssigned,
			models.OrderStatusInProgress, models.OrderStatusCompleted,
			models.OrderStatusCancelled:
			params.Status = status
		default:
			errs["status"] = "invalid order status"
		}
	}

	// NodeIDStart
	if nodeStart, ok := parseInt32Field("node_start"); ok {
		params.NodeIDStart = nodeStart
	}

	// NodeIDEnd
	if nodeEnd, ok := parseInt32Field("node_end"); ok {
		params.NodeIDEnd = nodeEnd
	}

	// If there are errors so far, return early (no price calculation needed)
	if len(errs) > 0 {
		return params, errs
	}

	// Calculate total price using the provided function.
	// The function signature from the original handler: calculateTotalPrice(distance, weight, status) string
	// We need to adapt it to return decimal.Decimal and an error.
	totalPrice := calculateTotalPrice(int(params.Distance), int(params.Weight), params.Status)
	// if err != nil {
	// 	errs["total_price"] = "unable to calculate price: " + err.Error()
	// 	return params, errs
	// }
	params.TotalPrice = totalPrice

	return params, errs
}
func (h Handler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id", "error", err)
		responseError(w, r, http.StatusBadRequest, "incorrect order id")
		return
	}
	if err := h.DB.SoftDeleteOrder(r.Context(), id); err != nil {
		slog.Error("can't delete order", "error", err, "id", id)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged, but ensure it uses the correct wrapper methods
}

func (h Handler) GetOrderTransportsHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged
}

func (h Handler) AssignOrderTransportsHandler(w http.ResponseWriter, r *http.Request) {
	// ... unchanged
}

func (h Handler) NewOrderPageHandler(w http.ResponseWriter, r *http.Request) {
	// ui.OrderCreateForm(nil).Render(r.Context(), w)
}

func (h Handler) EditOrderPageHandler(w http.ResponseWriter, r *http.Request) {
	// id, err := parseIDFromReq(r)
	// if err != nil {
	// 	slog.Error("can't parse id", "error", err)
	// 	responseError(w, r, http.StatusBadRequest, "incorrect order id")
	// 	return
	// }
	// order, err := h.DB.GetOrderByID(r.Context(), id)
	// if err != nil {
	// 	slog.Error("can't retrieve order for edit", "error", err, "id", id)
	// 	responseError(w, r, http.StatusNotFound, "order not found")
	// 	return
	// }
	// ui.OrderEditForm(order, nil).Render(r.Context(), w)
}

func (h Handler) ExportOrdersHandler(w http.ResponseWriter, r *http.Request) {
	// filter := parseOrderFilter(r)
	//
	// orders, err := h.DB.GetOrders(r.Context(), filter)
	// if err != nil {
	// 	slog.Error("can't retrieve orders for export", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }
	//
	// w.Header().Set("Content-Type", "text/csv")
	// w.Header().Set("Content-Disposition", "attachment;filename=orders.csv")
	//
	// fmt.Fprintln(w, "ID,Distance,Weight,Total Price,Status,Created At,Updated At,Deleted At")
	// for _, order := range orders {
	// 	deletedAt := ""
	// 	if !order.DeletedAt.IsZero() {
	// 		deletedAt = order.DeletedAt.Format(time.RFC3339)
	// 	}
	// 	fmt.Fprintf(w, "%d,%d,%d,%s,%s,%s,%s,%s\n",
	// 		order.OrderID,
	// 		order.Distance,
	// 		order.Weight,
	// 		order.TotalPrice,
	// 		order.Status,
	// 		order.CreatedAt.Format(time.RFC3339),
	// 		order.UpdatedAt.Format(time.RFC3339),
	// 		deletedAt,
	// 	)
	// }
}
