package generator

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"

	"zatpatsite/internal/hours"
	"zatpatsite/internal/model"
)

func sampleSite() model.Site {
	h := model.DefaultHours()
	h[1].Closed = true // Monday off
	return model.Site{
		ID: "test1", Name: "Noor Beauty Salon", Category: "salon",
		City: "Jaipur", Address: "Shop 12, MG Road, Shivaji Nagar, Jaipur - 302001",
		Phone: "9812345678", WhatsApp: "9812345678",
		Hours: h,
		Services: []model.Service{
			{Name: "Haircut (Women)", Price: 550},
			{Name: "Bridal Makeup", Price: 15000},
			{Name: "Hair Spa", Price: 900},
		},
		Rating: 4.6,
		Reviews: []model.Review{
			{Author: "Priya Sharma", Rating: 5, Text: "Best haircut in years."},
			{Author: "Kavita Reddy", Rating: 4, Text: "Always on time."},
		},
		Theme: "ivory", Lang: "en",
	}
}

// tuesday noon IST, within default hours.
var openNow = time.Date(2026, 7, 14, 12, 0, 0, 0, hours.IST)

var ldRe = regexp.MustCompile(`(?s)<script type="application/ld\+json">\s*(.*?)\s*</script>`)

func TestJSONLDValidAndCorrect(t *testing.T) {
	out := Generate(sampleSite(), openNow)
	m := ldRe.FindStringSubmatch(out)
	if m == nil {
		t.Fatal("no JSON-LD block in output")
	}
	var ld struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Name    string `json:"name"`
		Address struct {
			Locality string `json:"addressLocality"`
			Country  string `json:"addressCountry"`
		} `json:"address"`
		Specs []struct {
			DayOfWeek string `json:"dayOfWeek"`
			Opens     string `json:"opens"`
			Closes    string `json:"closes"`
		} `json:"openingHoursSpecification"`
		Agg *struct {
			Value float64 `json:"ratingValue"`
			Count int     `json:"reviewCount"`
		} `json:"aggregateRating"`
	}
	if err := json.Unmarshal([]byte(m[1]), &ld); err != nil {
		t.Fatalf("JSON-LD does not parse: %v", err)
	}
	if ld.Name != "Noor Beauty Salon" {
		t.Errorf("name = %q", ld.Name)
	}
	if ld.Type != "HairSalon" {
		t.Errorf("@type = %q, want HairSalon", ld.Type)
	}
	if ld.Address.Locality != "Jaipur" || ld.Address.Country != "IN" {
		t.Errorf("address wrong: %+v", ld.Address)
	}
	if len(ld.Specs) != 6 { // 7 days minus Monday off
		t.Fatalf("openingHoursSpecification count = %d, want 6", len(ld.Specs))
	}
	for _, sp := range ld.Specs {
		if sp.DayOfWeek == "Monday" {
			t.Error("closed Monday must not appear in openingHoursSpecification")
		}
		if sp.DayOfWeek == "Tuesday" && (sp.Opens != "09:30" || sp.Closes != "20:30") {
			t.Errorf("Tuesday spec = %+v", sp)
		}
		if sp.DayOfWeek == "Sunday" && (sp.Opens != "10:00" || sp.Closes != "14:00") {
			t.Errorf("Sunday spec = %+v", sp)
		}
	}
	if ld.Agg == nil || ld.Agg.Value != 4.6 || ld.Agg.Count != 2 {
		t.Errorf("aggregateRating = %+v", ld.Agg)
	}
}

// urlRe finds every absolute http(s) URL in the document.
var urlRe = regexp.MustCompile(`https?://[^"'\s)<>\\]+`)

// allowed prefixes: WhatsApp + Maps action links, plus two identifiers that
// are never fetched by a browser (SVG xmlns inside data: URIs, JSON-LD @context).
var allowedURL = []string{
	"https://wa.me/",
	"https://www.google.com/maps",
	"https://maps.google.com",
	"http://www.w3.org/2000/svg",
	"https://schema.org",
}

func externalViolations(doc string) []string {
	var bad []string
	for _, u := range urlRe.FindAllString(doc, -1) {
		ok := false
		for _, p := range allowedURL {
			if strings.HasPrefix(u, p) {
				ok = true
				break
			}
		}
		if !ok {
			bad = append(bad, u)
		}
	}
	return bad
}

func TestZeroExternalResourceLoads(t *testing.T) {
	for _, theme := range model.Themes {
		for _, lang := range model.Langs {
			for _, cat := range model.Categories {
				s := sampleSite()
				s.Theme, s.Lang, s.Category = theme, lang, cat
				out := Generate(s, openNow)
				if bad := externalViolations(out); len(bad) > 0 {
					t.Errorf("%s/%s/%s loads external URLs: %v", theme, lang, cat, bad)
				}
				// No resource-loading tags at all.
				for _, marker := range []string{"<img", "<link", `src="`, "src='", "@import", "url(http"} {
					if strings.Contains(out, marker) {
						t.Errorf("%s/%s/%s contains forbidden %q", theme, lang, cat, marker)
					}
				}
			}
		}
	}
}

func TestThemesAndLanguagesDistinctNonEmpty(t *testing.T) {
	outputs := map[string]string{}
	for _, theme := range model.Themes {
		for _, lang := range model.Langs {
			s := sampleSite()
			s.Theme, s.Lang = theme, lang
			out := Generate(s, openNow)
			if len(out) < 5000 {
				t.Errorf("%s/%s output suspiciously small: %d bytes", theme, lang, len(out))
			}
			if !strings.Contains(out, "Noor Beauty Salon") {
				t.Errorf("%s/%s missing business name", theme, lang)
			}
			outputs[theme+"/"+lang] = out
		}
	}
	keys := make([]string, 0, len(outputs))
	for k := range outputs {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if outputs[keys[i]] == outputs[keys[j]] {
				t.Errorf("outputs identical: %s vs %s", keys[i], keys[j])
			}
		}
	}
	if !strings.Contains(outputs["ivory/hi"], "हमारे बारे में") {
		t.Error("Hindi output missing Hindi about heading")
	}
	if !strings.Contains(outputs["ivory/en"], "About Us") {
		t.Error("English output missing About Us heading")
	}
	if !strings.Contains(outputs["ivory/hi"], `lang="hi"`) {
		t.Error("Hindi output must set lang attribute")
	}
}

func TestOpenBadgeServerRendered(t *testing.T) {
	s := sampleSite()
	if out := Generate(s, openNow); !strings.Contains(out, ">Open now<") {
		t.Error("expected Open now badge at Tuesday noon")
	}
	night := time.Date(2026, 7, 14, 23, 0, 0, 0, hours.IST)
	if out := Generate(s, night); !strings.Contains(out, ">Closed now<") {
		t.Error("expected Closed now badge at 23:00")
	}
	s.Lang = "hi"
	if out := Generate(s, openNow); !strings.Contains(out, ">अभी खुला है<") {
		t.Error("expected Hindi open badge")
	}
}

func TestINRAndWhatsAppLink(t *testing.T) {
	s := sampleSite()
	s.Services = append(s.Services, model.Service{Name: "Yearly Package", Price: 130000})
	out := Generate(s, openNow)
	if !strings.Contains(out, "₹550") {
		t.Error("missing plain rupee price ₹550")
	}
	if !strings.Contains(out, "₹15,000") {
		t.Error("missing Indian-grouped ₹15,000")
	}
	if !strings.Contains(out, "₹1.3 L") {
		t.Error("missing lakh-formatted ₹1.3 L")
	}
	if !strings.Contains(out, "https://wa.me/919812345678?text=") {
		t.Error("missing wa.me deep link with prefilled text")
	}
	if !strings.Contains(out, "tel:+919812345678") {
		t.Error("missing tel: link")
	}
}

func TestUserContentEscaped(t *testing.T) {
	s := sampleSite()
	s.Name = `Sharma <script>alert(1)</script> & Sons`
	s.About = `custom "about" <b>text</b>`
	out := Generate(s, openNow)
	if strings.Contains(out, "<script>alert(1)</script>") {
		t.Error("script injection not escaped")
	}
	if !strings.Contains(out, "Sharma &lt;script&gt;") {
		t.Error("expected escaped name in output")
	}
	if !strings.Contains(out, "custom &#34;about&#34; &lt;b&gt;text&lt;/b&gt;") {
		t.Error("expected escaped about text")
	}
}

func TestAutoAboutFillsEmpty(t *testing.T) {
	s := sampleSite()
	s.About = ""
	out := Generate(s, openNow)
	want := AutoAbout("Noor Beauty Salon", "Jaipur", "salon", "en")
	if !strings.Contains(out, want[:40]) {
		t.Error("empty about should be auto-filled from the category template")
	}
	s.About = "Our own story."
	if out := Generate(s, openNow); !strings.Contains(out, "Our own story.") {
		t.Error("custom about must be preserved")
	}
}
