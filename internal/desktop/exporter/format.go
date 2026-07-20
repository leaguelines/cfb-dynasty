package exporter

import (
	"strconv"
	"strings"
)

// fmtInt formats a plain int.
func fmtInt(v int) string { return strconv.Itoa(v) }

// fmtIntPtr formats a *int, returning "" for nil.
func fmtIntPtr(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}

// fmtUint32 formats a uint32, returning "" for zero.
func fmtUint32(v uint32) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(v), 10)
}

// FormatIntPtr formats a *int for display, returning "" for nil.
func FormatIntPtr(v *int) string { return fmtIntPtr(v) }

// fmtJoin joins strings for CSV/display.
func fmtJoin(vals []string) string { return JoinStrings(vals) }

// JoinStrings joins strings for display.
func JoinStrings(vals []string) string {
	if len(vals) == 0 {
		return ""
	}
	return strings.Join(vals, ", ")
}

// fmtBool formats a bool as "true"/"false".
func fmtBool(v bool) string { return strconv.FormatBool(v) }

// blanks returns a slice of n empty strings, used to pad nil nested records.
func blanks(n int) []string {
	out := make([]string, n)
	return out
}
