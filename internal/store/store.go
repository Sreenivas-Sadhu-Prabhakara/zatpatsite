// Package store is the in-memory site store with JSON snapshot persistence.
package store

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"zatpatsite/internal/model"
)

// ErrNotFound is returned when a site id does not exist.
var ErrNotFound = errors.New("site not found")

type snapshot struct {
	Sites []*model.Site `json:"sites"`
}

// Store keeps sites in memory behind a mutex and snapshots to a JSON file
// after every write.
type Store struct {
	mu    sync.Mutex
	path  string
	sites map[string]*model.Site
}

// Open loads the snapshot at path if present.
func Open(path string) (*Store, error) {
	s := &Store{path: path, sites: map[string]*model.Site{}}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	var snap snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return nil, fmt.Errorf("corrupt snapshot %s: %w", path, err)
	}
	for _, site := range snap.Sites {
		s.sites[site.ID] = site
	}
	return s, nil
}

// save writes the snapshot; caller must hold the lock.
func (s *Store) save() error {
	snap := snapshot{Sites: s.listLocked()}
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *Store) listLocked() []*model.Site {
	out := make([]*model.Site, 0, len(s.sites))
	for _, site := range s.sites {
		out = append(out, site)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].UpdatedAt.After(out[j].UpdatedAt)
	})
	return out
}

// clone deep-copies a site so callers never share slices with the store.
func clone(in *model.Site) *model.Site {
	c := *in
	c.Services = append([]model.Service(nil), in.Services...)
	c.Reviews = append([]model.Review(nil), in.Reviews...)
	return &c
}

// List returns all sites, most recently updated first.
func (s *Store) List() []*model.Site {
	s.mu.Lock()
	defer s.mu.Unlock()
	sites := s.listLocked()
	out := make([]*model.Site, len(sites))
	for i, site := range sites {
		out[i] = clone(site)
	}
	return out
}

// Get returns one site by id.
func (s *Store) Get(id string) (*model.Site, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	site, ok := s.sites[id]
	if !ok {
		return nil, ErrNotFound
	}
	return clone(site), nil
}

// Create assigns an id + timestamps and persists.
func (s *Store) Create(site *model.Site) (*model.Site, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	site.ID = newID()
	now := time.Now().UTC()
	site.CreatedAt, site.UpdatedAt = now, now
	s.sites[site.ID] = clone(site)
	return clone(site), s.save()
}

// Update replaces an existing site, keeping its CreatedAt.
func (s *Store) Update(site *model.Site) (*model.Site, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.sites[site.ID]
	if !ok {
		return nil, ErrNotFound
	}
	site.CreatedAt = old.CreatedAt
	site.UpdatedAt = time.Now().UTC()
	s.sites[site.ID] = clone(site)
	return clone(site), s.save()
}

// Delete removes a site.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sites[id]; !ok {
		return ErrNotFound
	}
	delete(s.sites, id)
	return s.save()
}

// Count reports the number of stored sites.
func (s *Store) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.sites)
}

func newID() string {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
