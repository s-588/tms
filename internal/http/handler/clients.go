package handler

import (
	"errors"
	"fmt"
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
	page := parsePagination(r)

	filter := parseClientFilters(r)

	clients, total, err := h.DB.GetClients(r.Context(), 1, models.ClientFilter{})
	if err != nil {
		slog.Error("can't retrieve list of clients", "error", err)
		ui.Toast("error", "Can't render clients page", "Something went wrong").Render(r.Context(), w)
		return
	}
	slog.Debug("clients page", "clients", clients)
	ui.ClientsPage(clients, page, total, filter).Render(r.Context(), w)
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
		ui.Toast("error", "Can't get clients data", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("retrieve clients from database", "filter", filter, "page", page,
		"total pages", total, "total clients", len(clients))
	ui.ClientsTable(clients, page, total, filter, true).Render(r.Context(), w)
}

// parseClientFilters function parse models.ClientFilter from request url values
// and return it.
func parseClientFilters(r *http.Request) models.ClientFilter {
	filter := models.ClientFilter{}
	q := r.URL.Query()

	if err := checkClientName(q.Get("name")); q.Has("name") && err == nil{
		filter.Name.SetValue(q.Get("name"))
	}
	if err := checkEmail("email"); q.Has("email") && err == nil{
		filter.Email.SetValue(q.Get("email"))
	}
	if err := checkPhone("phone"); q.Has("phone") && err == nil{
		filter.Phone.SetValue(q.Get("phone"))
	}
	if q.Has("email_verified") {
		emailVerifiedStr := q.Get("email_verified")
		switch emailVerifiedStr {
case "true":
			filter.EmailVerified.SetValue(true)
		case "false":
			filter.EmailVerified.SetValue(false)
		}
	}
	if q.Has("created_from") && q.Has("created_from") {
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
		sortOrder := "desc"
		filter.SortOrder.SetValue(sortOrder)
	} else {
		sortBy := "created_at"
		filter.SortBy.SetValue(sortBy)
		sortOrder := "asc"
		filter.SortOrder.SetValue(sortOrder)
	}

	return filter

}

// GetClientHandler handler parse id from path, retrieve client from database
// and return ClientDetails form to client.
func (h Handler) GetClientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Can't get client data", "Something went wrong").Render(r.Context(), w)
		return
	}
	client, err := h.DB.GetClient(r.Context(), id)
	if err != nil {
		slog.Error("can't retrieve client", "error", err, "id", id)
		ui.Toast("error", "Can't get client data", "Not found").Render(r.Context(), w)
		return
	}
	slog.Debug("retrieve client", "client", client)
	ui.ClientsViewSheetContent(client, ui.Form{}).Render(r.Context(), w)
}

// CreateClientHandler handler parse models.Client values from http form,
// trying to create client and return result. CreateSuccess if there is no errors
// and ClientsCreateForm with errors if user entered incorrect data.
func (h Handler) CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't get parse form", "Invalid form data").Render(r.Context(), w)
		return
	}

	// TODO: add token generation and insertion
	hasError, form := parseClientForm(r)
	if hasError != nil {
		slog.Debug("incorrect input data for adding client", "data", form)
		ui.ClientsAddContent(form).Render(r.Context(), w)
		return
	}

	_, err := h.DB.CreateClient(r.Context(),
		form["name"].Value, form["email"].Value, form["phone"].Value)
	if err != nil {
		slog.Error("can't create client", "error", err)
		ui.Toast("error", "Can't create user", "Something went wrong").Render(r.Context(), w)
		ui.ClientsAddContent(form).Render(r.Context(), w)
		return
	}

	slog.Debug("adding new client", "data", form)
	ui.Toast("success", "User created", fmt.Sprintf("User %s successfully created", form["name"].Value)).Render(r.Context(), w)
	// h.GetClients(w, r)
}

func parseClientForm(r *http.Request) (err error, form ui.Form) {
	form = make(ui.Form)
	name := r.PostForm.Get("name")
	form["name"] = ui.FormField{
		Value: name,
	}
	if err = checkClientName(name); err != nil {
		form["name"] = ui.FormField{
			Value: name,
			Err:   err,
		}
	}

	email := r.PostForm.Get("email")
	form["email"] = ui.FormField{
		Value: email,
	}
	if err = checkEmail(email); err != nil {
		form["email"] = ui.FormField{
			Value: email,
			Err:   err,
		}
	}

	phone := r.PostForm.Get("phone")
	form["phone"] = ui.FormField{
		Value: phone,
	}
	if err = checkPhone(phone); err != nil {
		form["phone"] = ui.FormField{
			Value: phone,
			Err:   err,
		}
	}

	return
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
		slog.Debug("incorrect phone format", "phone", phone)
		return errors.New("incorrect phone format")
	}
	return nil
}

// DeleteClient handler process soft delete of client with id parsed from path.
func (h Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Incorrect URL", "Can't parse id from URL path").Render(r.Context(), w)
		return
	}

	if err := h.DB.SoftDeleteClient(r.Context(), id); err != nil {
		slog.Error("can't delete client", "error", err, "id", id)
		ui.Toast("error", "Can't delete client", "Something went wrong").Render(r.Context(), w)
		return
	}

	slog.Debug("deleting client", "clientID", id)
	ui.Toast("success", "Deleted", "Client successfully deleted")
	h.GetClients(w, r)
}

// BulkDeleteClients handler soft delete multiple clients with ids parsed from
// http form. Return ClientsTable without deleted records.
func (h Handler) BulkDeleteClients(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ui.Toast("error", "Can't parse form", "Invalid form data").Render(r.Context(), w)
		return
	}

	selectedIDs := r.Form["selected_ids"]
	if len(selectedIDs) == 0 {
		ui.Toast("error", "Can't delete clients", "No clients selected").Render(r.Context(), w)
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

	h.GetClients(w, r)
}

// UpdateClient handler parse id from path and client's new values from form.
// Return ClientDetail form with new values.
func (h Handler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromReq(r)
	if err != nil {
		slog.Error("can't parse id from URL path", "error", err)
		ui.Toast("error", "Error", "Incorrect client ID")
		return
	}

	// TODO: process a not existing client error
	existingClient, err := h.DB.GetClient(r.Context(), id)
	if err != nil {
		slog.Error("can't recieve client", "error", err)
		ui.Toast("error", "Internal error", "Something wen wrong")
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("can't http form", "error", err)
		ui.Toast("error", "Bad request", "Invalid form format")
		return
	}

	err, form := parseClientForm(r)
	if err != nil {
		slog.Debug("can't update clients", "form",form, "err",err)
		ui.ClientsViewSheetContent(existingClient, form).Render(r.Context(), w)
		return
	}

	if err := h.DB.UpdateClient(r.Context(), id, form["name"].Value, form["email"].Value, form["phone"].Value); err != nil {
		slog.Error("can't update client", "error", err, "id", id)
		ui.Toast("error", "Internal error", "something went wrong")
		return
	}

	slog.Debug("update client", "form data", form)
	ui.Toast("success", "Client updated", "Client successfully updated")
	h.GetClientHandler(w, r)
	h.GetClients(w,r)
}

// VerifyEmail handler retrieve token from path and check if it exists in the
// database, if it is then email is verified.
func (h Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		ui.Toast("error", "Can't verify email", "Token is required").Render(r.Context(), w)
		return
	}

	if err := h.DB.VerifyClientEmail(r.Context(), token); err != nil {
		slog.Error("can't verify email", "error", err, "token", token)
		ui.Toast("error", "Can't verify email", "Expire or invalid email").Render(r.Context(), w)
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
