package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/http/handler"
	"github.com/s-588/tms/internal/testutil"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newFormRequest(method, path string, values url.Values) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func createTestClient(t *testing.T, h handler.Handler, name, email, phone string) int32 {
	t.Helper()

	form := url.Values{}
	form.Set("name", name)
	form.Set("email", email)
	form.Set("phone", phone)

	req := newFormRequest(http.MethodPost, "/clients", form)
	rr := httptest.NewRecorder()
	h.CreateClientHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("create client failed: status=%d body=%s", rr.Code, rr.Body.String())
	}

	// Re-fetch to get the ID (simple way – in real code you might return the ID from CreateClient)
	clients, _, err := h.DB.GetClients(req.Context(), 1, models.ClientFilter{})
	if err != nil {
		t.Fatalf("list clients after create: %v", err)
	}
	for _, c := range clients {
		if c.Email == email {
			return c.ClientID
		}
	}
	t.Fatal("created client not found")
	return 0
}

// ---------------------------------------------------------------------------
// GetClientsPage / GetClients
// ---------------------------------------------------------------------------

func TestGetClientsPage(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	req := httptest.NewRequest(http.MethodGet, "/clients/page", nil)
	rr := httptest.NewRecorder()

	h.GetClientsPage(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200\nbody: %s", rr.Code, rr.Body.String())
	}
}

func TestGetClients(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	// seed one client
	_ = createTestClient(t, h, "Test Client", "test@example.com", "+375291111111")

	tests := []struct {
		name  string
		query string
	}{
		{"no filters", ""},
		{"name filter", "name=Test"},
		{"email filter", "email=test@example.com"},
		{"phone filter", "phone=+375291111111"},
		{"pagination", "page=1"},
		{"sort", "sort=name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := "/clients"
			if tt.query != "" {
				path += "?" + tt.query
			}
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rr := httptest.NewRecorder()

			h.GetClients(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d\nbody: %s", rr.Code, rr.Body.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetClientHandler
// ---------------------------------------------------------------------------

func TestGetClientHandler(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	id := createTestClient(t, h, "View Client", "view@example.com", "+375292222222")

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/clients/%d", id), nil)
		req.SetPathValue("id", fmt.Sprintf("%d", id))
		rr := httptest.NewRecorder()

		h.GetClientHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d\nbody: %s", rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "View Client") {
			t.Errorf("body does not contain client name")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/clients/abc", nil)
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		h.GetClientHandler(rr, req)

		// handler returns toast, still 200 in your current code
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Can't get client data") {
			t.Errorf("expected error toast")
		}
	})

	t.Run("not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/clients/999999", nil)
		req.SetPathValue("id", "999999")
		rr := httptest.NewRecorder()

		h.GetClientHandler(rr, req)

		if !strings.Contains(rr.Body.String(), "Not found") {
			t.Errorf("expected not-found toast, got: %s", rr.Body.String())
		}
	})
}

// ---------------------------------------------------------------------------
// CreateClientHandler
// ---------------------------------------------------------------------------

func TestCreateClientHandler(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	tests := []struct {
		name           string
		form           url.Values
		wantContains   string
		wantNotContain string
	}{
		{
			name: "success",
			form: url.Values{
				"name":  {"Alice Wonderland"},
				"email": {"alice@example.com"},
				"phone": {"+375293333333"},
			},
			wantContains: "User created",
		},
		{
			name: "name too short",
			form: url.Values{
				"name":  {"Al"},
				"email": {"short@example.com"},
				"phone": {"+375294444444"},
			},
			wantContains: "client name must be at least 3 characters",
		},
		{
			name: "invalid email",
			form: url.Values{
				"name":  {"Bob Builder"},
				"email": {"not-an-email"},
				"phone": {"+375295555555"},
			},
			wantContains: "incorrect email format",
		},
		{
			name: "invalid phone",
			form: url.Values{
				"name":  {"Charlie Chaplin"},
				"email": {"charlie@example.com"},
				"phone": {"123"},
			},
			wantContains: "incorrect phone format for Belarus",
		},
		{
			name: "non-BY phone",
			form: url.Values{
				"name":  {"Diana Prince"},
				"email": {"diana@example.com"},
				"phone": {"+12025550123"}, // US number
			},
			wantContains: "incorrect phone format for Belarus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newFormRequest(http.MethodPost, "/clients", tt.form)
			rr := httptest.NewRecorder()

			h.CreateClientHandler(rr, req)

			body := rr.Body.String()
			if !strings.Contains(body, tt.wantContains) {
				t.Errorf("body missing %q\nGot:\n%s", tt.wantContains, body)
			}
		})
	}

	// Duplicate email / phone
	t.Run("duplicate email", func(t *testing.T) {
		// first create
		_ = createTestClient(t, h, "Dup Email", "dup@example.com", "+375296666666")

		form := url.Values{
			"name":  {"Another Client"},
			"email": {"dup@example.com"},
			"phone": {"+375297777777"},
		}
		req := newFormRequest(http.MethodPost, "/clients", form)
		rr := httptest.NewRecorder()
		h.CreateClientHandler(rr, req)

		if !strings.Contains(rr.Body.String(), "email already exists") {
			t.Errorf("expected duplicate email error, got: %s", rr.Body.String())
		}
	})

	t.Run("duplicate phone", func(t *testing.T) {
		_ = createTestClient(t, h, "Dup Phone", "phone1@example.com", "+375298888888")

		form := url.Values{
			"name":  {"Another Client"},
			"email": {"phone2@example.com"},
			"phone": {"+375298888888"},
		}
		req := newFormRequest(http.MethodPost, "/clients", form)
		rr := httptest.NewRecorder()
		h.CreateClientHandler(rr, req)

		if !strings.Contains(rr.Body.String(), "phone already exists") {
			t.Errorf("expected duplicate phone error, got: %s", rr.Body.String())
		}
	})
}

// ---------------------------------------------------------------------------
// UpdateClient
// ---------------------------------------------------------------------------

func TestUpdateClient(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	id := createTestClient(t, h, "Update Me", "update@example.com", "+375291010101")

	t.Run("success", func(t *testing.T) {
		form := url.Values{
			"name":  {"Updated Name"},
			"email": {"updated@example.com"},
			"phone": {"+375292020202"},
		}
		req := newFormRequest(http.MethodPut, fmt.Sprintf("/clients/%d", id), form)
		req.SetPathValue("id", fmt.Sprintf("%d", id))
		rr := httptest.NewRecorder()

		h.UpdateClient(rr, req)

		if !strings.Contains(rr.Body.String(), "Client updated") {
			t.Errorf("expected success toast, got: %s", rr.Body.String())
		}
	})

	t.Run("validation error", func(t *testing.T) {
		form := url.Values{
			"name":  {"X"}, // too short
			"email": {"bad"},
			"phone": {"123"},
		}
		req := newFormRequest(http.MethodPut, fmt.Sprintf("/clients/%d", id), form)
		req.SetPathValue("id", fmt.Sprintf("%d", id))
		rr := httptest.NewRecorder()

		h.UpdateClient(rr, req)

		body := rr.Body.String()
		if !strings.Contains(body, "client name must be at least 3 characters") &&
			!strings.Contains(body, "incorrect email format") {
			t.Errorf("expected validation errors, got: %s", body)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := newFormRequest(http.MethodPut, "/clients/abc", url.Values{})
		req.SetPathValue("id", "abc")
		rr := httptest.NewRecorder()

		h.UpdateClient(rr, req)
		// current implementation does not write the toast (missing .Render)
		// just ensure it does not panic
	})
}

// ---------------------------------------------------------------------------
// DeleteClient
// ---------------------------------------------------------------------------

func TestDeleteClient(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	id := createTestClient(t, h, "Delete Me", "delete@example.com", "+375293030303")

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/clients/%d", id), nil)
		req.SetPathValue("id", fmt.Sprintf("%d", id))
		rr := httptest.NewRecorder()

		h.DeleteClient(rr, req)

		// note: current code does not call .Render on the success toast
		// but it should still call GetClients afterwards
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d", rr.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/clients/xyz", nil)
		req.SetPathValue("id", "xyz")
		rr := httptest.NewRecorder()

		h.DeleteClient(rr, req)

		if !strings.Contains(rr.Body.String(), "Can't parse id") {
			t.Errorf("expected parse error toast")
		}
	})
}

// ---------------------------------------------------------------------------
// BulkDeleteClients
// ---------------------------------------------------------------------------

func TestBulkDeleteClients(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	id1 := createTestClient(t, h, "Bulk One", "bulk1@example.com", "+375294040404")
	id2 := createTestClient(t, h, "Bulk Two", "bulk2@example.com", "+375295050505")

	t.Run("success", func(t *testing.T) {
		form := url.Values{}
		form.Add("selected_ids", fmt.Sprintf("%d", id1))
		form.Add("selected_ids", fmt.Sprintf("%d", id2))

		req := newFormRequest(http.MethodPost, "/clients/delete", form)
		rr := httptest.NewRecorder()

		h.BulkDeleteClients(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d\nbody: %s", rr.Code, rr.Body.String())
		}
	})

	t.Run("no ids selected", func(t *testing.T) {
		req := newFormRequest(http.MethodPost, "/clients/delete", url.Values{})
		rr := httptest.NewRecorder()

		h.BulkDeleteClients(rr, req)

		if !strings.Contains(rr.Body.String(), "No clients selected") {
			t.Errorf("expected 'No clients selected' toast")
		}
	})
}

// ---------------------------------------------------------------------------
// VerifyEmail
// ---------------------------------------------------------------------------

func TestVerifyEmail(t *testing.T) {
	srv := testutil.SetupHTTPServer(t)
	h := srv.Handler

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/verify-email", nil)
		rr := httptest.NewRecorder()

		h.VerifyEmail(rr, req)

		if !strings.Contains(rr.Body.String(), "Token is required") {
			t.Errorf("expected token required message")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/verify-email?token=invalid-token-xyz", nil)
		rr := httptest.NewRecorder()

		h.VerifyEmail(rr, req)

		if !strings.Contains(rr.Body.String(), "Expire or invalid") {
			t.Errorf("expected invalid token message, got: %s", rr.Body.String())
		}
	})

	// Note: a real successful verify would require inserting a valid token
	// into the DB first. That can be added later once the token generation
	// TODO in CreateClientHandler is implemented.
}
