// Package baseline manages the expected (desired) state of containers
// derived from manifests, used as the reference point for drift detection.
package baseline

import (
	"fmt"
	"sync"

	"github.com/yourorg/driftwatch/internal/manifest"
)

// Entry holds the desired state for a single container.
type Entry struct {
	Name   string
	Image  string
	Env    map[string]string
	Labels map[string]string
}

// Store holds a thread-safe map of container name -> baseline Entry.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New creates an empty baseline Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
	}
}

// Load populates the store from a parsed manifest.
func (s *Store) Load(m *manifest.Manifest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := make(map[string]Entry, len(m.Containers))
	for _, c := range m.Containers {
		if c.Name == "" {
			return fmt.Errorf("baseline: container entry missing name")
		}
		env := make(map[string]string, len(c.Env))
		for k, v := range c.Env {
			env[k] = v
		}
		labels := make(map[string]string, len(c.Labels))
		for k, v := range c.Labels {
			labels[k] = v
		}
		next[c.Name] = Entry{
			Name:   c.Name,
			Image:  c.Image,
			Env:    env,
			Labels: labels,
		}
	}
	s.entries = next
	return nil
}

// Get returns the baseline Entry for the given container name.
func (s *Store) Get(name string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[name]
	return e, ok
}

// All returns a copy of all baseline entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of entries in the store.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
