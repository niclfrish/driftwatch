// Package history records drift detection results over time,
// allowing trend analysis and change detection between runs.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/your-org/driftwatch/internal/drift"
)

// Entry represents a single recorded drift check.
type Entry struct {
	Timestamp time.Time          `json:"timestamp"`
	Results   []drift.Result     `json:"results"`
}

// Store persists drift history to disk.
type Store struct {
	dir string
}

// New creates a Store that writes history files under dir.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Record saves a new entry with the current timestamp.
func (s *Store) Record(results []drift.Result) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Results:   results,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("history: marshal: %w", err)
	}
	name := entry.Timestamp.Format("20060102T150405Z") + ".json"
	path := filepath.Join(s.dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("history: write %s: %w", path, err)
	}
	return nil
}

// List returns all recorded entries sorted by timestamp ascending.
func (s *Store) List() ([]Entry, error) {
	matches, err := filepath.Glob(filepath.Join(s.dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("history: glob: %w", err)
	}
	var entries []Entry
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("history: read %s: %w", path, err)
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, fmt.Errorf("history: unmarshal %s: %w", path, err)
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})
	return entries, nil
}

// Latest returns the most recently recorded entry, or nil if none exist.
func (s *Store) Latest() (*Entry, error) {
	entries, err := s.List()
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	e := entries[len(entries)-1]
	return &e, nil
}
