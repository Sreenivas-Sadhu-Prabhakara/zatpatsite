// Package model defines the core domain types shared by the store,
// generator and API layers.
package model

import (
	"strings"
	"time"
)

// Categories supported by ZatpatSite. Order matters for UI display.
var Categories = []string{
	"salon", "restaurant", "kirana", "coaching", "clinic", "boutique", "gym", "bakery",
}

// Themes supported by the generator.
var Themes = []string{"ivory", "bazaar", "mint", "nightshade"}

// Langs supported by the generator.
var Langs = []string{"en", "hi"}

// DayHours describes opening hours for a single day.
// Times are "HH:MM" 24-hour strings. Index 0 = Sunday (matches JS Date.getDay()).
type DayHours struct {
	Closed bool   `json:"closed"`
	Open   string `json:"open"`
	Close  string `json:"close"`
}

// Service is a single line item on the services/menu section.
type Service struct {
	Name  string `json:"name"`
	Price int64  `json:"price"` // whole rupees
}

// Review is a customer review shown on the generated site.
type Review struct {
	Author string  `json:"author"`
	Rating float64 `json:"rating"`
	Text   string  `json:"text"`
}

// Site is everything needed to generate one website.
type Site struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Category  string      `json:"category"`
	City      string      `json:"city"`
	Address   string      `json:"address"`
	Phone     string      `json:"phone"`
	WhatsApp  string      `json:"whatsapp"`
	Hours     [7]DayHours `json:"hours"`
	Services  []Service   `json:"services"`
	About     string      `json:"about"`
	MapsURL   string      `json:"mapsUrl"`
	Theme     string      `json:"theme"`
	Lang      string      `json:"lang"`
	Rating    float64     `json:"rating"`
	Reviews   []Review    `json:"reviews"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// ValidCategory reports whether c is a known category.
func ValidCategory(c string) bool {
	for _, k := range Categories {
		if k == c {
			return true
		}
	}
	return false
}

// ValidTheme reports whether t is a known theme.
func ValidTheme(t string) bool {
	for _, k := range Themes {
		if k == t {
			return true
		}
	}
	return false
}

// ValidLang reports whether l is a supported language.
func ValidLang(l string) bool {
	for _, k := range Langs {
		if k == l {
			return true
		}
	}
	return false
}

// DefaultHours returns a sensible default week: 09:30-20:30 every day,
// Sunday closing early at 14:00.
func DefaultHours() [7]DayHours {
	var h [7]DayHours
	for i := range h {
		h[i] = DayHours{Open: "09:30", Close: "20:30"}
	}
	h[0] = DayHours{Open: "10:00", Close: "14:00"} // Sunday half day
	return h
}

// HoursEmpty reports whether the hours grid is entirely unset.
func HoursEmpty(h [7]DayHours) bool {
	for _, d := range h {
		if d.Closed || d.Open != "" || d.Close != "" {
			return false
		}
	}
	return true
}

// NormalizePhone keeps digits only and trims a leading country code so the
// result is a bare 10-digit Indian number where possible.
func NormalizePhone(p string) string {
	var b strings.Builder
	for _, r := range p {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	d := b.String()
	if len(d) == 12 && strings.HasPrefix(d, "91") {
		d = d[2:]
	}
	if len(d) == 11 && strings.HasPrefix(d, "0") {
		d = d[1:]
	}
	return d
}

// Normalize fills defaults and cleans up user input in place.
func (s *Site) Normalize() {
	s.Name = strings.TrimSpace(s.Name)
	s.City = strings.TrimSpace(s.City)
	s.Address = strings.TrimSpace(s.Address)
	s.About = strings.TrimSpace(s.About)
	s.MapsURL = strings.TrimSpace(s.MapsURL)
	s.Category = strings.ToLower(strings.TrimSpace(s.Category))
	s.Theme = strings.ToLower(strings.TrimSpace(s.Theme))
	s.Lang = strings.ToLower(strings.TrimSpace(s.Lang))
	if !ValidCategory(s.Category) {
		s.Category = "salon"
	}
	if !ValidTheme(s.Theme) {
		s.Theme = "ivory"
	}
	if !ValidLang(s.Lang) {
		s.Lang = "en"
	}
	if HoursEmpty(s.Hours) {
		s.Hours = DefaultHours()
	}
	s.Phone = NormalizePhone(s.Phone)
	if s.WhatsApp == "" {
		s.WhatsApp = s.Phone
	} else {
		s.WhatsApp = NormalizePhone(s.WhatsApp)
	}
	clean := make([]Service, 0, len(s.Services))
	for _, sv := range s.Services {
		sv.Name = strings.TrimSpace(sv.Name)
		if sv.Name != "" {
			clean = append(clean, sv)
		}
	}
	s.Services = clean
}
