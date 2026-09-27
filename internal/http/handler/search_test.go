package handler

import (
	"net/http"
	"testing"

	"github.com/s-588/tms/internal/db"
)

func TestHandler_SearchHandler(t *testing.T) {
	type fields struct {
		DB db.DB
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				DB: tt.fields.DB,
			}
			h.SearchHandler(tt.args.w, tt.args.r)
		})
	}
}
