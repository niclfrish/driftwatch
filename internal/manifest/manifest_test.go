package manifest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/manifest"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "manifest-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	raw := `
version: "1"
containers:
  - name: api
    image: myrepo/api:v1.2.3
    restartPolicy: always
    env:
      LOG_LEVEL: debug
    ports:
      - "8080:8080"
`
	path := writeTemp(t, raw)
	m, err := manifest.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m.Containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(m.Containers))
	}
	c := m.Containers[0]
	if c.Name != "api" {
		t.Errorf("expected name 'api', got %q", c.Name)
	}
	if c.Image != "myrepo/api:v1.2.3" {
		t.Errorf("unexpected image: %q", c.Image)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := manifest.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoContainers(t *testing.T) {
	path := writeTemp(t, "version: \"1\"\ncontainers: []\n")
	_, err := manifest.Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty containers")
	}
}

func TestLoad_MissingImage(t *testing.T) {
	raw := "version: \"1\"\ncontainers:\n  - name: api\n"
	path := writeTemp(t, raw)
	_, err := manifest.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing image")
	}
}

func TestByName(t *testing.T) {
	raw := `version: "1"
containers:
  - name: worker
    image: myrepo/worker:latest
`
	path := writeTemp(t, raw)
	m, _ := manifest.Load(path)

	c, err := m.ByName("worker")
	if err != nil || c.Image != "myrepo/worker:latest" {
		t.Errorf("ByName failed: err=%v, spec=%v", err, c)
	}

	_, err = m.ByName("missing")
	if err == nil {
		t.Error("expected error for unknown container name")
	}
}
