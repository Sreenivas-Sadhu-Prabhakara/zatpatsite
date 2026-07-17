package generator

import "strings"

// glyphs holds the inner stroke markup of each category emblem in a 48×48 box.
// Everything is stroke-based so themes can colour it with currentColor.
var glyphs = map[string]string{
	"salon":      `<circle cx="13" cy="13" r="5"/><circle cx="13" cy="35" r="5"/><path d="M17 16 38 35m-21-3L38 13m0 22 4 3M38 13l4-3"/>`,
	"restaurant": `<path d="M8 26h32a16 16 0 0 1-32 0Z"/><path d="M18 20c-2-3 2-5 0-8m10 8c-2-3 2-5 0-8"/><path d="M14 42h20"/>`,
	"kirana":     `<path d="M9 20h30l-4 18H13Z"/><path d="M18 20l6-11 6 11"/><path d="M19 26v6m5-6v6m5-6v6"/>`,
	"coaching":   `<path d="M24 12c-4-3-10-4-15-2v26c5-2 11-1 15 2 4-3 10-4 15-2V10c-5-2-11-1-15 2Z"/><path d="M24 12v26"/>`,
	"clinic":     `<path d="M24 6l14 5v10c0 9-6 17-14 21C16 38 10 30 10 21V11Z"/><path d="M24 16v14m-7-7h14"/>`,
	"boutique":   `<path d="M28 8a4 4 0 1 0-4 4v5"/><path d="M24 17 6 30c-2 1-1 4 1 4h34c2 0 3-3 1-4L24 17Z"/>`,
	"gym":        `<path d="M6 19h5v10H6Zm31 0h5v10h-5Zm-24-4h5v18h-5Zm17 0h5v18h-5Z"/><path d="M18 24h12"/>`,
	"bakery":     `<path d="M13 21c-2-11 24-11 22 0"/><path d="M10 21h28l-4 6H14Z"/><path d="M15 27l3 13h12l3-13"/>`,
}

func glyph(category string) string {
	if g, ok := glyphs[category]; ok {
		return g
	}
	return glyphs["salon"]
}

// Emblem returns an inline SVG for the category, coloured via currentColor.
func Emblem(category string) string {
	return `<svg class="emblem" viewBox="0 0 48 48" width="56" height="56" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">` +
		glyph(category) + `</svg>`
}

// waGlyph is a simple chat-bubble-with-handset mark for WhatsApp CTAs.
const waGlyph = `<svg viewBox="0 0 48 48" width="20" height="20" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M24 6a17 17 0 0 0-14.5 26L7 42l10.4-2.4A17 17 0 1 0 24 6Z"/><path d="M17 19c0 8 5 12 13 13l2-4-5-3-2 2c-2-1-4-3-5-5l2-2-3-5Z"/></svg>`

// phoneGlyph for call CTAs.
const phoneGlyph = `<svg viewBox="0 0 48 48" width="18" height="18" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M10 8h8l3 9-4 3c1.5 4 4.5 7 9 9l3-4 9 3v8c0 2-2 4-4 4C19 39 9 29 8 12c0-2 1-4 2-4Z"/></svg>`

// pinGlyph for the directions button.
const pinGlyph = `<svg viewBox="0 0 48 48" width="18" height="18" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M24 42S10 28.6 10 18a14 14 0 0 1 28 0c0 10.6-14 24-14 24Z"/><circle cx="24" cy="18" r="5"/></svg>`

// starGlyph (filled) for ratings.
const starGlyph = `<svg viewBox="0 0 24 24" width="14" height="14" fill="currentColor" aria-hidden="true"><path d="M12 2l2.9 6.3 6.9.8-5.1 4.7 1.4 6.8L12 17l-6.1 3.6 1.4-6.8L2.2 9.1l6.9-.8Z"/></svg>`

// svgEncode percent-encodes an SVG document for use inside a CSS data: URI.
func svgEncode(svg string) string {
	r := strings.NewReplacer(
		"%", "%25",
		"#", "%23",
		"<", "%3C",
		">", "%3E",
		"\"", "'",
		"{", "%7B",
		"}", "%7D",
		"\n", "",
	)
	return r.Replace(svg)
}

// PatternDataURI builds a subtle repeating background of the category glyph —
// the "decorative treatment" behind each hero. Pure SVG, zero requests.
func PatternDataURI(category, color string, opacity string) string {
	inner := glyph(category)
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="132" height="132" viewBox="0 0 132 132">` +
		`<g fill="none" stroke="` + color + `" stroke-opacity="` + opacity + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" transform="translate(10,10)">` + inner + `</g>` +
		`<g fill="none" stroke="` + color + `" stroke-opacity="` + opacity + `" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" transform="translate(76,76) scale(0.7)">` + inner + `</g>` +
		`</svg>`
	return "data:image/svg+xml," + svgEncode(svg)
}
