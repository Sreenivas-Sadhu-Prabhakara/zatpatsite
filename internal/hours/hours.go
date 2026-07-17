// Package hours implements the open-now rules used both server-side (initial
// badge state) and mirrored by the small inline script in generated sites.
package hours

import (
	"time"

	"zatpatsite/internal/model"
)

// IST is Indian Standard Time (UTC+5:30). Generated sites belong to Indian
// businesses, so "open now" is always evaluated in IST on the server.
var IST = time.FixedZone("IST", 5*3600+1800)

// minutes parses "HH:MM" into minutes since midnight. Returns -1 if invalid.
func minutes(s string) int {
	if len(s) != 5 || s[2] != ':' {
		return -1
	}
	h := int(s[0]-'0')*10 + int(s[1]-'0')
	m := int(s[3]-'0')*10 + int(s[4]-'0')
	if s[0] < '0' || s[0] > '9' || s[1] < '0' || s[1] > '9' ||
		s[3] < '0' || s[3] > '9' || s[4] < '0' || s[4] > '9' ||
		h > 23 || m > 59 {
		return -1
	}
	return h*60 + m
}

// IsOpen reports whether the business is open at t according to the weekly
// grid. Handles overnight ranges (close <= open means the window spills past
// midnight into the next day).
func IsOpen(h [7]model.DayHours, t time.Time) bool {
	day := int(t.Weekday()) // 0 = Sunday, matching the grid
	cur := t.Hour()*60 + t.Minute()

	// Today's window.
	today := h[day]
	if !today.Closed {
		o, c := minutes(today.Open), minutes(today.Close)
		if o >= 0 && c >= 0 {
			if c > o {
				if cur >= o && cur < c {
					return true
				}
			} else if cur >= o { // overnight: open until past midnight
				return true
			}
		}
	}

	// Yesterday's overnight spill (e.g. 18:00-01:00 keeps us open at 00:30).
	yest := h[(day+6)%7]
	if !yest.Closed {
		o, c := minutes(yest.Open), minutes(yest.Close)
		if o >= 0 && c >= 0 && c <= o && cur < c {
			return true
		}
	}
	return false
}

// TodayLabel returns today's index (0=Sunday) in IST for highlighting.
func TodayLabel(now time.Time) int {
	return int(now.In(IST).Weekday())
}
