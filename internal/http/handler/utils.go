package handler

import (
	"log/slog"
	"net/http"

	"github.com/s-588/tms/internal/ui"
)

func renderToast(w http.ResponseWriter, r *http.Request, t, title, message string) {
	err := ui.Toast(t, title, message).Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

// StringOrNil returns nil if s is empty, otherwise a pointer to s.
func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
