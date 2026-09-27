package handler

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
)

func TestHandler_GetOrdersPage(t *testing.T) {
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
			h.GetOrdersPage(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetOrders(t *testing.T) {
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
			h.GetOrders(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseOrderFilters(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want models.OrderFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOrderFilters(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOrderFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_CreateOrderHandler(t *testing.T) {
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
			h.CreateOrderHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseOrderCreateForm(t *testing.T) {
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
			gotForm, err := parseOrderCreateForm(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseOrderCreateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseOrderCreateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_GetOrderHandler(t *testing.T) {
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
			h.GetOrderHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_UpdateOrderHandler(t *testing.T) {
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
			h.UpdateOrderHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseOrderUpdateForm(t *testing.T) {
	type args struct {
		r        *http.Request
		existing models.Order
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
			gotForm, err := parseOrderUpdateForm(tt.args.r, tt.args.existing)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseOrderUpdateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseOrderUpdateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_DeleteOrderHandler(t *testing.T) {
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
			h.DeleteOrderHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_BulkDeleteOrdersHandler(t *testing.T) {
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
			h.BulkDeleteOrdersHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetOrderTransportsHandler(t *testing.T) {
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
			h.GetOrderTransportsHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_AssignOrderTransportsHandler(t *testing.T) {
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
			h.AssignOrderTransportsHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_checkOrderStatus(t *testing.T) {
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
			if err := checkOrderStatus(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("checkOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_int32Ptr(t *testing.T) {
	type args struct {
		i int32
	}
	tests := []struct {
		name string
		args args
		want *int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := int32Ptr(tt.args.i); got != tt.want {
				t.Errorf("int32Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GetOrderAddForm(t *testing.T) {
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
			h.GetOrderAddForm(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetOrderFilterForm(t *testing.T) {
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
			h.GetOrderFilterForm(tt.args.w, tt.args.r)
		})
	}
}

func Test_addListsToContext(t *testing.T) {
	type args struct {
		ctx context.Context
		db  db.DB
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addListsToContext(tt.args.ctx, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addListsToContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_OrdersExport(t *testing.T) {
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
			h := &Handler{
				DB: tt.fields.DB,
			}
			h.OrdersExport(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_fetchAllOrders(t *testing.T) {
	type fields struct {
		DB db.DB
	}
	type args struct {
		ctx    context.Context
		filter models.OrderFilter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				DB: tt.fields.DB,
			}
			got, err := h.fetchAllOrders(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Handler.fetchAllOrders() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.fetchAllOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GetOrdersForReport(t *testing.T) {
	type fields struct {
		DB db.DB
	}
	type args struct {
		ctx    context.Context
		filter models.OrderFilter
		period ReportPeriod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				DB: tt.fields.DB,
			}
			got, err := h.GetOrdersForReport(tt.args.ctx, tt.args.filter, tt.args.period)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Handler.GetOrdersForReport() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.GetOrdersForReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregateStats(t *testing.T) {
	type args struct {
		orders []models.Order
		period ReportPeriod
	}
	tests := []struct {
		name string
		args args
		want OrderStats
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := aggregateStats(tt.args.orders, tt.args.period); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("aggregateStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
