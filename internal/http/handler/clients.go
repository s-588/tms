package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/ui"
)

var (
	phoneRegex = regexp.MustCompile(`^[+]?[0-9\\-\\s()]{7,25}$`)
)

// GetClientsPage just retrieves clients from database and render clients page.
func (h Handler) GetClientsPage(w http.ResponseWriter, r *http.Request) {
	clients, total, err := h.DB.GetClients(r.Context(), 1, models.ClientFilter{})
	if err != nil {
		slog.Error("can't retrieve list of clients", "error", err)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}
	ui.ClientsPage(clients, 1, total, models.ClientFilter{}).Render(r.Context(), w)
}

// GetClients parses id from path; page and models.ClientFilter from
// path values, retrieve clients from database using filters and pagination
// and return ClientsTable form.
func (h Handler) GetClients(w http.ResponseWriter, r *http.Request) {
	page := parsePagination(r)

	filter := parseClientFilters(r)

	clients, total, err := h.DB.GetClients(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of clients", "error", err)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	slog.Debug("retrieve clients from database", "filter", filter, "page", page,
		"total pages", total, "total clients", len(clients))
	ui.ClientsTable(clients, page, total, filter).Render(r.Context(), w)
}

// parseClientFilters function parse models.ClientFilter from request url values
// and return it.
func parseClientFilters(r *http.Request) models.ClientFilter {
	filter := models.ClientFilter{}
	q := r.URL.Query()

	if q.Has("name") {
		filter.Name.SetValue(q.Get("name"))
	}
	if q.Has("email") {
		filter.Email.SetValue(q.Get("email"))
	}
	if q.Has("phone") {
		filter.Phone.SetValue(q.Get("phone"))
	}
	if q.Has("email_verified") {
		emailVerifiedStr := q.Get("email_verified")
		if emailVerifiedStr == "true" {
			filter.EmailVerified.SetValue(true)
		} else if emailVerifiedStr == "false" {
			filter.EmailVerified.SetValue(false)
		}
	}
	if q.Has("created_from") {
		if t, err := time.Parse("2006-01-02", q.Get("created_from")); err == nil {
			filter.CreatedFrom.SetValue(t)
		}
	}
	if q.Has("created_to") {
		if t, err := time.Parse("2006-01-02", q.Get("created_to")); err == nil {
			filter.CreatedTo.SetValue(t)
		}
	}
	if q.Has("updated_from") {
		if t, err := time.Parse("2006-01-02", q.Get("updated_from")); err == nil {
			filter.UpdatedFrom.SetValue(t)
		}
	}
	if q.Has("updated_to") {
		if t, err := time.Parse("2006-01-02", q.Get("updated_to")); err == nil {
			filter.UpdatedTo.SetValue(t)
		}
	}
	if q.Has("sort") {
		sortBy := q.Get("sort")
		filter.SortBy.SetValue(sortBy)
	} else {
		sortBy := "client_id"
		filter.SortBy.SetValue(sortBy)
	}
	if q.Has("order") {
		sortOrder := q.Get("order")
		filter.SortOrder.SetValue(sortOrder)
	} else {
		sortOrder := "desc"
		filter.SortOrder.SetValue(sortOrder)
	}

	if filter.SortBy.Value == "" {
		filter.SortBy.SetValue("client_id")
	}
	if filter.SortOrder.Value == "" {
		filter.SortOrder.SetValue("desc")
	}
	return filter

}

// // GetClientHandler handler parse id from path, retrieve client from database
// // and return ClientDetails form to client.
func (h Handler) GetClientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		responseError(w, r, http.StatusBadRequest, "incorrect client id")
		return
	}
	client, err := h.DB.GetClient(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve client", "error", err, "id", id)
		responseError(w, r, http.StatusNotFound, "client not found")
		return
	}
	ui.ClientsViewSheetContent(client, ui.Form{}).Render(r.Context(), w)
}

// CreateClientHandler handler parse models.Client values from http form,
// trying to create client and return result. CreateSuccess if there is no errors
// and ClientsCreateForm with errors if user entered incorrect data.
func (h Handler) CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responseError(w, r, http.StatusBadRequest, "invalid form data")
		return
	}

	// TODO: add token generation and insertion
	form := ui.Form{}
	name := r.PostForm.Get("name")
	form["name"] = ui.FormField{
		Value: name,
		Err:   checkClientName(name),
	}
	email := r.PostForm.Get("email")
	form["email"] = ui.FormField{
		Value: email,
		Err:   checkEmail(email),
	}

	phone := r.PostForm.Get("phone")
	form["phone"] = ui.FormField{
		Value: phone,
		Err:   checkPhone(phone),
	}
	if len(form) != 0 {
		ui.ClientsAddContent(form).Render(r.Context(), w)
		return
	}

	createdClient, err := h.DB.CreateClient(r.Context(), name, email, phone)
	if err != nil {
		slog.Error("can't create client", "error", err)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	ui.CreateSuccess("Client added", "clients", int(createdClient.ClientID)).Render(r.Context(), w)
}

func checkClientName(name string) error {
	if utf8.RuneCountInString(name) <= 3 {
		return errors.New("client name must be at least 3 characters")
	}
	return nil
}

// checkEmail validates email format using net/mail.ParseAddress
func checkEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("incorrect email format")
	}
	return nil
}

// checkPhone validates phone format using a regex
func checkPhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return errors.New("incorrect phone format")
	}
	return nil
}

// DeleteClient handler process soft delete of client with id parsed from path.
func (h Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		responseError(w, r, http.StatusBadRequest, "incorrect client id")
		return
	}

	if err := h.DB.SoftDeleteClient(r.Context(), id); err != nil {
		slog.Error("can't delete client", "error", err, "id", id)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.GetClientHandler(w, r)
}

// BulkDeleteClients handler soft delete multiple clients with ids parsed from
// http form. Return ClientsTable without deleted records.
func (h Handler) BulkDeleteClients(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		responseError(w, r, http.StatusBadRequest, "invalid form data")
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		responseError(w, r, http.StatusBadRequest, "no clients selected")
		return
	}

	var ids []int
	for _, idStr := range selectedIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("can't parse client id", "error", err, "id", idStr)
			continue
		}
		ids = append(ids, id)
	}

	if err := h.DB.BulkSoftDeleteClients(r.Context(), ids); err != nil {
		slog.Error("can't delete clients batch", "error", err, "batch", ids)
	}

	page := parsePagination(r)

	filter := parseClientFilters(r)

	clients, total, err := h.DB.GetClients(r.Context(), page, filter)
	if err != nil {
		slog.Error("can't retrieve list of clients after delete", "error", err)
		responseError(w, r, http.StatusInternalServerError, "something went wrong")
		return
	}

	ui.ClientsTable(clients, page, total, filter).Render(r.Context(), w)
}

// UpdateClient handler parse id from path and client's new values from form.
// Return ClientDetail form with new values.
// func (h Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
// 	id, err := parseIDFromReq(r)
// 	if err != nil {
// 		slog.Error("can't parse id from URL path", "error", err)
// 		responseError(w, r, http.StatusBadRequest, "incorrect client id")
// 		return
// 	}
//
// 	// TODO: process a not existing client error
// 	// existingClient, err := h.DB.GetClient(r.Context(), id)
// 	// if err != nil {
// 	// 	slog.Error("can't recieve client", "error", err)
// 	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
// 	// 	return
// 	// }
//
// 	if err := r.ParseForm(); err != nil {
// 		slog.Error("can't http form", "error", err)
// 		responseError(w, r, http.StatusBadRequest, "invalid form data")
// 		return
// 	}
// 	name, email, phone, errs := parseClientUpdateArgs(r)
// 	if len(errs) != 0 {
// 		// ui.ClientEditForm(existingClient, errs).Render(r.Context(), w)
// 	}
//
// 	if err := h.DB.UpdateClient(r.Context(), id, name, email, phone); err != nil {
// 		slog.Error("can't update client", "error", err, "id", id)
// 		responseError(w, r, http.StatusInternalServerError, "something went wrong")
// 		return
// 	}
//
// 	updatedClient, err := h.DB.GetClient(r.Context(), id)
// 	if err != nil {
// 		slog.Error("can't fetch updated client", "error", err, "id", id)
// 		responseError(w, r, http.StatusInternalServerError, "something went wrong")
// 		return
// 	}
//
// 	ui.ClientDetail(updatedClient).Render(r.Context(), w)
// }

// parseClientUpdateArgs function parse r.Form and return Optional name, email and phone
// args for client update.
// The return values are predefined to shorten return values types and don't
// define them in fucntion body.
func parseClientUpdateArgs(r *http.Request) (name, email, phone models.Optional[string], errs map[string]string) {
	if n := r.PostForm.Get("name"); r.PostForm.Has("name") {
		if utf8.RuneCountInString(n) <= 3 {
			errs["name"] = "client name must be at least 3 characters"
		} else {
			name.SetValue(n)
		}
	}

	if e := r.PostForm.Get("email"); r.PostForm.Has("email") {
		if _, err := mail.ParseAddress(e); err != nil {
			errs["email"] = "incorrect email format"
		} else {
			email.SetValue(e)
		}
	}

	if p := r.PostForm.Get("phone"); r.PostForm.Has("phone") {
		if !phoneRegex.MatchString(p) {
			errs["email"] = "incorrect email format"
		} else {
			phone.SetValue(p)
		}
	}

	// return defined in the function header variables
	return
}

// VerifyEmail handler retrieve token from path and check if it exists in the
// database, if it is then email is verified.
func (h Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		responseError(w, r, http.StatusBadRequest, "token is required")
		return
	}

	if err := h.DB.VerifyClientEmail(r.Context(), token); err != nil {
		slog.Error("can't verify email", "error", err, "token", token)
		responseError(w, r, http.StatusBadRequest, "invalid or expired token")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email verified successfully"))
}

// GetClientOrders handler parse id from path; order filter, limit and offset
// from url values and retrieve all client orders with filters and pagination.
// Return OrdersTable.
func (h Handler) GetClientOrders(w http.ResponseWriter, r *http.Request) {
	// id, err := parseIDFromReq(r)
	// if err != nil {
	// 	slog.Error("can't parse id from URL path", "error", err)
	// 	responseError(w, r, http.StatusBadRequest, "incorrect client id")
	// 	return
	// }

	// limit, offset := parsePagination(r)
	// orders, total, err := h.DB.GetClientOrders(r.Context(), id, limit, offset)
	// if err != nil {
	// 	slog.Error("can't retrieve client orders", "error", err, "client_id", id)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }
	//
	// filter := parseOrderFilter(r)

	// ui.OrdersTable(orders, limit, offset, int(total), filter, []models.Transport{}).Render(r.Context(), w)
}

// NewClientPageHandler return create form.
func (h Handler) NewClientPageHandler(w http.ResponseWriter, r *http.Request) {
	// ui.ClientCreateForm(nil).Render(r.Context(), w)
}

// EditClientPageHandler parse id from path, retrive client by id and return
// edit form.
func (h Handler) EditClientPageHandler(w http.ResponseWriter, r *http.Request) {
	// id, err := parseIDFromReq(r)
	// if err != nil {
	// 	slog.Error("can't parse id from URL path", "error", err)
	// 	responseError(w, r, http.StatusBadRequest, "incorrect client id")
	// 	return
	// }
	//
	// client, err := h.DB.GetClient(r.Context(), id)
	// if err != nil {
	// 	slog.Error("can't retrieve client for edit", "error", err, "id", id)
	// 	responseError(w, r, http.StatusNotFound, "client not found")
	// 	return
	// }

	// ui.ClientEditForm(client, nil).Render(r.Context(), w)
}

// ExportClientsHandler generate export file and send it to user.
// TODO: implement
func (h Handler) ExportClientsHandler(w http.ResponseWriter, r *http.Request) {
	// // Parse filter parameters
	// filter := ui.ClientFilter{
	// 	Name:          r.URL.Query().Get("name"),
	// 	Email:         r.URL.Query().Get("email"),
	// 	Phone:         r.URL.Query().Get("phone"),
	// 	EmailVerified: r.URL.Query().Get("email_verified"),
	// 	CreatedFrom:   r.URL.Query().Get("created_from"),
	// 	CreatedTo:     r.URL.Query().Get("created_to"),
	// 	UpdatedFrom:   r.URL.Query().Get("updated_from"),
	// 	UpdatedTo:     r.URL.Query().Get("updated_to"),
	// 	SortBy:        r.URL.Query().Get("sort"),
	// 	SortOrder:     r.URL.Query().Get("order"),
	// }
	//
	// // Get all clients with filters (no pagination for export)
	// clients, err := h.DB.GetClients(r.Context(), filter)
	// if err != nil {
	// 	slog.Error("can't retrieve clients for export", "error", err)
	// 	responseError(w, r, http.StatusInternalServerError, "something went wrong")
	// 	return
	// }
	//
	// // Set CSV headers
	// w.Header().Set("Content-Type", "text/csv")
	// w.Header().Set("Content-Disposition", "attachment;filename=clients.csv")
	//
	// // Write CSV header
	// fmt.Fprintln(w, "ID,Name,Email,Email Verified,Phone,Created At,Updated At")
	//
	// // Write data
	// for _, client := range clients {
	// 	emailVerified := "false"
	// 	if client.EmailVerified {
	// 		emailVerified = "true"
	// 	}
	// 	fmt.Fprintf(w, "%d,%s,%s,%s,%s,%s,%s\n",
	// 		client.ClientID,
	// 		escapeCSV(client.Name),
	// 		escapeCSV(client.Email),
	// 		emailVerified,
	// 		escapeCSV(client.Phone),
	// 		client.CreatedAt.Format(time.RFC3339),
	// 		client.UpdatedAt.Format(time.RFC3339),
	// 	)
	// }
}
