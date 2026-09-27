package handler

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/s-588/tms/internal/db"
)

func TestHandler_buildReportData(t *testing.T) {
	type fields struct {
		DB db.DB
	}
	type args struct {
		ctx     context.Context
		orderID int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ReportData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				DB: tt.fields.DB,
			}
			got, err := h.buildReportData(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Handler.buildReportData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Handler.buildReportData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateBytes(t *testing.T) {
	type args struct {
		templatePath string
		data         ReportData
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateBytes(tt.args.templatePath, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("generateBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_DownloadContract(t *testing.T) {
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
			h.DownloadContract(tt.args.w, tt.args.r)
		})
	}
}

func TestHandler_DownloadAct(t *testing.T) {
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
			h.DownloadAct(tt.args.w, tt.args.r)
		})
	}
}

func TestGenerateOrdersReport(t *testing.T) {
	type args struct {
		stats OrderStats
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateOrdersReport(tt.args.stats)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateOrdersReport() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateOrdersReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
