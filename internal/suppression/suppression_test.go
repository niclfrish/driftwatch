package suppression_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/suppression"
)

func tempStore(t *testing.T) (*suppression.Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "suppressions.json")
	s, err := suppression.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s, path
}

func TestIsSuppressed_NoRules(t *testing.T) {
	s, _ := tempStore(t)
	if s.IsSuppressed("web", "image") {
		t.Error("expected not suppressed with no rules")
	}
}

func TestIsSuppressed_ActiveRule(t *testing.T) {
	s, _ := tempStore(t)
	err := s.Add(suppression.Rule{
		Container: "web",
		Field:     "image",
		Until:     time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if !s.IsSuppressed("web", "image") {
		t.Error("expected suppressed")
	}
}

func TestIsSuppressed_ExpiredRule(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.Add(suppression.Rule{
		Container: "web",
		Field:     "image",
		Until:     time.Now().Add(-time.Minute),
	})
	if s.IsSuppressed("web", "image") {
		t.Error("expected not suppressed for expired rule")
	}
}

func TestIsSuppressed_WildcardField(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.Add(suppression.Rule{
		Container: "db",
		Until:     time.Now().Add(time.Hour),
	})
	if !s.IsSuppressed("db", "env") {
		t.Error("expected wildcard suppression to match env")
	}
	if !s.IsSuppressed("db", "image") {
		t.Error("expected wildcard suppression to match image")
	}
}

func TestPrune_RemovesExpired(t *testing.T) {
	s, path := tempStore(t)
	_ = s.Add(suppression.Rule{Container: "a", Until: time.Now().Add(-time.Minute)})
	_ = s.Add(suppression.Rule{Container: "b", Until: time.Now().Add(time.Hour)})
	if err := s.Prune(); err != nil {
		t.Fatalf("Prune: %v", err)
	}
	s2, err := suppression.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if s2.IsSuppressed("a", "") {
		t.Error("expected 'a' pruned")
	}
	if !s2.IsSuppressed("b", "") {
		t.Error("expected 'b' still active")
	}
}

func TestNew_MissingFile_OK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	_, err := suppression.New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestNew_CorruptFile_Error(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := suppression.New(path)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
