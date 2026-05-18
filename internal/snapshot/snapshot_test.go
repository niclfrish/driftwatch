package snapshot_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/snapshot"
)

func makeRecord(name string, drifted bool) snapshot.Record {
	return snapshot.Record{
		ContainerName: name,
		Timestamp:     time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Result: drift.Result{
			ContainerName: name,
			Drifted:       drifted,
			Diffs:         []drift.Diff{{Field: "env.PORT", Expected: "8080", Actual: "9090"}},
		},
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	store, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rec := makeRecord("web", true)
	if err := store.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load("web")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got.ContainerName != rec.ContainerName {
		t.Errorf("ContainerName: got %q, want %q", got.ContainerName, rec.ContainerName)
	}
	if got.Result.Drifted != rec.Result.Drifted {
		t.Errorf("Drifted: got %v, want %v", got.Result.Drifted, rec.Result.Drifted)
	}
	if len(got.Result.Diffs) != 1 {
		t.Errorf("Diffs len: got %d, want 1", len(got.Result.Diffs))
	}
}

func TestLoad_NotExist(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)

	_, err := store.Load("nonexistent")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected ErrNotExist, got %v", err)
	}
}

func TestSave_OverwritesPrevious(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)

	first := makeRecord("api", false)
	_ = store.Save(first)

	second := makeRecord("api", true)
	_ = store.Save(second)

	got, err := store.Load("api")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !got.Result.Drifted {
		t.Error("expected Drifted=true after overwrite")
	}
}

func TestNew_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/nested/snapshots"
	_, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
