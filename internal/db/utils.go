package db

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db/generated"
	"github.com/shopspring/decimal"
)

// ToInt32Ptr converts models.Optional[int32] to *int32.
func ToInt32Ptr(o models.Optional[int32]) *int32 {
	if o.Set {
		return &o.Value
	}
	return nil
}

// ToInt16Ptr converts models.Optional[int16] to *int16.
func ToInt16Ptr(o models.Optional[int16]) *int16 {
	if o.Set {
		return &o.Value
	}
	return nil
}

// ToStringPtr converts models.Optional[string] to *string.
func ToStringPtr(o models.Optional[string]) *string {
	if o.Set {
		return &o.Value
	}
	return nil
}

// ToBoolPtr converts models.Optional[bool] to *bool.
func ToBoolPtr(o models.Optional[bool]) *bool {
	if o.Set {
		return &o.Value
	}
	return nil
}

// ToInt16PtrFromUint8 converts models.Optional[uint8] to *int16.
func ToInt16PtrFromUint8(o models.Optional[uint8]) *int16 {
	if o.Set {
		v := int16(o.Value)
		return &v
	}
	return nil
}

// ToInt16PtrFromInt converts models.Optional[int] to *int16.
func ToInt16PtrFromInt(o models.Optional[int]) *int16 {
	if o.Set {
		v := int16(o.Value)
		return &v
	}
	return nil
}

// ToPgTypeNumeric converts models.Optional[decimal.Decimal] to pgtype.Numeric.
func ToPgTypeNumeric(o models.Optional[decimal.Decimal]) pgtype.Numeric {
	var n pgtype.Numeric
	if o.Set {
		n.Int = o.Value.Coefficient()
		n.Exp = o.Value.Exponent()
		n.Valid = true
	}
	return n
}

// ToPgTypeNumericFromDecimal converts a required decimal.Decimal to pgtype.Numeric.
func ToPgTypeNumericFromDecimal(d decimal.Decimal) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   d.Coefficient(),
		Exp:   d.Exponent(),
		Valid: true,
	}
}

// ToPgTypeDate converts models.Optional[time.Time] to pgtype.Date.
func ToPgTypeDate(o models.Optional[time.Time]) pgtype.Date {
	var d pgtype.Date
	if o.Set {
		d.Time = o.Value
		d.Valid = true
	}
	return d
}

// ToPgTypeDateFromTime converts a required time.Time to pgtype.Date.
func ToPgTypeDateFromTime(t time.Time) pgtype.Date {
	return pgtype.Date{Time: t, Valid: true}
}

// ToPgTypeTimestamptz converts models.Optional[time.Time] to pgtype.Timestamptz.
func ToPgTypeTimestamptz(o models.Optional[time.Time]) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	if o.Set {
		ts.Time = o.Value
		ts.Valid = true
	}
	return ts
}

// ToPgTypeTimestamptzFromTime converts a required time.Time to pgtype.Timestamptz.
func ToPgTypeTimestamptzFromTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// ToNullOrderStatus converts models.Optional[models.OrderStatus] to generated.NullOrderStatus.
func ToNullOrderStatus(o models.Optional[models.OrderStatus]) generated.NullOrderStatus {
	if o.Set {
		return generated.NullOrderStatus{
			OrderStatus: generated.OrderStatus(o.Value),
			Valid:       true,
		}
	}
	return generated.NullOrderStatus{Valid: false}
}

// ToNullEmployeeStatus converts models.Optional[models.EmployeeStatus] to generated.NullEmployeeStatus.
func ToNullEmployeeStatus(o models.Optional[models.EmployeeStatus]) generated.NullEmployeeStatus {
	if o.Set {
		return generated.NullEmployeeStatus{
			EmployeeStatus: generated.EmployeeStatus(o.Value),
			Valid:          true,
		}
	}
	return generated.NullEmployeeStatus{Valid: false}
}

// ToNullInspectionStatus converts models.Optional[models.InspectionStatus] to generated.NullInspectionStatus.
func ToNullInspectionStatus(o models.Optional[models.InspectionStatus]) generated.NullInspectionStatus {
	if o.Set {
		return generated.NullInspectionStatus{
			InspectionStatus: generated.InspectionStatus(o.Value),
			Valid:            true,
		}
	}
	return generated.NullInspectionStatus{Valid: false}
}

// ToNullEmployeeJobTitle converts models.Optional[models.EmployeeJobTitle] to generated.NullEmployeeJobTitle.
func ToNullEmployeeJobTitle(o models.Optional[models.EmployeeJobTitle]) generated.NullEmployeeJobTitle {
	if o.Set {
		return generated.NullEmployeeJobTitle{
			EmployeeJobTitle: generated.EmployeeJobTitle(o.Value),
			Valid:            true,
		}
	}
	return generated.NullEmployeeJobTitle{Valid: false}
}

// fromPgTimestamptz converts pgtype.Timestamptz to time.Time (zero if invalid).
func fromPgTimestamptz(ts pgtype.Timestamptz) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}

// fromPgDate converts pgtype.Date to time.Time (zero if invalid).
func fromPgDate(d pgtype.Date) time.Time {
	if d.Valid {
		return d.Time
	}
	return time.Time{}
}

// fromPgNumeric converts pgtype.Numeric to decimal.Decimal (zero if invalid).
func fromPgNumeric(n pgtype.Numeric) decimal.Decimal {
	if n.Valid {
		return decimal.NewFromBigInt(n.Int, n.Exp)
	}
	return decimal.Zero
}

// fromStringPtr converts *string to string, returning empty if nil.
func fromStringPtr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// fromInt16PtrToUint8 converts *int16 to uint8, returning 0 if nil.
func fromInt16PtrToUint8(p *int16) uint8 {
	if p != nil {
		return uint8(*p)
	}
	return 0
}

// fromInt32PtrToInt converts *int32 to int, returning 0 if nil.
func fromInt32PtrToInt(p *int32) int {
	if p != nil {
		return int(*p)
	}
	return 0
}

// convertIntSliceToInt32 converts []int to []int32.
func convertIntSliceToInt32(ids []int) []int32 {
	result := make([]int32, len(ids))
	for i, v := range ids {
		result[i] = int32(v)
	}
	return result
}

// timeToPgDatePtr converts a *time.Time to pgtype.Date (NULL if nil)
func timeToPgDatePtr(t *time.Time) pgtype.Date {
	var d pgtype.Date
	if t != nil {
		d.Time = *t
		d.Valid = true
	}
	return d
}

// optionalTimeToPgDate converts models.Optional[time.Time] to pgtype.Date
func optionalTimeToPgDate(o models.Optional[time.Time]) pgtype.Date {
	var d pgtype.Date
	if o.Set {
		d.Time = o.Value
		d.Valid = true
	}
	return d
}

// optionalTimeToPgTimestamptz converts models.Optional[time.Time] to pgtype.Timestamptz
func optionalTimeToPgTimestamptz(o models.Optional[time.Time]) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	if o.Set {
		ts.Time = o.Value
		ts.Valid = true
	}
	return ts
}

// fromDecimalPtrToPgNumeric converts a *decimal.Decimal to pgtype.Numeric (NULL if nil)
func fromDecimalPtrToPgNumeric(d *decimal.Decimal) pgtype.Numeric {
	var n pgtype.Numeric
	if d != nil {
		_ = n.Scan(d.String())
	}
	return n
}

// optionalDecimalToPgNumeric converts models.Optional[decimal.Decimal] to pgtype.Numeric
func optionalDecimalToPgNumeric(o models.Optional[decimal.Decimal]) pgtype.Numeric {
	var n pgtype.Numeric
	if o.Set {
		_ = n.Scan(o.Value.String())
	}
	return n
}

// stringToPoint converts a string like "(x,y)" to pgtype.Point
func stringToPoint(s string) (pgtype.Point, error) {
	var p pgtype.Point
	err := p.Scan(s)
	return p, err
}

// pointToString converts pgtype.Point to its string representation
func pointToString(p pgtype.Point) string {
	if !p.Valid {
		return ""
	}
	return fmt.Sprintf("(%f,%f)", p.P.X, p.P.Y)
}
