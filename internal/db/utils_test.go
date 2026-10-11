package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

func TestToInt32Ptr(t *testing.T) {
	type args struct {
		o models.Optional[int32]
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
			if got := ToInt32Ptr(tt.args.o); got != tt.want {
				t.Errorf("ToInt32Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt16Ptr(t *testing.T) {
	type args struct {
		o models.Optional[int16]
	}
	tests := []struct {
		name string
		args args
		want *int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt16Ptr(tt.args.o); got != tt.want {
				t.Errorf("ToInt16Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStringPtr(t *testing.T) {
	type args struct {
		o models.Optional[string]
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
			if got := ToStringPtr(tt.args.o); got != tt.want {
				t.Errorf("ToStringPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBoolPtr(t *testing.T) {
	type args struct {
		o models.Optional[bool]
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToBoolPtr(tt.args.o); got != tt.want {
				t.Errorf("ToBoolPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt16PtrFromUint8(t *testing.T) {
	type args struct {
		o models.Optional[uint8]
	}
	tests := []struct {
		name string
		args args
		want *int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt16PtrFromUint8(tt.args.o); got != tt.want {
				t.Errorf("ToInt16PtrFromUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt16PtrFromInt(t *testing.T) {
	type args struct {
		o models.Optional[int]
	}
	tests := []struct {
		name string
		args args
		want *int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToInt16PtrFromInt(tt.args.o); got != tt.want {
				t.Errorf("ToInt16PtrFromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPgTypeNumeric(t *testing.T) {
	type args struct {
		o models.Optional[decimal.Decimal]
	}
	tests := []struct {
		name string
		args args
		want pgtype.Numeric
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPgTypeNumeric(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPgTypeNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPgTypeNumericFromDecimal(t *testing.T) {
	type args struct {
		d decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want pgtype.Numeric
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPgTypeNumericFromDecimal(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPgTypeNumericFromDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPgTypeDate(t *testing.T) {
	type args struct {
		o models.Optional[time.Time]
	}
	tests := []struct {
		name string
		args args
		want pgtype.Date
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPgTypeDate(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPgTypeDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPgTypeDateFromTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want pgtype.Date
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPgTypeDateFromTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPgTypeDateFromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPgTypeTimestamptz(t *testing.T) {
	type args struct {
		o models.Optional[time.Time]
	}
	tests := []struct {
		name string
		args args
		want pgtype.Timestamptz
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPgTypeTimestamptz(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPgTypeTimestamptz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPgTypeTimestamptzFromTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want pgtype.Timestamptz
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToPgTypeTimestamptzFromTime(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToPgTypeTimestamptzFromTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullOrderStatus(t *testing.T) {
	type args struct {
		o models.Optional[models.OrderStatus]
	}
	tests := []struct {
		name string
		args args
		want generated.NullOrderStatus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNullOrderStatus(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNullOrderStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullEmployeeStatus(t *testing.T) {
	type args struct {
		o models.Optional[models.EmployeeStatus]
	}
	tests := []struct {
		name string
		args args
		want generated.NullEmployeeStatus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNullEmployeeStatus(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNullEmployeeStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullInspectionStatus(t *testing.T) {
	type args struct {
		o models.Optional[models.InspectionStatus]
	}
	tests := []struct {
		name string
		args args
		want generated.NullInspectionStatus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNullInspectionStatus(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNullInspectionStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullEmployeeJobTitle(t *testing.T) {
	type args struct {
		o models.Optional[models.EmployeeJobTitle]
	}
	tests := []struct {
		name string
		args args
		want generated.NullEmployeeJobTitle
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToNullEmployeeJobTitle(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToNullEmployeeJobTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromPgTimestamptz(t *testing.T) {
	type args struct {
		ts pgtype.Timestamptz
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromPgTimestamptz(tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromPgTimestamptz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromPgDate(t *testing.T) {
	type args struct {
		d pgtype.Date
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromPgDate(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromPgDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromPgNumeric(t *testing.T) {
	type args struct {
		n pgtype.Numeric
	}
	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromPgNumeric(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromPgNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromStringPtr(t *testing.T) {
	type args struct {
		s *string
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
			if got := fromStringPtr(tt.args.s); got != tt.want {
				t.Errorf("fromStringPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromInt16PtrToUint8(t *testing.T) {
	type args struct {
		p *int16
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromInt16PtrToUint8(tt.args.p); got != tt.want {
				t.Errorf("fromInt16PtrToUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromInt32PtrToInt(t *testing.T) {
	type args struct {
		p *int32
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromInt32PtrToInt(tt.args.p); got != tt.want {
				t.Errorf("fromInt32PtrToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertIntSliceToInt32(t *testing.T) {
	type args struct {
		ids []int
	}
	tests := []struct {
		name string
		args args
		want []int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertIntSliceToInt32(tt.args.ids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertIntSliceToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_timeToPgDatePtr(t *testing.T) {
	type args struct {
		t *time.Time
	}
	tests := []struct {
		name string
		args args
		want pgtype.Date
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timeToPgDatePtr(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeToPgDatePtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optionalTimeToPgDate(t *testing.T) {
	type args struct {
		o models.Optional[time.Time]
	}
	tests := []struct {
		name string
		args args
		want pgtype.Date
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optionalTimeToPgDate(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("optionalTimeToPgDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optionalTimeToPgTimestamptz(t *testing.T) {
	type args struct {
		o models.Optional[time.Time]
	}
	tests := []struct {
		name string
		args args
		want pgtype.Timestamptz
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optionalTimeToPgTimestamptz(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("optionalTimeToPgTimestamptz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fromDecimalPtrToPgNumeric(t *testing.T) {
	type args struct {
		d *decimal.Decimal
	}
	tests := []struct {
		name string
		args args
		want pgtype.Numeric
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fromDecimalPtrToPgNumeric(tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromDecimalPtrToPgNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optionalDecimalToPgNumeric(t *testing.T) {
	type args struct {
		o models.Optional[decimal.Decimal]
	}
	tests := []struct {
		name string
		args args
		want pgtype.Numeric
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := optionalDecimalToPgNumeric(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("optionalDecimalToPgNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringToPoint(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    pgtype.Point
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToPoint(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Fatalf("stringToPoint() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pointToString(t *testing.T) {
	type args struct {
		p pgtype.Point
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
			if got := pointToString(tt.args.p); got != tt.want {
				t.Errorf("pointToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64Ptr(t *testing.T) {
	type args struct {
		o models.Optional[float64]
	}
	tests := []struct {
		name string
		args args
		want *float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToFloat64Ptr(tt.args.o); got != tt.want {
				t.Errorf("ToFloat64Ptr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toFloat64(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toFloat64(tt.args.v); got != tt.want {
				t.Errorf("toFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
