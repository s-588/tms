package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/s-588/tms/internal/db"
	
)

type Handler struct {
	DB db.DB
}

func NewHandler(db db.DB) Handler {
	return Handler{
		DB: db,
	}
}

func parseIDFromReq(r *http.Request) (int32, error) {
	idStr := r.PathValue("id")
	if idStr != "" {
		id, err := strconv.ParseInt(idStr,10,32)
		if err != nil {
			return 0, fmt.Errorf("can't parse ID from path: %w", err)
		}
		return int32(id), nil
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