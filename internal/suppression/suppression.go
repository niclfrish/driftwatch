// Package suppression provides a mechanism to suppress drift alerts
// for specific containers or fields for a configured duration.
package suppression

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Rule defines a suppression rule for a container and optional field.
type Rule struct {
	Container string    `json:"container"`
	Field     string    `json:"field,omitempty"` // empty means suppress all fields
	Until     time.Time `json:"until"`
}

// Store holds active suppression rules.
type Store struct {
	mu    sync.RWMutex
	rules []Rule
	path  string
}

// New creates a new Store, loading persisted rules from path if it exists.
func New(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Add inserts a new suppression rule and persists the store.
func (s *Store) Add(r Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules = append(s.rules, r)
	return s.save()
}

// IsSuppressed reports whether the given container+field combination
// is currently suppressed.
func (s *Store) IsSuppressed(container, field string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	for _, r := range s.rules {
		if r.Container != container {
			continue
		}
		if r.Until.Before(now) {
			continue
		}
		if r.Field == "" || r.Field == field {
			return true
		}
	}
	return false
}

// Prune removes expired rules and persists the store.
func (s *Store) Prune() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	active := s.rules[:0]
	for _, r := range s.rules {
		if r.Until.After(now) {
			active = append(active, r)
		}
	}
	s.rules = active
	return s.save()
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.rules)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.rules, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
