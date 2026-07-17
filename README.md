# ZatpatSite — Instant Website Generator

**Zatpat** (झटपट) = *instant*. In March 2024 Google killed free Business Profile
websites and orphaned millions of Indian small businesses overnight. ZatpatSite
gives them a way back: **one form → one genuinely beautiful, fully
self-contained one-page website** — services with ₹ prices, opening hours with
a live open-now badge, customer reviews, a Google Maps button and WhatsApp
call-to-actions, in English or Hindi, across four distinct designs.

**Who pays:** the shop owner — ₹2,999/yr for hosting + domain. The generated
site *is* the product: a single `index.html` with every byte of CSS inline,
inline SVG art, zero external requests. It works from a pen drive.

## Quickstart

```bash
make run        # serves on http://localhost:8104
```

Open http://localhost:8104 — the builder loads with two seeded demo sites
(Noor Beauty Salon, Jaipur and Annapurna Bhojanalaya, Indore) under **My sites**.

Type a business name, press **Fetch from Google** (deterministic mock — same
name + city always returns the same profile), watch the live preview, flip
themes/language, then **Save** or **Download index.html**.

```bash
make test       # go test ./...
make build      # bin/zatpatsite
```

## The four themes

| Theme | Design language |
|---|---|
| **Ivory** | Elegant serif, cream/charcoal, thin rules, dotted price leaders |
| **Bazaar** | Bold color-blocked, chunky sans, hard shadows, sticker cards |
| **Mint** | Airy rounded cards, fresh green, pill nav, soft shadows |
| **Nightshade** | Dark premium, gold accents, serif display, hairline ornaments |

These are different layout systems, not palette swaps.

## API

Base: `http://localhost:8104/api/v1`

| Method & path | What it does |
|---|---|
| `GET /health` | `{"status":"ok","gbp_provider":"mock","sites":N}` |
| `GET /meta` | Categories, themes (with swatches), languages |
| `POST /sites` | Create a site (JSON body; name required) |
| `GET /sites` | List sites, most recently updated first |
| `GET /sites/{id}` | Fetch one site |
| `PUT /sites/{id}` | Update a site |
| `DELETE /sites/{id}` | Delete a site |
| `GET /sites/{id}/preview` | The generated website (text/html) |
| `GET /sites/{id}/download` | Same HTML as `Content-Disposition: attachment` |
| `POST /render` | Generate HTML from an unsaved payload (live preview) |
| `POST /gbp/fetch` | Mock Google Business Profile autofill (`{name, city, category, lang}`) |

Storage: in-memory behind a mutex, JSON-snapshotted to `./data/store.json`
after every write and reloaded on boot. No database.

## Generated-site guarantees (enforced by tests)

- Valid `LocalBusiness` JSON-LD with `openingHoursSpecification`, OpenGraph tags, meta description.
- **Zero external resource loads** — the only outbound links are `wa.me`, `tel:` and Google Maps hrefs.
- Live open-now badge (overnight hours like 18:00–01:00 handled), Indian ₹ formatting (`₹1.2 Cr / ₹36.5 L / ₹12,500`).
- All 4 themes × 2 languages render non-empty and pairwise distinct; user content is HTML-escaped.

## Upgrade to live integrations

| Env var | Default | Live behaviour |
|---|---|---|
| `PORT` | `8104` | HTTP port |
| `GBP_PROVIDER` | `mock` | `live` would call Google Places API (interface `internal/gbp.Provider` is ready; only the mock ships — the server logs and falls back) |
| `GOOGLE_PLACES_API_KEY` | — | Required once `GBP_PROVIDER=live` exists |

Zero third-party Go dependencies — standard library only.
