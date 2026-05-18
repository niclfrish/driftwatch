package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
)

// Record holds a timestamped drift result for a single container.
type Record struct {
	ContainerName string            `json:"container_name"`
	Timestamp     time.Time         `json:"timestamp"`
	Result        drift.Result      `json:"result"`
}

// Store persists and retrieves drift snapshots on disk.
type Store struct {
	dir string
}

// New creates a Store that writes snapshot files under dir.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir %q: %w", dir, err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a Record for the given container to disk, overwriting any
// previous snapshot for that container.
func (s *Store) Save(r Record) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := s.filePath(r.ContainerName)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %q: %w", path, err)
	}
	return nil
}

// Load reads the most recent snapshot for the given container name.
// Returns os.ErrNotExist if no snapshot has been saved yet.
func (s *Store) Load(containerName string) (Record, error) {
	var rec Record
	path := s.filePath(containerName)
	data, err := os.ReadFile(path)
	if err != nil {
		return rec, err
	}
	if err := json.Unmarshal(data, &rec); err != nil {
		return rec, fmt.Errorf("snapshot: unmarshal %q: %w", path, err)
	}
	return rec, nil
}

func (s *Store) filePath(containerName string) string {
	return filepath.Join(s.dir, containerName+".json")
}
