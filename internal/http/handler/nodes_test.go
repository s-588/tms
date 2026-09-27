package handler

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
)

func TestHandler_GetNodesPage(t *testing.T) {
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
			h.GetNodesPage(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetNodes(t *testing.T) {
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
			h.GetNodes(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseNodeFilters(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want models.NodeFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNodeFilters(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseNodeFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_CreateNodeHandler(t *testing.T) {
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
			h.CreateNodeHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseNodeCreateForm(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name     string
		args     args
		wantForm ui.Form
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotForm, err := parseNodeCreateForm(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseNodeCreateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseNodeCreateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_GetNodeHandler(t *testing.T) {
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
			h.GetNodeHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_UpdateNodeHandler(t *testing.T) {
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
			h.UpdateNodeHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseNodeUpdateForm(t *testing.T) {
	type args struct {
		r        *http.Request
		existing models.Node
	}
	tests := []struct {
		name     string
		args     args
		wantForm ui.Form
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotForm, err := parseNodeUpdateForm(tt.args.r, tt.args.existing)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseNodeUpdateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseNodeUpdateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_DeleteNodeHandler(t *testing.T) {
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
			h.DeleteNodeHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_BulkDeleteNodesHandler(t *testing.T) {
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
			h.BulkDeleteNodesHandler(tt.args.w, tt.args.r)
		})
	}
}
