package handler

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
)

func TestHandler_GetInspectionsPage(t *testing.T) {
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
			h.GetInspectionsPage(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetInspections(t *testing.T) {
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
			h.GetInspections(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseInspectionFilters(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want models.InspectionFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInspectionFilters(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseInspectionFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_CreateInspectionHandler(t *testing.T) {
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
			h.CreateInspectionHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseInspectionCreateForm(t *testing.T) {
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
			gotForm, err := parseInspectionCreateForm(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseInspectionCreateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseInspectionCreateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_GetInspectionHandler(t *testing.T) {
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
			h.GetInspectionHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_UpdateInspectionHandler(t *testing.T) {
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
			h.UpdateInspectionHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseInspectionUpdateForm(t *testing.T) {
	type args struct {
		r        *http.Request
		existing models.Inspection
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
			gotForm, err := parseInspectionUpdateForm(tt.args.r, tt.args.existing)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseInspectionUpdateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseInspectionUpdateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_DeleteInspectionHandler(t *testing.T) {
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
			h.DeleteInspectionHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_BulkDeleteInspectionsHandler(t *testing.T) {
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
			h.BulkDeleteInspectionsHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_checkInspectionStatus(t *testing.T) {
	type args struct {
		status string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkInspectionStatus(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("checkInspectionStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
