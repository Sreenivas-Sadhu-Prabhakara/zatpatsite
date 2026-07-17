// ZatpatSite server: JSON API + embedded builder UI in one binary.
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"zatpatsite/internal/gbp"
	"zatpatsite/internal/generator"
	"zatpatsite/internal/model"
	"zatpatsite/internal/store"
	"zatpatsite/web"
)

const defaultPort = "8104"

type app struct {
	store    *store.Store
	provider gbp.Provider
}

func main() {
	st, err := store.Open("data/store.json")
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	a := &app{store: st, provider: pickProvider()}
	a.seed()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", a.health)
	mux.HandleFunc("GET /api/v1/meta", a.meta)
	mux.HandleFunc("POST /api/v1/sites", a.createSite)
	mux.HandleFunc("GET /api/v1/sites", a.listSites)
	mux.HandleFunc("GET /api/v1/sites/{id}", a.getSite)
	mux.HandleFunc("PUT /api/v1/sites/{id}", a.updateSite)
	mux.HandleFunc("DELETE /api/v1/sites/{id}", a.deleteSite)
	mux.HandleFunc("GET /api/v1/sites/{id}/preview", a.previewSite)
	mux.HandleFunc("GET /api/v1/sites/{id}/download", a.downloadSite)
	mux.HandleFunc("POST /api/v1/render", a.render)
	mux.HandleFunc("POST /api/v1/gbp/fetch", a.gbpFetch)
	mux.Handle("GET /", http.FileServerFS(web.FS))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log.Printf("ZatpatSite listening on http://localhost:%s (gbp provider: %s)", port, a.provider.Mode())
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

// pickProvider honours GBP_PROVIDER; only the mock ships in this build.
func pickProvider() gbp.Provider {
	if os.Getenv("GBP_PROVIDER") == "live" {
		log.Println("GBP_PROVIDER=live requested but no live client is built; using mock")
	}
	return gbp.Mock{}
}

// seed creates two demo sites on a fresh store so the gallery is never empty.
func (a *app) seed() {
	if a.store.Count() > 0 {
		return
	}
	demos := []struct {
		name, city, category, theme, lang string
	}{
		{"Noor Beauty Salon", "Jaipur", "salon", "ivory", "en"},
		{"Annapurna Bhojanalaya", "Indore", "restaurant", "bazaar", "hi"},
	}
	for _, d := range demos {
		p, err := a.provider.Fetch(d.name, d.city, d.category)
		if err != nil {
			continue
		}
		site := &model.Site{
			Name: p.Name, Category: p.Category, City: p.City,
			Address: p.Address, Phone: p.Phone, WhatsApp: p.WhatsApp,
			MapsURL: p.MapsURL, Hours: p.Hours, Services: p.Services,
			Rating: p.Rating, Reviews: p.Reviews,
			Theme: d.theme, Lang: d.lang,
		}
		site.Normalize()
		if _, err := a.store.Create(site); err != nil {
			log.Printf("seed %s: %v", d.name, err)
		}
	}
	log.Printf("seeded %d demo sites", a.store.Count())
}

// ---- helpers ----

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func decodeSite(r *http.Request) (*model.Site, error) {
	var s model.Site
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	if err := dec.Decode(&s); err != nil {
		return nil, err
	}
	s.Normalize()
	if s.Name == "" {
		return nil, errors.New("business name is required")
	}
	return &s, nil
}

// ---- handlers ----

func (a *app) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":       "ok",
		"gbp_provider": a.provider.Mode(),
		"sites":        a.store.Count(),
	})
}

func (a *app) meta(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"categories": model.Categories,
		"themes":     generator.ThemeList,
		"langs":      model.Langs,
	})
}

func (a *app) createSite(w http.ResponseWriter, r *http.Request) {
	s, err := decodeSite(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := a.store.Create(s)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (a *app) listSites(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"sites": a.store.List()})
}

func (a *app) getSite(w http.ResponseWriter, r *http.Request) {
	s, err := a.store.Get(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "site not found")
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (a *app) updateSite(w http.ResponseWriter, r *http.Request) {
	s, err := decodeSite(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	s.ID = r.PathValue("id")
	saved, err := a.store.Update(s)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeErr(w, http.StatusNotFound, "site not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, saved)
}

func (a *app) deleteSite(w http.ResponseWriter, r *http.Request) {
	if err := a.store.Delete(r.PathValue("id")); err != nil {
		writeErr(w, http.StatusNotFound, "site not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (a *app) previewSite(w http.ResponseWriter, r *http.Request) {
	s, err := a.store.Get(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "site not found")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(generator.Generate(*s, time.Now())))
}

func (a *app) downloadSite(w http.ResponseWriter, r *http.Request) {
	s, err := a.store.Get(r.PathValue("id"))
	if err != nil {
		writeErr(w, http.StatusNotFound, "site not found")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="index.html"`)
	_, _ = w.Write([]byte(generator.Generate(*s, time.Now())))
}

// render generates HTML for an unsaved payload — powers the live preview.
func (a *app) render(w http.ResponseWriter, r *http.Request) {
	s, err := decodeSite(r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(generator.Generate(*s, time.Now())))
}

func (a *app) gbpFetch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		City     string `json:"city"`
		Category string `json:"category"`
		Lang     string `json:"lang"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<16)).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Name == "" {
		writeErr(w, http.StatusBadRequest, "business name is required to fetch")
		return
	}
	if !model.ValidLang(req.Lang) {
		req.Lang = "en"
	}
	p, err := a.provider.Fetch(req.Name, req.City, req.Category)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"profile": p,
		"about":   generator.AutoAbout(p.Name, p.City, p.Category, req.Lang),
	})
}
