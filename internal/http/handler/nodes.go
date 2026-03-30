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

func (h Handler) GetNodesPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseNodeFilters(r)

	nodes, total, err := h.DB.GetNodes(r.Context(), 1, models.NodeFilter{})
	if err != nil {
		slog.Error("can't retrieve list of nodes", "error", err)
		ui.Toast("error", "Can't render nodes page", "Something went wrong").Render(r.Context(), w)
		return
	}
	ui.NodesPage(nodes, page, total, filter).Render(r.Context(), w)
}

func (h Handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseNodeFilters(r)

	nodes, total, err := h.DB.GetNodes(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of nodes", "error", err)
		ui.Toast("error", "Can't get nodes data", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("retrieve nodes from database", "filter", filter, "page", page,
		"total pages", total, "total nodes", len(nodes))
	ui.NodesTable(nodes, page, total, filter, true).Render(r.Context(), w)
}

// ============================================================================
// Filter Parsing
// ============================================================================

func parseNodeFilters(r *http.Request) models.NodeFilter {
	filter := models.NodeFilter{}
	q := r.URL.Query()

	if name := q.Get("name"); name != "" {
		filter.Name.SetValue(name)
	}
	if sortBy := q.Get("sort"); sortBy != "" {
		filter.SortBy.SetValue(sortBy)
	} else {
		filter.SortBy.SetValue("node_id")
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

func (h Handler) CreateNodeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't parse form", "Invalid form data").Render(r.Context(), w)
		return
	}

	hasError, form := parseNodeCreateForm(r)
	if hasError != nil {
		slog.Debug("incorrect input data for adding node", "data", form)
		ui.NodesAddContent(form).Render(r.Context(), w)
		return
	}

	// Convert coordinates
	x, _ := strconv.ParseFloat(form["x"].Value, 64)
	y, _ := strconv.ParseFloat(form["y"].Value, 64)

	_, err := h.DB.CreateNode(r.Context(),
		db.CreateNodeArgs{
			Name:    models.Optional[string]{Value: form["name"].Value},
			Address: form["address"].Value,
			Geom:    models.Point{X: x, Y: y},
		})
	if err != nil {
		if errors.Is(err, db.ErrDuplicateNodeAddress) {
			// We can attach the error to any field, or create a general error message
			// Let's attach it to cargo_type for simplicity
			form["address"] = ui.FormField{Value: form["address"].Value, Err: errors.New("address already exists")}
			ui.NodesAddContent(form).Render(r.Context(), w)
			return
		}
		slog.Error("can't create node", "error", err)
		ui.Toast("error", "Can't create node", "Something went wrong").Render(r.Context(), w)
		ui.NodesAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("adding new node", "data", form)
	ui.Toast("success", "Node created", "Node successfully created").Render(r.Context(), w)
	h.GetNodes(w, r)
}

func parseNodeCreateForm(r *http.Request) (err error, form ui.Form) {
	form = make(ui.Form)

	// Name
	name := strings.TrimSpace(r.PostForm.Get("name"))
	form["name"] = ui.FormField{Value: name}

	address := strings.TrimSpace(r.PostForm.Get("address"))
	form["address"] = ui.FormField{Value: name}
	if address == "" {
		err = errors.New("address is required")
		form["address"] = ui.FormField{Value: name, Err: err}
	}

	// X coordinate
	xStr := strings.TrimSpace(r.PostForm.Get("x"))
	form["x"] = ui.FormField{Value: xStr}
	if xStr == "" {
		err = errors.New("X coordinate is required")
		form["x"] = ui.FormField{Value: xStr, Err: err}
	} else {
		if _, e := strconv.ParseFloat(xStr, 64); e != nil {
			err = errors.New("X coordinate must be a valid number")
			form["x"] = ui.FormField{Value: xStr, Err: err}
		}
	}

	// Y coordinate
	yStr := strings.TrimSpace(r.PostForm.Get("y"))
	form["y"] = ui.FormField{Value: yStr}
	if yStr == "" {
		err = errors.New("Y coordinate is required")
		form["y"] = ui.FormField{Value: yStr, Err: err}
	} else {
		if _, e := strconv.ParseFloat(yStr, 64); e != nil {
			err = errors.New("Y coordinate must be a valid number")
			form["y"] = ui.FormField{Value: yStr, Err: err}
		}
	}

	return
}

// ============================================================================
// Read (single node for sheet)
// ============================================================================

func (h Handler) GetNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Can't get node data", "Something went wrong").Render(r.Context(), w)
		return
	}
	node, err := h.DB.GetNodeByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't retrieve node", "error", err, "id", id)
		ui.Toast("error", "Can't get node data", "Not found").Render(r.Context(), w)
		return
	}
	slog.Debug("retrieve node", "node", node)
	ui.NodesViewSheetContent(node, ui.Form{}).Render(r.Context(), w)
}

// ============================================================================
// Update
// ============================================================================

func (h Handler) UpdateNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Error", "Incorrect node ID").Render(r.Context(), w)
		return
	}

	existing, err := h.DB.GetNodeByID(r.Context(), int32(id))
	if err != nil {
		slog.Error("can't receive node", "error", err)
		ui.Toast("error", "Internal error", "Something went wrong").Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		ui.Toast("error", "Bad request", "Invalid form format").Render(r.Context(), w)
		return
	}

	err, form := parseNodeUpdateForm(r, existing)
	if err != nil {
		slog.Debug("can't update node", "form", form, "err", err)
		ui.NodesViewSheetContent(existing, form).Render(r.Context(), w)
		return
	}

	x, _ := strconv.ParseFloat(form["x"].Value, 64)
	y, _ := strconv.ParseFloat(form["y"].Value, 64)

	nameOpt := models.Optional[string]{}
	if nameVal := form["name"].Value; nameVal != "" {
		nameOpt.SetValue(nameVal)
	}

	if err := h.DB.UpdateNode(r.Context(), db.UpdateNodeArgs{
		NodeID:  int32(id),
		Name:    nameOpt,
		Geom:    models.Point{X: x, Y: y},
		Address: form["address"].Value, // from updated form
	}); err != nil {
		if errors.Is(err, db.ErrDuplicateNodeAddress) {
			// We can attach the error to any field, or create a general error message
			// Let's attach it to cargo_type for simplicity
			form["address"] = ui.FormField{Value: form["address"].Value, Err: errors.New("address already exists")}
			ui.NodesAddContent(form).Render(r.Context(), w)
			return
		}
		slog.Error("can't update node", "error", err)
		ui.Toast("error", "Can't update node", "Something went wrong").Render(r.Context(), w)
		ui.NodesAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("update node", "form data", form)
	ui.Toast("success", "Node updated", "Node successfully updated").Render(r.Context(), w)
	h.GetNodeHandler(w, r)
	h.GetNodes(w, r)
}

func parseNodeUpdateForm(r *http.Request, existing models.Node) (err error, form ui.Form) {
	form = make(ui.Form)

	getValue := func(key string, defaultValue string) string {
		if val := r.PostForm.Get(key); val != "" {
			return val
		}
		return defaultValue
	}

	// Name
	name := getValue("name", existing.Name)
	form["name"] = ui.FormField{Value: name}
	if name == "" {
		err = errors.New("name is required")
		form["name"] = ui.FormField{Value: name, Err: err}
	}

	// X coordinate
	xStr := getValue("x", strconv.FormatFloat(existing.Geom.X, 'f', -1, 64))
	form["x"] = ui.FormField{Value: xStr}
	if xStr == "" {
		err = errors.New("X coordinate is required")
		form["x"] = ui.FormField{Value: xStr, Err: err}
	} else {
		if _, e := strconv.ParseFloat(xStr, 64); e != nil {
			err = errors.New("X coordinate must be a valid number")
			form["x"] = ui.FormField{Value: xStr, Err: err}
		}
	}

	// Y coordinate
	yStr := getValue("y", strconv.FormatFloat(existing.Geom.Y, 'f', -1, 64))
	form["y"] = ui.FormField{Value: yStr}
	if yStr == "" {
		err = errors.New("Y coordinate is required")
		form["y"] = ui.FormField{Value: yStr, Err: err}
	} else {
		if _, e := strconv.ParseFloat(yStr, 64); e != nil {
			err = errors.New("Y coordinate must be a valid number")
			form["y"] = ui.FormField{Value: yStr, Err: err}
		}
	}

	address := getValue("address", existing.Address)
	form["address"] = ui.FormField{Value: address}
	if address == "" {
		err = errors.New("address is required")
		form["address"] = ui.FormField{Value: address, Err: err}
	}

	return
}

// ============================================================================
// Delete
// ============================================================================

func (h Handler) DeleteNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Incorrect URL", "Can't parse id from URL path").Render(r.Context(), w)
		return
	}

	if err := h.DB.SoftDeleteNode(r.Context(), int32(id)); err != nil {
		slog.Error("can't delete node", "error", err, "id", id)
		ui.Toast("error", "Can't delete node", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("deleting node", "nodeID", id)
	ui.Toast("success", "Deleted", "Node successfully deleted").Render(r.Context(), w)
	h.GetNodes(w, r)
}

// ============================================================================
// Bulk Delete
// ============================================================================

func (h Handler) BulkDeleteNodesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't delete nodes", "Can't parse form").Render(r.Context(), w)
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		ui.Toast("error", "Can't delete nodes", "No nodes selected").Render(r.Context(), w)
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("can't parse node id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, int32(id))
	}

	if err := h.DB.BulkSoftDeleteNodes(r.Context(), ids); err != nil {
		slog.Error("can't delete nodes batch", "error", err)
	}

	h.GetNodes(w, r)
}
