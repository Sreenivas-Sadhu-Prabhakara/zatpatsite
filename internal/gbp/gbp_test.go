package gbp

import (
	"regexp"
	"strings"
	"testing"
)

func TestFetchDeterministic(t *testing.T) {
	m := Mock{}
	a, err := m.Fetch("Noor Beauty Salon", "Jaipur", "salon")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := m.Fetch("Noor Beauty Salon", "Jaipur", "salon")
	if a.Address != b.Address || a.Phone != b.Phone || a.Rating != b.Rating ||
		len(a.Reviews) != len(b.Reviews) || a.Reviews[0] != b.Reviews[0] {
		t.Error("same input must yield identical profile")
	}
	c, _ := m.Fetch("Gupta Kirana Store", "Indore", "kirana")
	if a.Address == c.Address && a.Phone == c.Phone {
		t.Error("different inputs should yield different profiles")
	}
}

func TestFetchShape(t *testing.T) {
	p, err := Mock{}.Fetch("Iyer Tiffin Centre", "Chennai", "restaurant")
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^98\d{8}$`).MatchString(p.Phone) {
		t.Errorf("phone %q must be a 10-digit 98-series mobile", p.Phone)
	}
	if !strings.Contains(p.Address, "Chennai") {
		t.Errorf("address %q must contain the city", p.Address)
	}
	if p.Rating < 3.9 || p.Rating > 4.9 {
		t.Errorf("rating %v out of believable band", p.Rating)
	}
	if len(p.Services) < 5 {
		t.Errorf("expected at least 5 services, got %d", len(p.Services))
	}
	if len(p.Reviews) != 3 {
		t.Errorf("expected 3 reviews, got %d", len(p.Reviews))
	}
	open := false
	for _, d := range p.Hours {
		if !d.Closed && d.Open != "" {
			open = true
		}
	}
	if !open {
		t.Error("hours grid should have open days")
	}
	if !strings.HasPrefix(p.MapsURL, "https://www.google.com/maps/") {
		t.Errorf("unexpected maps url %q", p.MapsURL)
	}
}

func TestFetchInfersCategory(t *testing.T) {
	p, _ := Mock{}.Fetch("Some Shop", "Pune", "")
	q, _ := Mock{}.Fetch("Some Shop", "Pune", "not-a-category")
	if p.Category != q.Category {
		t.Error("category inference must be deterministic")
	}
	if p.Category == "" {
		t.Error("category must be inferred when absent")
	}
}
