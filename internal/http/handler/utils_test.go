package handler

import (
	"reflect"
	"testing"
	"time"

	"github.com/s-588/tms/cmd/models"
)

func TestStringOrNil(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringOrNil(tt.args.s); got != tt.want {
				t.Errorf("StringOrNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_escapeCSV(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeCSV(tt.args.s); got != tt.want {
				t.Errorf("escapeCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOptionalTime(t *testing.T) {
	type args struct {
		timeStr string
	}
	tests := []struct {
		name string
		args args
		want models.Optional[time.Time]
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOptionalTime(tt.args.timeStr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOptionalTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOptionalString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want models.Optional[string]
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOptionalString(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOptionalString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseOptionalInt32(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want models.Optional[int32]
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseOptionalInt32(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOptionalInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}
