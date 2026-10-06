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

func (h Handler) GetNodesPage(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseNodeFilters(r)

	nodes, total, err := h.DB.GetNodes(r.Context(), 1, models.NodeFilter{})
	if err != nil {
		slog.Error("can't retrieve list of nodes", "error", err)
		renderToast(w, r, "error", "Can't render nodes page", "Something went wrong")
		return
	}
	err = ui.NodesPage(nodes, page, total, filter).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

func (h Handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)
	filter := parseNodeFilters(r)

	nodes, total, err := h.DB.GetNodes(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of nodes", "error", err)
		renderToast(w, r, "error", "Can't get nodes data", "Something went wrong")
		return
	}

	slog.Debug("retrieve nodes from database", "filter", filter, "page", page,
		"total pages", total, "total nodes", len(nodes))
	err = ui.NodesTable(nodes, page, total, filter, true).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

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

func (h Handler) CreateNodeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't parse form", "Invalid form data")
		return
	}

	form, err := parseNodeCreateForm(r)
	if err != nil {
		slog.Debug("incorrect input data for adding node", "data", form)
		err := ui.NodesAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	// Convert coordinates
	x, err := strconv.ParseFloat(form["x"].Value, 64)
	if err != nil {
		slog.Error("can't parse X coordinate", "error", err)
		renderToast(w, r, "error", "Can't create node", "Invalid X coordinate")
		return
	}
	y, err := strconv.ParseFloat(form["y"].Value, 64)
	if err != nil {
		slog.Error("can't parse Y coordinate", "error", err)
		renderToast(w, r, "error", "Can't create node", "Invalid Y coordinate")
		return
	}

	_, err = h.DB.CreateNode(r.Context(),
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
			err = ui.NodesAddContent(form).Render(r.Context(), w)
			if err != nil {
				slog.Error("can't render response", "error", err)
			}
			return
		}
		slog.Error("can't create node", "error", err)
		renderToast(w, r, "error", "Can't create node", "Something went wrong")
		err = ui.NodesAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	slog.Debug("adding new node", "data", form)
	renderToast(w, r, "success", "Node created", "Node successfully created")
	h.GetNodes(w, r)
}

func parseNodeCreateForm(r *http.Request) (form ui.Form, err error) {
	form = make(ui.Form)

	// Name
	name := strings.TrimSpace(r.PostForm.Get("name"))
	form["name"] = ui.FormField{Value: name}

	address := strings.TrimSpace(r.PostForm.Get("address"))
	form["address"] = ui.FormField{Value: address}
	if address == "" {
		err = errors.New("address is required")
		form["address"] = ui.FormField{Value: address, Err: err}
	}

	// X coordinate
	xStr := strings.TrimSpace(r.PostForm.Get("x"))
	form["x"] = ui.FormField{Value: xStr}
	if xStr == "" {
		err = errors.New("x coordinate is required")
		form["x"] = ui.FormField{Value: xStr, Err: err}
	} else {
		if _, e := strconv.ParseFloat(xStr, 64); e != nil {
			slog.Error("can't parse x coordinate", "error", e)
			err = errors.New("x coordinate must be a valid number")
			form["x"] = ui.FormField{Value: xStr, Err: err}
		}
	}

	// Y coordinate
	yStr := strings.TrimSpace(r.PostForm.Get("y"))
	form["y"] = ui.FormField{Value: yStr}
	if yStr == "" {
		err = errors.New("y coordinate is required")
		form["y"] = ui.FormField{Value: yStr, Err: err}
	} else {
		if _, e := strconv.ParseFloat(yStr, 64); e != nil {
			slog.Error("can't parse y coordinate", "error", e)
			err = errors.New("y coordinate must be a valid number")
			form["y"] = ui.FormField{Value: yStr, Err: err}
		}
	}

	return
}

func (h Handler) GetNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Can't get node data", "Something went wrong")
		return
	}
	node, err := h.DB.GetNodeByID(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve node", "error", err, "id", id)
		renderToast(w, r, "error", "Can't get node data", "Not found")
		return
	}
	slog.Debug("retrieve node", "node", node)
	err = ui.NodesViewSheetContent(node, ui.Form{}).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

func (h Handler) UpdateNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Error", "Incorrect node ID")
		return
	}

	existing, err := h.DB.GetNodeByID(r.Context(), id)
	if err != nil {
		slog.Error("can't receive node", "error", err)
		renderToast(w, r, "error", "Internal error", "Something went wrong")
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't parse http form", "error", err)
		renderToast(w, r, "error", "Bad request", "Invalid form format")
		return
	}

	form, args, err := parseNodeUpdateForm(r, id, existing)
	if err != nil {
		slog.Debug("can't update node", "form", form, "err", err)
		err = ui.NodesViewSheetContent(existing, form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	if err := h.DB.UpdateNode(r.Context(), args); err != nil {
		if errors.Is(err, db.ErrDuplicateNodeAddress) {
			form["address"] = ui.FormField{Value: form["address"].Value, Err: errors.New("address already exists")}
			err = ui.NodesAddContent(form).Render(r.Context(), w)
			if err != nil {
				slog.Error("can't render response", "error", err)
			}
			return
		}
		slog.Error("can't update node", "error", err)
		renderToast(w, r, "error", "Can't update node", "Something went wrong")
		err = ui.NodesAddContent(form).Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	slog.Debug("update node", "form data", form)
	renderToast(w, r, "success", "Node updated", "Node successfully updated")
	h.GetNodeHandler(w, r)
	h.GetNodes(w, r)
}

func parseNodeUpdateForm(r *http.Request, id int32, existing models.Node) (form ui.Form, args db.UpdateNodeArgs, err error) {
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
	var x float64
	if xStr == "" {
		err = errors.New("x coordinate is required")
		form["x"] = ui.FormField{Value: xStr, Err: err}
	} else {
		var e error
		x, e = strconv.ParseFloat(xStr, 64)
		if e != nil {
			slog.Error("can't parse x coordinate", "error", e)
			err = errors.New("x coordinate must be a valid number")
			form["x"] = ui.FormField{Value: xStr, Err: err}
		}
	}

	// Y coordinate
	yStr := getValue("y", strconv.FormatFloat(existing.Geom.Y, 'f', -1, 64))
	form["y"] = ui.FormField{Value: yStr}
	var y float64
	if yStr == "" {
		err = errors.New("y coordinate is required")
		form["y"] = ui.FormField{Value: yStr, Err: err}
	} else {
		var e error
		y, e = strconv.ParseFloat(yStr, 64)
		if e != nil {
			slog.Error("can't parse y coordinate", "error", e)
			err = errors.New("y coordinate must be a valid number")
			form["y"] = ui.FormField{Value: yStr, Err: err}
		}
	}

	// Address
	address := getValue("address", existing.Address)
	form["address"] = ui.FormField{Value: address}
	if address == "" {
		err = errors.New("address is required")
		form["address"] = ui.FormField{Value: address, Err: err}
	}

	if err != nil {
		return form, args, err
	}

	nameOpt := models.Optional[string]{}
	if name != "" {
		nameOpt.SetValue(name)
	}

	args = db.UpdateNodeArgs{
		NodeID:  id,
		Name:    nameOpt,
		Geom:    models.Point{X: x, Y: y},
		Address: address,
	}

	return form, args, nil
}

func (h Handler) DeleteNodeHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		renderToast(w, r, "error", "Incorrect URL", "Can't parse id from URL path")
		return
	}

	if err := h.DB.SoftDeleteNode(r.Context(), id); err != nil {
		slog.Error("can't delete node", "error", err, "id", id)
		renderToast(w, r, "error", "Can't delete node", "Something went wrong")
		return
	}

	slog.Debug("deleting node", "nodeID", id)
	renderToast(w, r, "success", "Deleted", "Node successfully deleted")
	h.GetNodes(w, r)
}

func (h Handler) BulkDeleteNodesHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderToast(w, r, "error", "Can't delete nodes", "Can't parse form")
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		err := ui.Toast("error", "Can't delete nodes", "No nodes selected").Render(r.Context(), w)
		if err != nil {
			slog.Error("can't render response", "error", err)
		}
		return
	}

	var ids []int32
	for _, idStr := range selectedIDs {
		id, err := strconv.ParseInt(idStr, 10, 32)
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
