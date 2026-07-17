package money

import "testing"

func TestFormatINR(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "₹0"},
		{499, "₹499"},
		{2999, "₹2,999"},
		{12500, "₹12,500"},
		{99999, "₹99,999"},
		{1_00_000, "₹1 L"},
		{3_70_000, "₹3.7 L"},
		{36_50_000, "₹36.5 L"},
		{99_99_999, "₹100 L"}, // rounds within lakh band
		{1_00_00_000, "₹1 Cr"},
		{1_20_00_000, "₹1.2 Cr"},
		{25_50_00_000, "₹25.5 Cr"},
		{-2999, "-₹2,999"},
	}
	for _, c := range cases {
		if got := FormatINR(c.in); got != c.want {
			t.Errorf("FormatINR(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGroupIndian(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{5, "5"},
		{999, "999"},
		{1000, "1,000"},
		{12500, "12,500"},
		{123456, "1,23,456"},
		{1234567, "12,34,567"},
		{123456789, "12,34,56,789"},
	}
	for _, c := range cases {
		if got := GroupIndian(c.in); got != c.want {
			t.Errorf("GroupIndian(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}
