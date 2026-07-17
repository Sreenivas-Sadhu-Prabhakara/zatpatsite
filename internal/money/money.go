// Package money formats amounts the way Indians read them:
// ₹1.2 Cr, ₹36.5 L, ₹12,500.
package money

import (
	"fmt"
	"strings"
)

// FormatINR renders whole rupees using Indian conventions.
//
//	>= 1 crore  -> ₹1.2 Cr
//	>= 1 lakh   -> ₹36.5 L
//	otherwise   -> ₹12,500 (Indian digit grouping)
func FormatINR(rupees int64) string {
	neg := ""
	if rupees < 0 {
		neg = "-"
		rupees = -rupees
	}
	switch {
	case rupees >= 1_00_00_000:
		return neg + "₹" + trimUnit(float64(rupees)/1_00_00_000) + " Cr"
	case rupees >= 1_00_000:
		return neg + "₹" + trimUnit(float64(rupees)/1_00_000) + " L"
	default:
		return neg + "₹" + GroupIndian(rupees)
	}
}

// trimUnit renders one decimal place, dropping a trailing ".0".
func trimUnit(v float64) string {
	s := fmt.Sprintf("%.1f", v)
	return strings.TrimSuffix(s, ".0")
}

// GroupIndian inserts Indian-style digit separators: last three digits,
// then groups of two (e.g. 1234567 -> 12,34,567).
func GroupIndian(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	head, tail := s[:len(s)-3], s[len(s)-3:]
	var parts []string
	for len(head) > 2 {
		parts = append([]string{head[len(head)-2:]}, parts...)
		head = head[:len(head)-2]
	}
	parts = append([]string{head}, parts...)
	return strings.Join(parts, ",") + "," + tail
}
