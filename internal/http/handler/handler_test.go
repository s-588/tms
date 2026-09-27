package handler

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/s-588/tms/internal/db"
)

func TestNewHandler(t *testing.T) {
	type args struct {
		db db.DB
	}
	tests := []struct {
		name string
		args args
		want Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseIDFromReq(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIDFromReq(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseIDFromReq() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("parseIDFromReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePagination(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePagination(tt.args.r); got != tt.want {
				t.Errorf("parsePagination() = %v, want %v", got, tt.want)
			}
		})
	}
}
