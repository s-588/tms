package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
)

type Handler struct {
	DB db.DB
}

func NewHandler(db db.DB) Handler {
	return Handler{
		DB: db,
	}
}

func parseIDFromReq(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return 0, fmt.Errorf("can't parse ID from path: %w", err)
		}
		return id, nil
	}
	return 0, fmt.Errorf("no id value in path")
}

func parsePagination(r *http.Request) int32 {
	pageStr := r.URL.Query().Get("page")

	if pageStr == "" {
		return 1
	} else {
		p, _ := strconv.ParseInt(pageStr, 10, 32)
		if p <= 0 {
			// TODO: exclude this hardcoded limit to config value
			p = 1
		}
		return int32(p)
	}
}

func responseError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.WriteHeader(code)
	ui.ErrorMessage(msg).Render(r.Context(), w)
}
