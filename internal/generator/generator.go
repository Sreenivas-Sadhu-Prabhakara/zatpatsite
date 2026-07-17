// Package generator turns a model.Site into one fully self-contained
// index.html: every byte of CSS inline, inline SVG art, zero external
// resource loads. The only outbound hrefs are wa.me, tel: and Google Maps.
package generator

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"zatpatsite/internal/hours"
	"zatpatsite/internal/model"
	"zatpatsite/internal/money"
)

func esc(s string) string { return html.EscapeString(s) }

// waNumber renders the digits used in a wa.me link (adds 91 for 10-digit numbers).
func waNumber(digits string) string {
	if len(digits) == 10 {
		return "91" + digits
	}
	return digits
}

// FormatPhone renders a bare 10-digit number as +91-XXXXXXXXXX for display.
func FormatPhone(digits string) string {
	if len(digits) == 10 {
		return "+91-" + digits
	}
	return digits
}

// WALink builds the wa.me deep link with the prefilled message.
func WALink(s model.Site) string {
	str := langFor(s.Lang)
	msg := strings.ReplaceAll(str.WAMessage, "{name}", s.Name)
	return "https://wa.me/" + waNumber(s.WhatsApp) + "?text=" + url.QueryEscape(msg)
}

// MapsHref returns the owner-supplied maps URL, or a search link built from
// the name + address.
func MapsHref(s model.Site) string {
	u := strings.TrimSpace(s.MapsURL)
	if strings.HasPrefix(u, "https://") || strings.HasPrefix(u, "http://") {
		return u
	}
	q := strings.TrimSpace(s.Name + " " + s.Address + " " + s.City)
	return "https://www.google.com/maps/search/?api=1&query=" + url.QueryEscape(q)
}

var schemaTypes = map[string]string{
	"salon": "HairSalon", "restaurant": "Restaurant", "kirana": "GroceryStore",
	"coaching": "EducationalOrganization", "clinic": "MedicalClinic",
	"boutique": "ClothingStore", "gym": "ExerciseGym", "bakery": "Bakery",
}

var schemaDays = [7]string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}

type ldSpec struct {
	Type      string `json:"@type"`
	DayOfWeek string `json:"dayOfWeek"`
	Opens     string `json:"opens"`
	Closes    string `json:"closes"`
}

type ldAddress struct {
	Type            string `json:"@type"`
	StreetAddress   string `json:"streetAddress"`
	AddressLocality string `json:"addressLocality"`
	AddressCountry  string `json:"addressCountry"`
}

type ldRating struct {
	Type        string  `json:"@type"`
	RatingValue float64 `json:"ratingValue"`
	ReviewCount int     `json:"reviewCount"`
}

type ldBusiness struct {
	Context         string    `json:"@context"`
	Type            string    `json:"@type"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Telephone       string    `json:"telephone,omitempty"`
	Address         ldAddress `json:"address"`
	OpeningHours    []ldSpec  `json:"openingHoursSpecification"`
	AggregateRating *ldRating `json:"aggregateRating,omitempty"`
	PriceRange      string    `json:"priceRange,omitempty"`
}

// JSONLD builds the LocalBusiness structured data block.
func JSONLD(s model.Site, description string) (string, error) {
	typ, ok := schemaTypes[s.Category]
	if !ok {
		typ = "LocalBusiness"
	}
	specs := make([]ldSpec, 0, 7)
	for i, d := range s.Hours {
		if d.Closed || d.Open == "" || d.Close == "" {
			continue
		}
		specs = append(specs, ldSpec{
			Type: "OpeningHoursSpecification", DayOfWeek: schemaDays[i],
			Opens: d.Open, Closes: d.Close,
		})
	}
	biz := ldBusiness{
		Context: "https://schema.org", Type: typ,
		Name: s.Name, Description: description,
		Telephone: FormatPhone(s.Phone),
		Address: ldAddress{
			Type: "PostalAddress", StreetAddress: s.Address,
			AddressLocality: s.City, AddressCountry: "IN",
		},
		OpeningHours: specs,
		PriceRange:   "₹₹",
	}
	if s.Rating > 0 && len(s.Reviews) > 0 {
		biz.AggregateRating = &ldRating{
			Type: "AggregateRating", RatingValue: s.Rating, ReviewCount: len(s.Reviews),
		}
	}
	b, err := json.MarshalIndent(biz, "", " ")
	return string(b), err
}

// stars renders n filled star glyphs (rounded from the rating).
func stars(rating float64) string {
	n := int(rating + 0.5)
	if n < 1 {
		n = 1
	}
	if n > 5 {
		n = 5
	}
	return strings.Repeat(starGlyph, n)
}

// Generate renders the complete standalone index.html for a site.
// now is used for the server-rendered open/closed badge (evaluated in IST);
// the inline script keeps it live on the visitor's device.
func Generate(s model.Site, now time.Time) string {
	s.Normalize()
	str := langFor(s.Lang)
	catLabel := CategoryLabel(s.Category, s.Lang)
	tagline := Tagline(s.Category, s.Lang)
	about := s.About
	if about == "" {
		about = AutoAbout(s.Name, s.City, s.Category, s.Lang)
	}
	svcHeading := servicesHeading(s.Category, s.Lang)
	desc := strings.NewReplacer(
		"{name}", s.Name, "{category}", catLabel,
		"{city}", s.City, "{tagline}", tagline,
	).Replace(str.MetaDesc)

	nowIST := now.In(hours.IST)
	isOpen := hours.IsOpen(s.Hours, nowIST)
	today := int(nowIST.Weekday())
	badgeClass, badgeText := "closed", str.ClosedNow
	if isOpen {
		badgeClass, badgeText = "open", str.OpenNow
	}

	wa := WALink(s)
	tel := "tel:+91" + s.Phone
	maps := MapsHref(s)
	ld, _ := JSONLD(s, desc)
	title := s.Name + " — " + catLabel
	if s.City != "" {
		title += ", " + s.City
	}

	var b strings.Builder
	b.Grow(64 << 10)
	w := func(parts ...string) {
		for _, p := range parts {
			b.WriteString(p)
		}
	}

	// ---- head ----
	w(`<!DOCTYPE html>
<html lang="`, str.LangCode, `">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>`, esc(title), `</title>
<meta name="description" content="`, esc(desc), `">
<meta property="og:title" content="`, esc(title), `">
<meta property="og:description" content="`, esc(desc), `">
<meta property="og:type" content="business.business">
<meta property="og:locale" content="`, map[string]string{"en": "en_IN", "hi": "hi_IN"}[str.LangCode], `">
<script type="application/ld+json">
`, ld, `
</script>
<style>`, themeCSS(s.Theme, s.Category), `</style>
</head>
<body class="t-`, s.Theme, ` cat-`, s.Category, `">
`)

	// ---- nav ----
	w(`<nav class="nav"><div class="nav-inner"><a class="brand" href="#top">`, esc(s.Name), `</a><div class="navlinks">`)
	if len(s.Services) > 0 {
		w(`<a href="#services">`, esc(svcHeading), `</a>`)
	}
	w(`<a href="#about">`, str.HeadingAbout, `</a><a href="#hours">`, str.HeadingHours, `</a>`)
	if len(s.Reviews) > 0 {
		w(`<a href="#reviews">`, str.HeadingReviews, `</a>`)
	}
	w(`<a href="#visit">`, str.HeadingVisit, `</a></div></div></nav>
`)

	// ---- hero ----
	w(`<header class="hero" id="top"><div class="hero-inner">
<div class="emblem-wrap">`, Emblem(s.Category), `</div>
<p class="kicker">`, esc(catLabel))
	if s.City != "" {
		w(` · `, esc(s.City))
	}
	w(`</p>
<h1>`, esc(s.Name), `</h1>
<p class="tagline">`, esc(tagline), `</p>
<div class="orn" aria-hidden="true"><span></span></div>
<div class="badge-row"><span id="open-badge" class="badge `, badgeClass, `">`, badgeText, `</span>`)
	if s.Rating > 0 {
		w(`<span class="rating-chip">`, starGlyph, ` `, fmt.Sprintf("%.1f", s.Rating), ` · `, str.ReviewsOn, `</span>`)
	}
	w(`</div>
<div class="cta-row"><a class="btn btn-wa" href="`, esc(wa), `" rel="noopener">`, waGlyph, ` `, str.WhatsAppCTA, `</a>`)
	if s.Phone != "" {
		w(`<a class="btn btn-call" href="`, esc(tel), `">`, phoneGlyph, ` `, str.CallCTA, `</a>`)
	}
	w(`</div>
</div></header>
<main>
`)

	// ---- services ----
	if len(s.Services) > 0 {
		w(`<section class="sec" id="services"><div class="wrap"><h2 class="sec-title">`, esc(svcHeading), `</h2><ul class="svc">`)
		for _, sv := range s.Services {
			price := money.FormatINR(sv.Price)
			if sv.Price == 0 {
				if s.Lang == "hi" {
					price = "मुफ़्त"
				} else {
					price = "Free"
				}
			}
			w(`<li class="svc-item"><span class="svc-name">`, esc(sv.Name), `</span><span class="leader"></span><span class="svc-price">`, esc(price), `</span></li>`)
		}
		w(`</ul></div></section>
`)
	}

	// ---- about ----
	w(`<section class="sec sec-alt" id="about"><div class="wrap"><h2 class="sec-title">`, str.HeadingAbout, `</h2><p class="about-text">`, esc(about), `</p></div></section>
`)

	// ---- hours ----
	w(`<section class="sec" id="hours"><div class="wrap"><h2 class="sec-title">`, str.HeadingHours, `</h2>`)
	w(`<p class="hours-note"><span id="hours-badge" class="badge `, badgeClass, `">`, badgeText, `</span></p>`)
	w(`<div class="hours-card"><table class="hours"><tbody>`)
	for _, di := range [7]int{1, 2, 3, 4, 5, 6, 0} { // Monday-first display
		d := s.Hours[di]
		cls := ""
		if di == today {
			cls = ` class="today"`
		}
		hiddenAttr := " hidden"
		if di == today {
			hiddenAttr = ""
		}
		w(`<tr data-day="`, fmt.Sprintf("%d", di), `"`, cls, `><td class="day">`, str.Days[di],
			`<span class="today-tag"`, hiddenAttr, `> · `, str.Today, `</span></td><td class="time">`)
		if d.Closed || d.Open == "" || d.Close == "" {
			w(`<span class="day-closed">`, str.ClosedDay, `</span>`)
		} else {
			w(esc(d.Open), ` – `, esc(d.Close))
		}
		w(`</td></tr>`)
	}
	w(`</tbody></table></div></div></section>
`)

	// ---- reviews ----
	if len(s.Reviews) > 0 {
		w(`<section class="sec sec-alt" id="reviews"><div class="wrap"><h2 class="sec-title">`, str.HeadingReviews, `</h2><div class="reviews">`)
		for _, r := range s.Reviews {
			w(`<figure class="review"><div class="stars">`, stars(r.Rating), `</div><blockquote>“`, esc(r.Text), `”</blockquote><figcaption>`, esc(r.Author), `</figcaption></figure>`)
		}
		w(`</div></div></section>
`)
	}

	// ---- visit ----
	w(`<section class="sec" id="visit"><div class="wrap"><div class="visit-block"><h2 class="sec-title">`, str.HeadingVisit, `</h2>`)
	if s.Address != "" {
		w(`<address>`, esc(s.Address), `</address>`)
	}
	if s.Phone != "" {
		w(`<p class="visit-phone">`, str.PhoneLabel, `: `, esc(FormatPhone(s.Phone)), `</p>`)
	}
	w(`<div class="cta-row"><a class="btn" href="`, esc(maps), `" target="_blank" rel="noopener">`, pinGlyph, ` `, str.DirectionsCTA, `</a><a class="btn btn-wa" href="`, esc(wa), `" rel="noopener">`, waGlyph, ` `, str.ChatOnWhatsApp, `</a></div>`)
	w(`</div></div></section>
</main>
`)

	// ---- footer + floating WA ----
	year := nowIST.Year()
	w(`<footer class="footer"><p>© `, fmt.Sprintf("%d", year), ` `, esc(s.Name), `</p><p class="made">`, str.MadeWith, `</p></footer>
<a class="wa-float" href="`, esc(wa), `" rel="noopener" aria-label="WhatsApp">`, waGlyph, `</a>
`)

	// ---- inline script: live open-now ----
	hoursJSON, _ := json.Marshal(s.Hours)
	w(`<script>
(function(){
var H=`, string(hoursJSON), `;
var OPEN=`, jsString(str.OpenNow), `,CLOSED=`, jsString(str.ClosedNow), `;
function m(s){if(!s||s.length!==5)return -1;var h=+s.slice(0,2),mm=+s.slice(3);if(isNaN(h)||isNaN(mm))return -1;return h*60+mm}
function isOpen(d){
 var day=d.getDay(),cur=d.getHours()*60+d.getMinutes();
 var t=H[day];
 if(t&&!t.closed){var o=m(t.open),c=m(t.close);
  if(o>=0&&c>=0){if(c>o){if(cur>=o&&cur<c)return true}else if(cur>=o)return true}}
 var y=H[(day+6)%7];
 if(y&&!y.closed){var o2=m(y.open),c2=m(y.close);
  if(o2>=0&&c2>=0&&c2<=o2&&cur<c2)return true}
 return false}
function paint(){
 var open=isOpen(new Date());
 ["open-badge","hours-badge"].forEach(function(id){
  var el=document.getElementById(id);if(!el)return;
  el.textContent=open?OPEN:CLOSED;
  el.classList.toggle("open",open);el.classList.toggle("closed",!open)});
 var today=new Date().getDay();
 document.querySelectorAll(".hours tr").forEach(function(r){
  var isToday=+r.getAttribute("data-day")===today;
  r.classList.toggle("today",isToday);
  var tag=r.querySelector(".today-tag");if(tag)tag.hidden=!isToday})}
paint();setInterval(paint,60000);
})();
</script>
</body>
</html>
`)
	return b.String()
}

// jsString safely embeds a UI string in the inline script.
func jsString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
