package handler

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
)

func TestHandler_GetPricesPage(t *testing.T) {
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
			h.GetPricesPage(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetPrices(t *testing.T) {
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
			h.GetPrices(tt.args.w, tt.args.r)
		})
	}
}

func Test_parsePriceFilters(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want models.PriceFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePriceFilters(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePriceFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_CreatePriceHandler(t *testing.T) {
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
			h.CreatePriceHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parsePriceCreateForm(t *testing.T) {
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
			gotForm, err := parsePriceCreateForm(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parsePriceCreateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parsePriceCreateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_GetPriceHandler(t *testing.T) {
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
			h.GetPriceHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_UpdatePriceHandler(t *testing.T) {
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
			h.UpdatePriceHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parsePriceUpdateForm(t *testing.T) {
	type args struct {
		r        *http.Request
		existing models.Price
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
			gotForm, err := parsePriceUpdateForm(tt.args.r, tt.args.existing)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parsePriceUpdateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parsePriceUpdateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_DeletePriceHandler(t *testing.T) {
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
			h.DeletePriceHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_BulkDeletePricesHandler(t *testing.T) {
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
			h.BulkDeletePricesHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_NewPricePageHandler(t *testing.T) {
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
			h.NewPricePageHandler(tt.args.w, tt.args.r)
		})
	}
}
