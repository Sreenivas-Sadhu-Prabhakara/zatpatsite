package hours

import (
	"testing"
	"time"

	"zatpatsite/internal/model"
)

// at builds a time on a fixed known week in IST.
// 2026-07-12 is a Sunday.
func at(weekday int, hh, mm int) time.Time {
	return time.Date(2026, 7, 12+weekday, hh, mm, 0, 0, IST)
}

func week(open, close string) [7]model.DayHours {
	var h [7]model.DayHours
	for i := range h {
		h[i] = model.DayHours{Open: open, Close: close}
	}
	return h
}

func TestIsOpenBasic(t *testing.T) {
	h := week("09:30", "20:30")
	cases := []struct {
		name string
		t    time.Time
		want bool
	}{
		{"mid-morning open", at(1, 11, 0), true},
		{"exactly at open", at(2, 9, 30), true},
		{"one minute before open", at(2, 9, 29), false},
		{"exactly at close is closed", at(3, 20, 30), false},
		{"one minute before close", at(3, 20, 29), true},
		{"late night closed", at(4, 23, 0), false},
		{"early morning closed", at(5, 6, 0), false},
	}
	for _, c := range cases {
		if got := IsOpen(h, c.t); got != c.want {
			t.Errorf("%s: IsOpen = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestIsOpenClosedDay(t *testing.T) {
	h := week("09:30", "20:30")
	h[1].Closed = true // Monday off
	if IsOpen(h, at(1, 12, 0)) {
		t.Error("expected closed on weekly off day")
	}
	if !IsOpen(h, at(2, 12, 0)) {
		t.Error("expected open on Tuesday")
	}
}

func TestIsOpenOvernight(t *testing.T) {
	// Dhaba open 18:00 - 01:00 every day.
	h := week("18:00", "01:00")
	cases := []struct {
		name string
		t    time.Time
		want bool
	}{
		{"evening open", at(1, 22, 0), true},
		{"past midnight still open (yesterday's window)", at(2, 0, 30), true},
		{"after overnight close", at(2, 1, 30), false},
		{"afternoon closed", at(2, 15, 0), false},
	}
	for _, c := range cases {
		if got := IsOpen(h, c.t); got != c.want {
			t.Errorf("%s: IsOpen = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestIsOpenOvernightRespectsClosedYesterday(t *testing.T) {
	h := week("18:00", "01:00")
	h[1].Closed = true // Monday closed: Tuesday 00:30 must NOT inherit spill
	if IsOpen(h, at(2, 0, 30)) {
		t.Error("Tuesday 00:30 should be closed when Monday is a day off")
	}
}

func TestIsOpenInvalidTimes(t *testing.T) {
	var h [7]model.DayHours // all blank
	if IsOpen(h, at(1, 12, 0)) {
		t.Error("blank grid should read as closed")
	}
	h = week("9:30", "20:30") // malformed open (4 chars)
	if IsOpen(h, at(1, 12, 0)) {
		t.Error("malformed time should read as closed, not panic")
	}
}
