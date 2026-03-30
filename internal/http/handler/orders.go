package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/tms"
	"github.com/s-588/tms/internal/ui"
	"github.com/shopspring/decimal"
)

// ============================================================================
// Page & Table Handlers
// ============================================================================

func (h Handler) GetOrdersPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseOrderFilters(r)

	orders, total, err := h.DB.GetOrders(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of orders", "error", err)
		ui.Toast("error", "Can't render orders page", "Something went wrong").Render(r.Context(), w)
		return
	}

	ctx := addListsToContext(r.Context(), h.DB)
	ctx = context.WithValue(ctx, ui.FilterKey, filter)

	ui.OrdersPage(orders, page, total).Render(ctx, w)
}

func (h Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseOrderFilters(r)

	orders, total, err := h.DB.GetOrders(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of orders", "error", err)
		ui.Toast("error", "Can't get orders data", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("retrieve orders from database", "filter", filter, "page", page,
		"total pages", total, "total orders", len(orders))
	ui.OrdersTable(orders, page, total, true).Render(r.Context(), w)
}

func parseOrderFilters(r *http.Request) models.OrderFilter {
	filter := models.OrderFilter{}
	q := r.URL.Query()

	if clientID := q.Get("client_id"); clientID != "" {
		if val, err := strconv.Atoi(clientID); err == nil && val > 0 {
			filter.ClientID.SetValue(int32(val))
		}
	}
	if transportID := q.Get("transport_id"); transportID != "" {
		if val, err := strconv.Atoi(transportID); err == nil && val > 0 {
			filter.TransportID.SetValue(int32(val))
		}
	}
	if employeeID := q.Get("employee_id"); employeeID != "" {
		if val, err := strconv.Atoi(employeeID); err == nil && val > 0 {
			filter.EmployeeID.SetValue(int32(val))
		}
	}
	if priceID := q.Get("price_id"); priceID != "" {
		if val, err := strconv.Atoi(priceID); err == nil && val > 0 {
			filter.PriceID.SetValue(int32(val))
		}
	}
	if distanceMin := q.Get("distance_min"); distanceMin != "" {
		if val, err := strconv.ParseFloat(distanceMin, 10); err == nil && val > 0 {
			filter.DistanceMin.SetValue(val)
		}
	}
	if distanceMax := q.Get("distance_max"); distanceMax != "" {
		if val, err := strconv.ParseFloat(distanceMax, 10); err == nil && val > 0 {
			filter.DistanceMax.SetValue(val)
		}
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
	if priceMin := q.Get("price_min"); priceMin != "" {
		if d, err := decimal.NewFromString(priceMin); err == nil && d.IsPositive() {
			filter.TotalPriceMin.SetValue(d)
		}
	}
	if priceMax := q.Get("price_max"); priceMax != "" {
		if d, err := decimal.NewFromString(priceMax); err == nil && d.IsPositive() {
			filter.TotalPriceMax.SetValue(d)
		}
	}
	if gradeMin := q.Get("grade_min"); gradeMin != "" {
		if val, err := strconv.Atoi(gradeMin); err == nil && val >= 0 && val <= 5 {
			filter.GradeMin.SetValue(uint8(val))
		}
	}
	if gradeMax := q.Get("grade_max"); gradeMax != "" {
		if val, err := strconv.Atoi(gradeMax); err == nil && val >= 0 && val <= 5 {
			filter.GradeMax.SetValue(uint8(val))
		}
	}
	if status := q.Get("status"); status != "" {
		if err := checkOrderStatus(status); err == nil {
			filter.Status.SetValue(models.OrderStatus(status))
		}
	}
	if sortBy := q.Get("sort"); sortBy != "" {
		filter.SortBy.SetValue(sortBy)
	} else {
		filter.SortBy.SetValue("order_id")
	}
	if sortOrder := q.Get("order"); sortOrder != "" {
		filter.SortOrder.SetValue(sortOrder)
	} else {
		filter.SortOrder.SetValue("desc")
	}
	return filter
}

func (h Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't parse form", "Invalid form data").Render(r.Context(), w)
		return
	}

	hasError, form := parseOrderCreateForm(r)
	ctx := addListsToContext(r.Context(), h.DB)

	if hasError != nil {
		slog.Debug("incorrect input data for adding order", "data", form)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersAddContent(true).Render(ctx, w)
		return
	}

	// Convert form values
	clientID, _ := strconv.ParseInt(form["client_id"].Value, 10, 32)
	transportID, _ := strconv.ParseInt(form["transport_id"].Value, 10, 32)
	employeeID, _ := strconv.ParseInt(form["employee_id"].Value, 10, 32)
	priceID, _ := strconv.ParseInt(form["price_id"].Value, 10, 32)
	weight, _ := strconv.ParseInt(form["weight"].Value, 10, 32)
	nodeStart, _ := strconv.ParseInt(form["node_start"].Value, 10, 32)
	nodeEnd, _ := strconv.ParseInt(form["node_end"].Value, 10, 32)

	// Set default values for fields no longer in form
	grade := uint8(0)
	distance, err := h.DB.CalculateDistance(r.Context(), int32(nodeStart), int32(nodeEnd)) // will be calculated server‑side
	if err != nil {
		slog.Error("can't calculate distance", "error", err)
		ui.Toast("error", "Can't create order", "Something went wrong").Render(ctx, w)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersAddContent(true).Render(ctx, w)
		return
	}

	transport, err := h.DB.GetTransportByID(ctx, int32(transportID))
	if err != nil {
		slog.Error("can't get transport", "error", err)
		ui.Toast("error", "Can't create order", "Something went wrong").Render(ctx, w)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersAddContent(true).Render(ctx, w)
		return
	}

	if transport.PayloadCapacity < int32(weight) {
		form["transport_id"] = ui.FormField{
			Value: strconv.FormatInt(transportID, 10),
			Err:   fmt.Errorf("cannot ship more than vehicle can borrow"),
		}
		ui.Toast("error", "Can't create order", "Cannot ship more than vehicle can borrow").Render(ctx, w)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersAddContent(true).Render(ctx, w)
		return
	}

	employee, err := h.DB.GetEmployeeByID(ctx, int32(employeeID))
	if err != nil {
		slog.Error("can't get employee", "error", err)
		ui.Toast("error", "Can't create order", "Something went wrong").Render(ctx, w)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersAddContent(true).Render(ctx, w)
		return
	}

	if employee.Status != models.EmployeeStatusAvailable || employee.JobTitle != models.EmployeeJobTitleDriver {
		form["employee_id"] = ui.FormField{
			Value: strconv.FormatInt(employeeID, 10),
			Err:   fmt.Errorf("unavailable clients cannot be assigned"),
		}
		ui.Toast("error", "Can't create order", "Unavailable clients cannot be assigned").Render(ctx, w)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersAddContent(true).Render(ctx, w)
		return
	}

	totalPrice, err := tms.CalculateOrderCost(ctx, h.DB, tms.CalculateOrderCostArgs{
		ClientID:        int32(clientID),
		PriceID:         int32(priceID),
		Weight:          weight,
		FuelConsumption: transport.FuelConsumption,
		PayloadCapacity: transport.PayloadCapacity,
		NodeStartID:     int32(nodeStart),
		NodeEndID:       int32(nodeEnd),
	})

	order, err := h.DB.CreateOrder(ctx, db.CreateOrderArg{
		ClientID:    int32(clientID),
		TransportID: int32(transportID),
		EmployeeID:  int32(employeeID),
		Grade:       grade,
		Distance:    distance,
		Weight:      int32(weight),
		TotalPrice:  totalPrice,
		PriceID:     int32(priceID),
		Status:      models.OrderStatus(form["status"].Value),
		NodeIDStart: int32(nodeStart),
		NodeIDEnd:   int32(nodeEnd),
	})
	if err != nil {
		slog.Error("can't create order", "error", err)
		slog.Debug("can't create order",
			"ClientID", int32(clientID),
			"TransportID", int32(transportID),
			"EmployeeID", int32(employeeID),
			"Grade", grade,
			"Distance", distance,
			"Weight", int32(weight),
			"TotalPrice", totalPrice,
			"PriceID", int32(priceID),
			"Status", models.OrderStatus(form["status"].Value),
			"NodeIDStart", int32(nodeStart),
			"NodeIDEnd", int32(nodeEnd))
		ui.Toast("error", "Can't create order", "Something went wrong").Render(ctx, w)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersAddContent(true).Render(ctx, w)
		return
	}
	slog.Info("new order created", "order", order)

	ui.Toast("success", "Order created", "Order successfully created").Render(ctx, w)
	ui.OrdersAddContent(true).Render(ctx, w)
	h.GetOrders(w, r)
}

func parseOrderCreateForm(r *http.Request) (err error, form ui.Form) {
	form = make(ui.Form)

	// ClientID
	clientID := strings.TrimSpace(r.PostForm.Get("client_id"))
	form["client_id"] = ui.FormField{Value: clientID}
	if val, e := strconv.Atoi(clientID); e != nil || val <= 0 {
		err = errors.New("client must be selected")
		form["client_id"] = ui.FormField{Value: clientID, Err: err}
	}

	// TransportID
	transportID := strings.TrimSpace(r.PostForm.Get("transport_id"))
	form["transport_id"] = ui.FormField{Value: transportID}
	if val, e := strconv.Atoi(transportID); e != nil || val <= 0 {
		err = errors.New("transport must be selected")
		form["transport_id"] = ui.FormField{Value: transportID, Err: err}
	}

	// EmployeeID
	employeeID := strings.TrimSpace(r.PostForm.Get("employee_id"))
	form["employee_id"] = ui.FormField{Value: employeeID}
	if val, e := strconv.Atoi(employeeID); e != nil || val <= 0 {
		err = errors.New("employee must be selected")
		form["employee_id"] = ui.FormField{Value: employeeID, Err: err}
	}

	// PriceID
	priceID := strings.TrimSpace(r.PostForm.Get("price_id"))
	form["price_id"] = ui.FormField{Value: priceID}
	if val, e := strconv.Atoi(priceID); e != nil || val <= 0 {
		err = errors.New("price configuration must be selected")
		form["price_id"] = ui.FormField{Value: priceID, Err: err}
	}

	// Weight
	weight := strings.TrimSpace(r.PostForm.Get("weight"))
	form["weight"] = ui.FormField{Value: weight}
	if val, e := strconv.Atoi(weight); e != nil || val <= 0 {
		err = errors.New("weight must be a positive integer")
		form["weight"] = ui.FormField{Value: weight, Err: err}
	}

	// Status
	status := strings.TrimSpace(r.PostForm.Get("status"))
	form["status"] = ui.FormField{Value: status}
	if e := checkOrderStatus(status); e != nil {
		err = e
		form["status"] = ui.FormField{Value: status, Err: err}
	}

	// NodeStart (required)
	nodeStart := strings.TrimSpace(r.PostForm.Get("node_start"))
	form["node_start"] = ui.FormField{Value: nodeStart}
	if val, e := strconv.Atoi(nodeStart); e != nil || val <= 0 {
		err = errors.New("start node must be selected")
		form["node_start"] = ui.FormField{Value: nodeStart, Err: err}
	}

	// NodeEnd (required)
	nodeEnd := strings.TrimSpace(r.PostForm.Get("node_end"))
	form["node_end"] = ui.FormField{Value: nodeEnd}
	if val, e := strconv.Atoi(nodeEnd); e != nil || val <= 0 {
		err = errors.New("end node must be selected")
		form["node_end"] = ui.FormField{Value: nodeEnd, Err: err}
	}

	return
}

func (h Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Can't get order data", "Something went wrong").Render(r.Context(), w)
		return
	}
	order, err := h.DB.GetOrderByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't retrieve order", "error", err, "id", id)
		ui.Toast("error", "Can't get order data", "Not found").Render(r.Context(), w)
		return
	}
	slog.Debug("retrieve order", "order", order)

	ctx := addListsToContext(r.Context(), h.DB)
	ctx = context.WithValue(ctx, ui.FormKey, ui.Form{}) // empty form

	ui.OrdersViewSheetContent(order).Render(ctx, w)
}

func (h Handler) UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Error", "Incorrect order ID").Render(r.Context(), w)
		return
	}

	existing, err := h.DB.GetOrderByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't receive order", "error", err)
		ui.Toast("error", "Internal error", "Something went wrong").Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		ui.Toast("error", "Bad request", "Invalid form format").Render(r.Context(), w)
		return
	}

	err, form := parseOrderUpdateForm(r, existing)
	ctx := addListsToContext(r.Context(), h.DB)

	if err != nil {
		slog.Debug("can't update order", "form", form, "err", err)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		ui.OrdersViewSheetContent(existing).Render(ctx, w)
		return
	}

	// Convert form values
	clientID, _ := strconv.ParseInt(form["client_id"].Value, 10, 32)
	transportID, _ := strconv.ParseInt(form["transport_id"].Value, 10, 32)
	employeeID, _ := strconv.ParseInt(form["employee_id"].Value, 10, 32)
	priceID, _ := strconv.ParseInt(form["price_id"].Value, 10, 32)
	weight, _ := strconv.ParseInt(form["weight"].Value, 10, 32)

	// Preserve existing grade and total price; distance will be recalculated server‑side
	grade := existing.Grade
	distance := existing.Distance // keep current distance until recalculated
	totalPrice := existing.TotalPrice

	if err := h.DB.UpdateOrder(ctx, db.UpdateOrderArgs{
		OrderID:     id,
		ClientID:    int32(clientID),
		TransportID: int32(transportID),
		EmployeeID:  int32(employeeID),
		Grade:       grade,
		Distance:    distance,
		Weight:      int32(weight),
		TotalPrice:  totalPrice,
		PriceID:     int32(priceID),
		Status:      models.OrderStatus(form["status"].Value),
	}); err != nil {
		slog.Error("can't update order", "error", err, "id", id)
		ui.Toast("error", "Internal error", "something went wrong").Render(ctx, w)
		ui.OrdersViewSheetContent(existing).Render(ctx, w)
		return
	}

	slog.Debug("update order", "form data", form)
	ui.Toast("success", "Order updated", "Order successfully updated").Render(ctx, w)
	h.GetOrderHandler(w, r)
	h.GetOrders(w, r)
}

func parseOrderUpdateForm(r *http.Request, existing models.Order) (err error, form ui.Form) {
	form = make(ui.Form)

	// Helper to get value or default
	getValue := func(key string) string {
		return r.PostForm.Get(key)
	}

	// ClientID
	clientID := getValue("client_id")
	form["client_id"] = ui.FormField{Value: clientID}
	if val, e := strconv.Atoi(clientID); e != nil || val <= 0 {
		err = errors.New("client must be selected")
		form["client_id"] = ui.FormField{Value: clientID, Err: err}
	}

	// TransportID
	transportID := getValue("transport_id")
	form["transport_id"] = ui.FormField{Value: transportID}
	if val, e := strconv.Atoi(transportID); e != nil || val <= 0 {
		err = errors.New("transport must be selected")
		form["transport_id"] = ui.FormField{Value: transportID, Err: err}
	}

	// EmployeeID
	employeeID := getValue("employee_id")
	form["employee_id"] = ui.FormField{Value: employeeID}
	if val, e := strconv.Atoi(employeeID); e != nil || val <= 0 {
		err = errors.New("employee must be selected")
		form["employee_id"] = ui.FormField{Value: employeeID, Err: err}
	}

	// PriceID
	priceID := getValue("price_id")
	form["price_id"] = ui.FormField{Value: priceID}
	if val, e := strconv.Atoi(priceID); e != nil || val <= 0 {
		err = errors.New("price configuration must be selected")
		form["price_id"] = ui.FormField{Value: priceID, Err: err}
	}

	// Weight
	weight := getValue("weight")
	form["weight"] = ui.FormField{Value: weight}
	if val, e := strconv.Atoi(weight); e != nil || val <= 0 {
		err = errors.New("weight must be a positive integer")
		form["weight"] = ui.FormField{Value: weight, Err: err}
	}

	// Status
	status := getValue("status")
	form["status"] = ui.FormField{Value: status}
	if e := checkOrderStatus(status); e != nil {
		err = e
		form["status"] = ui.FormField{Value: status, Err: err}
	}

	// NodeStart (required)
	nodeStart := getValue("node_start")
	form["node_start"] = ui.FormField{Value: nodeStart}
	if val, e := strconv.Atoi(nodeStart); e != nil || val <= 0 {
		err = errors.New("start node must be selected")
		form["node_start"] = ui.FormField{Value: nodeStart, Err: err}
	}

	// NodeEnd (required)
	nodeEnd := getValue("node_end")
	form["node_end"] = ui.FormField{Value: nodeEnd}
	if val, e := strconv.Atoi(nodeEnd); e != nil || val <= 0 {
		err = errors.New("end node must be selected")
		form["node_end"] = ui.FormField{Value: nodeEnd, Err: err}
	}

	return
}

// ============================================================================
// Delete
// ============================================================================

func (h Handler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Incorrect URL", "Can't parse id from URL path").Render(r.Context(), w)
		return
	}

	if err := h.DB.SoftDeleteOrder(r.Context(), int32(id)); err != nil {
		slog.Error("can't delete order", "error", err, "id", id)
		ui.Toast("error", "Can't delete order", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("deleting order", "orderID", id)
	ui.Toast("success", "Deleted", "Order successfully deleted").Render(r.Context(), w)
	h.GetOrders(w, r)
}

// ============================================================================
// Bulk Delete
// ============================================================================

func (h Handler) BulkDeleteOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't delete orders", "Can't parse form").Render(r.Context(), w)
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		ui.Toast("error", "Can't delete orders", "No orders selected").Render(r.Context(), w)
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("can't parse order id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, int32(id))
	}

	if err := h.DB.BulkSoftDeleteOrders(r.Context(), ids); err != nil {
		slog.Error("can't delete orders batch", "error", err)
	}

	h.GetOrders(w, r)
}

// ============================================================================
// Additional Handlers (optional)
// ============================================================================

func (h Handler) GetOrderTransportsHandler(w http.ResponseWriter, r *http.Request) {
	// To be implemented if needed
}

func (h Handler) AssignOrderTransportsHandler(w http.ResponseWriter, r *http.Request) {
	// To be implemented if needed
}

// ============================================================================
// Validation Helpers
// ============================================================================

func checkOrderStatus(status string) error {
	switch models.OrderStatus(status) {
	case models.OrderStatusPending, models.OrderStatusAssigned,
		models.OrderStatusInProgress, models.OrderStatusCompleted,
		models.OrderStatusCancelled:
		return nil
	default:
		return errors.New("invalid order status")
	}
}

func int32Ptr(i int32) *int32 {
	return &i
}

func (h Handler) GetOrderAddForm(w http.ResponseWriter, r *http.Request) {
	ctx := addListsToContext(r.Context(), h.DB)
	ctx = context.WithValue(ctx, ui.FormKey, ui.Form{})

	ui.OrdersAddContent(true).Render(ctx, w)
}

func (h Handler) GetOrderFilterForm(w http.ResponseWriter, r *http.Request) {
	ctx := addListsToContext(r.Context(), h.DB)
	filter := parseOrderFilters(r)
	ctx = context.WithValue(ctx, ui.FilterKey, filter)

	ui.OrdersFilter().Render(ctx, w)
}

// addListsToContext fetches all reference lists and stores them in the context.
func addListsToContext(ctx context.Context, db db.DB) context.Context {
	clients, _ := db.ListClients(ctx)
	employees, _ := db.ListFreeDrivers(ctx)
	transports, _ := db.ListFreeTransports(ctx)
	prices, _ := db.ListPrices(ctx)
	nodes, _ := db.ListNodes(ctx)

	ctx = context.WithValue(ctx, ui.ClientsKey, clients)
	ctx = context.WithValue(ctx, ui.EmployeesKey, employees)
	ctx = context.WithValue(ctx, ui.TransportsKey, transports)
	ctx = context.WithValue(ctx, ui.PricesKey, prices)
	ctx = context.WithValue(ctx, ui.NodesKey, nodes)

	return ctx
}

func (h *Handler) OrdersExport(w http.ResponseWriter, r *http.Request) {
	filter := parseOrderFilters(r)

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	var period ReportPeriod
	if startStr != "" {
		period.Start, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr != "" {
		period.End, _ = time.Parse("2006-01-02", endStr)
	}

	// Fetch all orders matching the filter and date range
	orders, err := h.GetOrdersForReport(r.Context(), filter, period)
	if err != nil {
		slog.Error("failed to fetch orders for export", "error", err)
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	stats := aggregateStats(orders, period)

	bytes, err := GenerateOrdersReport(stats)
	if err != nil {
		slog.Error("failed to generate report", "error", err)
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="orders_report.xlsx"`)
	w.Write(bytes)
}

// fetchAllOrders — получает ВСЕ заказы (обходит пагинацию)
func (h *Handler) fetchAllOrders(ctx context.Context, filter models.OrderFilter) ([]models.Order, error) {
	var all []models.Order
	page := int32(1)

	for {
		orders, totalPages, err := h.DB.GetOrders(ctx, page, filter)
		if err != nil {
			return nil, err
		}
		all = append(all, orders...)

		if len(orders) == 0 || page >= totalPages {
			break
		}
		page++
	}
	return all, nil
}

// GetOrdersForReport — все заказы + фильтр по периоду (в памяти)
// GetOrdersForReport — все заказы + фильтр по периоду (в памяти)
func (h *Handler) GetOrdersForReport(ctx context.Context, filter models.OrderFilter, period ReportPeriod) ([]models.Order, error) {
	all, err := h.fetchAllOrders(ctx, filter)
	if err != nil {
		return nil, err
	}

	var filtered []models.Order
	for _, o := range all {
		// Lower bound: skip if start is set and order is before it
		if !period.Start.IsZero() && o.CreatedAt.Before(period.Start) {
			continue
		}
		// Upper bound: skip if end is set and order is after the end day
		if !period.End.IsZero() {
			endOfPeriod := period.End.Add(24 * time.Hour) // include the entire end day
			if !o.CreatedAt.Before(endOfPeriod) {
				continue
			}
		}
		filtered = append(filtered, o)
	}
	return filtered, nil
}
// aggregateStats — основная функция агрегации
func aggregateStats(orders []models.Order, period ReportPeriod) OrderStats {
	stats := OrderStats{
		Period:         period,
		Orders:         orders,
		TotalOrders:    len(orders),
		OrdersByStatus: make(map[models.OrderStatus]int),
		OrdersByMonth:  make(map[string]int),
	}

	if len(orders) == 0 {
		return stats
	}

	var totalRevenue decimal.Decimal
	var totalDistance, totalWeight float64
	clientMap := make(map[string]struct {
		Revenue decimal.Decimal
		Count   int
	})

	for _, o := range orders {
		totalRevenue = totalRevenue.Add(o.TotalPrice)
		totalDistance += o.Distance
		totalWeight += float64(o.Weight)

		stats.OrdersByStatus[o.Status]++

		monthKey := o.CreatedAt.Format("2006-01")
		stats.OrdersByMonth[monthKey]++

		c := clientMap[o.ClientName]
		c.Revenue = c.Revenue.Add(o.TotalPrice)
		c.Count++
		clientMap[o.ClientName] = c
	}

	stats.TotalRevenue = totalRevenue
	if stats.TotalOrders > 0 {
		stats.AvgPrice = totalRevenue.Div(decimal.NewFromInt(int64(stats.TotalOrders))).Round(2)
		stats.AvgDistance = totalDistance / float64(stats.TotalOrders)
		stats.AvgWeight = totalWeight / float64(stats.TotalOrders)
	}

	// Топ-10 клиентов
	for name, v := range clientMap {
		stats.TopClients = append(stats.TopClients, TopClient{
			Name:    name,
			Revenue: v.Revenue,
			Count:   v.Count,
		})
	}
	sort.Slice(stats.TopClients, func(i, j int) bool {
		return stats.TopClients[i].Revenue.GreaterThan(stats.TopClients[j].Revenue)
	})
	if len(stats.TopClients) > 10 {
		stats.TopClients = stats.TopClients[:10]
	}

	return stats
}

