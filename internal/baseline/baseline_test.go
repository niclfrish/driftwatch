package baseline_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/baseline"
	"github.com/yourorg/driftwatch/internal/manifest"
)

func makeManifest(containers []manifest.Container) *manifest.Manifest {
	return &manifest.Manifest{Containers: containers}
}

func TestLoad_PopulatesEntries(t *testing.T) {
	s := baseline.New()
	m := makeManifest([]manifest.Container{
		{Name: "web", Image: "nginx:1.25", Env: map[string]string{"PORT": "8080"}},
		{Name: "db", Image: "postgres:15"},
	})
	if err := s.Load(m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", s.Len())
	}
}

func TestGet_KnownEntry(t *testing.T) {
	s := baseline.New()
	m := makeManifest([]manifest.Container{
		{Name: "web", Image: "nginx:1.25", Env: map[string]string{"PORT": "8080"}},
	})
	_ = s.Load(m)

	e, ok := s.Get("web")
	if !ok {
		t.Fatal("expected entry for 'web'")
	}
	if e.Image != "nginx:1.25" {
		t.Errorf("expected image nginx:1.25, got %s", e.Image)
	}
	if e.Env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", e.Env["PORT"])
	}
}

func TestGet_UnknownEntry(t *testing.T) {
	s := baseline.New()
	_ = s.Load(makeManifest(nil))
	_, ok := s.Get("missing")
	if ok {
		t.Error("expected no entry for unknown container")
	}
}

func TestLoad_MissingName_ReturnsError(t *testing.T) {
	s := baseline.New()
	m := makeManifest([]manifest.Container{
		{Name: "", Image: "nginx:1.25"},
	})
	if err := s.Load(m); err == nil {
		t.Error("expected error for missing container name")
	}
}

func TestLoad_ReplacesExistingEntries(t *testing.T) {
	s := baseline.New()
	_ = s.Load(makeManifest([]manifest.Container{
		{Name: "web", Image: "nginx:1.24"},
	}))
	_ = s.Load(makeManifest([]manifest.Container{
		{Name: "api", Image: "myapp:2.0"},
	}))
	if s.Len() != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", s.Len())
	}
	if _, ok := s.Get("web"); ok {
		t.Error("old entry 'web' should have been replaced")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := baseline.New()
	_ = s.Load(makeManifest([]manifest.Container{
		{Name: "a", Image: "img:1"},
		{Name: "b", Image: "img:2"},
	}))
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2, got %d", len(all))
	}
}
