package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/s-588/tms/db"
	"github.com/s-588/tms/static/template"
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

func parsePagination(r *http.Request) (limit, offset int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr == "" {
		limit = 10
	} else {
		limit, _ = strconv.Atoi(limitStr)
		if limit <= 0 {
			// TODO: exclude this hardcoded limit to config value
			limit = 10
		}
	}

	offset, _ = strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

func responseError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.WriteHeader(code)
	template.ErrorMessage(msg).Render(r.Context(), w)
}
