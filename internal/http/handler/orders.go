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
		renderToast(w, r, "error", "Can't render orders page", "Something went wrong")
		return
	}

	ctx := addListsToContext(r.Context(), h.DB)
	ctx = context.WithValue(ctx, ui.FilterKey, filter)

	err = ui.OrdersPage(orders, page, total).Render(ctx, w)
	if err != nil {
		slog.Error("can't render orders page", "error", err)
	}
}

func (h Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseOrderFilters(r)

	orders, total, err := h.DB.GetOrders(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of orders", "error", err)
		renderToast(w, r, "error", "Can't get orders data", "Something went wrong")
		return
	}

	slog.Debug("retrieve orders from database", "filter", filter, "page", page,
		"total pages", total, "total orders", len(orders))
	err = ui.OrdersTable(orders, page, total, true).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render orders table", "error", err)
	}
}

//nolint:funlen,gocyclo // function is clear and understandable
func parseOrderFilters(r *http.Request) models.OrderFilter {
	filter := models.OrderFilter{}
	q := r.URL.Query()

	if clientID := q.Get("client_id"); clientID != "" {
		if val, err := parseID(clientID); err == nil {
			filter.ClientID.SetValue(val)
		}
	}
	if transportID := q.Get("transport_id"); transportID != "" {
		if val, err := parseID(transportID); err == nil {
			filter.TransportID.SetValue(val)
		}
	}
	if employeeID := q.Get("employee_id"); employeeID != "" {
		if val, err := parseID(employeeID); err == nil {
			filter.EmployeeID.SetValue(val)
		}
	}
	if priceID := q.Get("price_id"); priceID != "" {
		if val, err := parseID(priceID); err == nil {
			filter.PriceID.SetValue(val)
		}
	}
	if distanceMin := q.Get("distance_min"); distanceMin != "" {
		if val, err := strconv.ParseFloat(distanceMin, 64); err == nil && val > 0 {
			filter.DistanceMin.SetValue(val)
		}
	}
	if distanceMax := q.Get("distance_max"); distanceMax != "" {
		if val, err := strconv.ParseFloat(distanceMax, 64); err == nil && val > 0 {
			filter.DistanceMax.SetValue(val)
		}
	}
	if weightMin := q.Get("weight_min"); weightMin != "" {
		if val, err := strconv.ParseInt(weightMin, 10, 32); err == nil && val > 0 {
			filter.WeightMin.SetValue(int32(val))
		}
	}
	if weightMax := q.Get("weight_max"); weightMax != "" {
		if val, err := strconv.ParseInt(weightMax, 10, 32); err == nil && val > 0 {
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

func parseID(clientID string) (int32, error) {
	val, err := strconv.ParseInt(clientID, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("can't parse int: %w", err)
	} else if val < 0 {
		return 0, fmt.Errorf("id less than zero")
	}
	return int32(val), nil
}

func (h Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't parse form", "Invalid form data")
		return
	}

	form, args, err := parseOrderCreateForm(r)
	ctx := addListsToContext(r.Context(), h.DB)
	if err != nil {
		h.renderOrderFormError(w, ctx, form, "", "") // no toast, just re-render form
		return
	}

	args.Grade = 0

	distance, err := h.DB.CalculateDistance(r.Context(), args.NodeIDStart, args.NodeIDEnd)
	if err != nil {
		slog.Error("can't calculate distance", "error", err)
		h.renderOrderFormError(w, ctx, form, "Can't create order", "Something went wrong")
		return
	}
	args.Distance = distance

	transport, err := h.DB.GetTransportByID(ctx, args.TransportID)
	if err != nil {
		slog.Error("can't get transport", "error", err)
		h.renderOrderFormError(w, ctx, form, "Can't create order", "Something went wrong")
		return
	}

	if transport.PayloadCapacity < args.Weight {
		form["transport_id"] = ui.FormField{
			Value: strconv.FormatInt(int64(args.TransportID), 10),
			Err:   fmt.Errorf("cannot ship more than vehicle can borrow"),
		}
		h.renderOrderFormError(w, ctx, form, "Can't create order", "Cannot ship more than vehicle can borrow")
		return
	}

	employee, err := h.DB.GetEmployeeByID(ctx, args.EmployeeID)
	if err != nil {
		slog.Error("can't get employee", "error", err)
		h.renderOrderFormError(w, ctx, form, "Can't create order", "Something went wrong")
		return
	}

	if employee.Status != models.EmployeeStatusAvailable || employee.JobTitle != models.EmployeeJobTitleDriver {
		form["employee_id"] = ui.FormField{
			Value: strconv.FormatInt(int64(args.EmployeeID), 10),
			Err:   fmt.Errorf("unavailable clients cannot be assigned"),
		}
		h.renderOrderFormError(w, ctx, form, "Can't create order", "Unavailable clients cannot be assigned")
		return
	}

	totalPrice, err := tms.CalculateOrderCost(ctx, h.DB, tms.CalculateOrderCostArgs{
		ClientID:        args.ClientID,
		PriceID:         args.PriceID,
		Weight:          int64(args.Weight),
		FuelConsumption: transport.FuelConsumption,
		PayloadCapacity: transport.PayloadCapacity,
		NodeStartID:     args.NodeIDStart,
		NodeEndID:       args.NodeIDEnd,
	})
	if err != nil {
		slog.Error("can't calculate order cost", "error", err)
		h.renderOrderFormError(w, ctx, form, "Can't create order", "Something went wrong")
		return
	}
	args.TotalPrice = totalPrice

	order, err := h.DB.CreateOrder(ctx, args)
	if err != nil {
		slog.Error("can't create order", "error", err)
		h.renderOrderFormError(w, ctx, form, "Can't create order", "Something went wrong")
		return
	}
	slog.Info("new order created", "order", order)

	if err := ui.Toast("success", "Order created", "Order successfully created").Render(ctx, w); err != nil {
		slog.Error("can't render toast", "error", err)
	}
	if err := ui.OrdersAddContent(true).Render(ctx, w); err != nil {
		slog.Error("can't render orders add content", "error", err)
	}
	h.GetOrders(w, r)
}

func (h Handler) renderOrderFormError(w http.ResponseWriter, ctx context.Context, form ui.Form, toastTitle, toastMsg string) {
	if toastTitle != "" {
		if err := ui.Toast("error", toastTitle, toastMsg).Render(ctx, w); err != nil {
			slog.Error("can't render toast", "error", err)
		}
	}
	ctx = context.WithValue(ctx, ui.FormKey, form)
	if err := ui.OrdersAddContent(true).Render(ctx, w); err != nil {
		slog.Error("can't render orders add content", "error", err)
	}
}

func parseOrderCreateForm(r *http.Request) (form ui.Form, args db.CreateOrderArg, err error) {
	form = make(ui.Form)

	setStr := func(key, val string, check func(string) error) {
		form[key] = ui.FormField{Value: val}
		if e := check(val); e != nil {
			form[key] = ui.FormField{Value: val, Err: e}
			err = e
		}
	}

	setInt := func(key, val string, dest *int32, msg string) {
		form[key] = ui.FormField{Value: val}
		n, e := strconv.ParseInt(val, 10, 32)
		if e != nil || n <= 0 {
			e = errors.New(msg)
			form[key] = ui.FormField{Value: val, Err: e}
			err = e
			return
		}
		*dest = int32(n)
	}

	// ClientID
	setInt("client_id", strings.TrimSpace(r.PostForm.Get("client_id")), &args.ClientID, "client must be selected")

	// TransportID
	setInt("transport_id", strings.TrimSpace(r.PostForm.Get("transport_id")), &args.TransportID, "transport must be selected")

	// EmployeeID
	setInt("employee_id", strings.TrimSpace(r.PostForm.Get("employee_id")), &args.EmployeeID, "employee must be selected")

	// PriceID
	setInt("price_id", strings.TrimSpace(r.PostForm.Get("price_id")), &args.PriceID, "price configuration must be selected")

	// Weight
	setInt("weight", strings.TrimSpace(r.PostForm.Get("weight")), &args.Weight, "weight must be a positive integer")

	// Status
	status := strings.TrimSpace(r.PostForm.Get("status"))
	setStr("status", status, checkOrderStatus)
	args.Status = models.OrderStatus(status)

	// NodeStart
	setInt("node_start", strings.TrimSpace(r.PostForm.Get("node_start")), &args.NodeIDStart, "start node must be selected")

	// NodeEnd
	setInt("node_end", strings.TrimSpace(r.PostForm.Get("node_end")), &args.NodeIDEnd, "end node must be selected")

	if err != nil {
		return form, args, err
	}

	// Grade, Distance, TotalPrice are filled later in the handler
	return form, args, nil
}

func (h Handler) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Can't get order data", "Something went wrong")
		return
	}
	order, err := h.DB.GetOrderByID(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve order", "error", err, "id", id)
		renderToast(w, r, "error", "Can't get order data", "Not found")
		return
	}
	slog.Debug("retrieve order", "order", order)

	ctx := addListsToContext(r.Context(), h.DB)
	ctx = context.WithValue(ctx, ui.FormKey, ui.Form{}) // empty form

	err = ui.OrdersViewSheetContent(order).Render(ctx, w)
	if err != nil {
		slog.Error("can't render orders view sheet content", "error", err)
	}
}

func (h Handler) UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Error", "Incorrect order ID")
		return
	}

	existing, err := h.DB.GetOrderByID(r.Context(), id)
	if err != nil {
		slog.Error("can't receive order", "error", err)
		renderToast(w, r, "error", "Internal error", "Something went wrong")
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		renderToast(w, r, "error", "Bad request", "Invalid form format")
		return
	}

	form, args, err := parseOrderUpdateForm(r, id, existing)
	ctx := addListsToContext(r.Context(), h.DB)
	if err != nil {
		slog.Debug("can't update order", "form", form, "err", err)
		ctx = context.WithValue(ctx, ui.FormKey, form)
		err = ui.OrdersViewSheetContent(existing).Render(ctx, w)
		if err != nil {
			slog.Error("can't render orders view sheet content", "error", err)
		}
		return
	}

	if err := h.DB.UpdateOrder(ctx, args); err != nil {
		slog.Error("can't update order", "error", err, "id", id)
		err = ui.Toast("error", "Internal error", "something went wrong").Render(ctx, w)
		if err != nil {
			slog.Error("can't render toast", "error", err)
		}
		err = ui.OrdersViewSheetContent(existing).Render(ctx, w)
		if err != nil {
			slog.Error("can't render orders view sheet content", "error", err)
		}
		return
	}

	slog.Debug("update order", "form data", form)
	err = ui.Toast("success", "Order updated", "Order successfully updated").Render(ctx, w)
	if err != nil {
		slog.Error("can't render toast", "error", err)
	}
	h.GetOrderHandler(w, r)
	h.GetOrders(w, r)
}

func parseOrderUpdateForm(r *http.Request, id int32, existing models.Order) (form ui.Form, args db.UpdateOrderArgs, err error) {
	form = make(ui.Form)

	setStr := func(key, val string, check func(string) error) {
		form[key] = ui.FormField{Value: val}
		if e := check(val); e != nil {
			form[key] = ui.FormField{Value: val, Err: e}
			err = e
		}
	}

	setInt := func(key, val string, dest *int32, msg string) {
		form[key] = ui.FormField{Value: val}
		n, e := strconv.ParseInt(val, 10, 32)
		if e != nil || n <= 0 {
			e = errors.New(msg)
			form[key] = ui.FormField{Value: val, Err: e}
			err = e
			return
		}
		*dest = int32(n)
	}

	// ClientID
	setInt("client_id", r.PostForm.Get("client_id"), &args.ClientID, "client must be selected")

	// TransportID
	setInt("transport_id", r.PostForm.Get("transport_id"), &args.TransportID, "transport must be selected")

	// EmployeeID
	setInt("employee_id", r.PostForm.Get("employee_id"), &args.EmployeeID, "employee must be selected")

	// PriceID
	setInt("price_id", r.PostForm.Get("price_id"), &args.PriceID, "price configuration must be selected")

	// Weight
	setInt("weight", r.PostForm.Get("weight"), &args.Weight, "weight must be a positive integer")

	// Status
	status := r.PostForm.Get("status")
	setStr("status", status, checkOrderStatus)
	args.Status = models.OrderStatus(status)

	// NodeStart (validated only)
	setInt("node_start", r.PostForm.Get("node_start"), new(int32), "start node must be selected")

	// NodeEnd (validated only)
	setInt("node_end", r.PostForm.Get("node_end"), new(int32), "end node must be selected")

	if err != nil {
		return form, args, err
	}

	// Preserve existing grade and total price; distance will be recalculated server-side
	args.OrderID = id
	args.Grade = existing.Grade
	args.Distance = existing.Distance
	args.TotalPrice = existing.TotalPrice

	return form, args, nil
}

func (h Handler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Incorrect URL", "Can't parse id from URL path")
		return
	}

	if err := h.DB.SoftDeleteOrder(r.Context(), id); err != nil {
		slog.Error("can't delete order", "error", err, "id", id)
		renderToast(w, r, "error", "Can't delete order", "Something went wrong")
		return
	}

	slog.Debug("deleting order", "orderID", id)
	renderToast(w, r, "success", "Deleted", "Order successfully deleted")
	h.GetOrders(w, r)
}

func (h Handler) BulkDeleteOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't delete orders", "Can't parse form")
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		err := ui.Toast("error", "Can't delete orders", "No orders selected").Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render toast", "error", err)
		}
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.ParseInt(idStr, 10, 32)
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

func (h Handler) GetOrderTransportsHandler(w http.ResponseWriter, r *http.Request) {
	// To be implemented if needed
}

func (h Handler) AssignOrderTransportsHandler(w http.ResponseWriter, r *http.Request) {
	// To be implemented if needed
}

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

func (h Handler) GetOrderAddForm(w http.ResponseWriter, r *http.Request) {
	ctx := addListsToContext(r.Context(), h.DB)
	ctx = context.WithValue(ctx, ui.FormKey, ui.Form{})

	err := ui.OrdersAddContent(true).Render(ctx, w)
	if err != nil {
		slog.Error("can't render order add form", "error", err)
	}
}

func (h Handler) GetOrderFilterForm(w http.ResponseWriter, r *http.Request) {
	ctx := addListsToContext(r.Context(), h.DB)
	filter := parseOrderFilters(r)
	ctx = context.WithValue(ctx, ui.FilterKey, filter)

	err := ui.OrdersFilter().Render(ctx, w)
	if err != nil {
		slog.Error("can't render order filter form", "error", err)
	}
}

// addListsToContext fetches all reference lists and stores them in the context.
func addListsToContext(ctx context.Context, db db.DB) context.Context {
	clients, err := db.ListClients(ctx)
	if err != nil {
		slog.Error("can't fetch clients", "error", err)
	}
	employees, err := db.ListFreeDrivers(ctx)
	if err != nil {
		slog.Error("can't fetch employees", "error", err)
	}
	transports, err := db.ListFreeTransports(ctx)
	if err != nil {
		slog.Error("can't fetch transports", "error", err)
	}
	prices, err := db.ListPrices(ctx)
	if err != nil {
		slog.Error("can't fetch prices", "error", err)
	}
	nodes, err := db.ListNodes(ctx)
	if err != nil {
		slog.Error("can't fetch nodes", "error", err)
	}

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
	_, err = w.Write(bytes)
	if err != nil {
		slog.Error("can't write report to response", "error", err)
	}
}

// fetchAllOrders returns all orders without pagination
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
