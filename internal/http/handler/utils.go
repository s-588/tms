package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/s-588/tms/cmd/models"
)

// StringOrNil returns nil if s is empty, otherwise a pointer to s.
func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// escapeCSV escapes double quotes and wraps field in quotes if needed.
func escapeCSV(s string) string {
	if strings.ContainsAny(s, `",`) {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

// parseOptionalTime parses RFC3339 time string into models.Optional[time.Time].
func parseOptionalTime(timeStr string) models.Optional[time.Time] {
	var opt models.Optional[time.Time]
	if timeStr != "" {
		if t, err := time.Parse(time.RFC3339, timeStr); err == nil && !t.IsZero() {
			opt.SetValue(t)
		}
	}
	return opt
}

// parseOptionalString creates models.Optional[string] from a non-empty string.
func parseOptionalString(str string) models.Optional[string] {
	var opt models.Optional[string]
	if str != "" {
		opt.SetValue(str)
	}
	return opt
}

// parseOptionalInt32 converts a string to int32 and wraps in models.Optional[int32].
// If parsing fails or value <= 0, the Optional is left unset.
func parseOptionalInt32(str string) models.Optional[int32] {
	var opt models.Optional[int32]
	if str != "" {
		if val, err := strconv.Atoi(str); err == nil && val > 0 {
			opt.SetValue(int32(val))
		}
	}
	return opt
}
