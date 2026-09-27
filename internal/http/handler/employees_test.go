package handler

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/ui"
)

func TestHandler_GetEmployeesPage(t *testing.T) {
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
			h.GetEmployeesPage(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_GetEmployees(t *testing.T) {
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
			h.GetEmployees(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseEmployeeFilters(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want models.EmployeeFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseEmployeeFilters(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEmployeeFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_CreateEmployeeHandler(t *testing.T) {
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
			h.CreateEmployeeHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseEmployeeCreateForm(t *testing.T) {
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
			gotForm, err := parseEmployeeCreateForm(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseEmployeeCreateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseEmployeeCreateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_GetEmployeeHandler(t *testing.T) {
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
			h.GetEmployeeHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_UpdateEmployeeHandler(t *testing.T) {
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
			h.UpdateEmployeeHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_parseEmployeeUpdateForm(t *testing.T) {
	type args struct {
		r        *http.Request
		existing models.Employee
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
			gotForm, err := parseEmployeeUpdateForm(tt.args.r, tt.args.existing)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseEmployeeUpdateForm() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(gotForm, tt.wantForm) {
				t.Errorf("parseEmployeeUpdateForm() = %v, want %v", gotForm, tt.wantForm)
			}
		})
	}
}

func TestHandler_DeleteEmployeeHandler(t *testing.T) {
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
			h.DeleteEmployeeHandler(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_BulkDeleteEmployeesHandler(t *testing.T) {
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
			h.BulkDeleteEmployeesHandler(tt.args.w, tt.args.r)
		})
	}
}

func Test_checkEmployeeName(t *testing.T) {
	type args struct {
		name string
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
			if err := checkEmployeeName(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("checkEmployeeName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkEmployeeNameFilter(t *testing.T) {
	type args struct {
		name string
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
			if err := checkEmployeeNameFilter(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("checkEmployeeNameFilter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkEmployeeStatus(t *testing.T) {
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
			if err := checkEmployeeStatus(tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("checkEmployeeStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkEmployeeJobTitle(t *testing.T) {
	type args struct {
		job string
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
			if err := checkEmployeeJobTitle(tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("checkEmployeeJobTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
