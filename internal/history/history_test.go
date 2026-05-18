package history_test

import (
	"os"
	"testing"
	"time"

	"github.com/your-org/driftwatch/internal/drift"
	"github.com/your-org/driftwatch/internal/history"
)

func makeResults(drifted bool) []drift.Result {
	return []drift.Result{
		{
			ContainerName: "web",
			Drifted:       drifted,
			Diffs:         nil,
		},
	}
}

func TestNew_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/history"
	_, err := history.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
}

func TestRecord_And_List_RoundTrip(t *testing.T) {
	s, err := history.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	results := makeResults(true)
	if err := s.Record(results); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(entries[0].Results))
	}
	if !entries[0].Results[0].Drifted {
		t.Error("expected Drifted=true")
	}
}

func TestList_SortedAscending(t *testing.T) {
	s, err := history.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 3; i++ {
		if err := s.Record(makeResults(i%2 == 0)); err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Millisecond)
	}
	entries, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for i := 1; i < len(entries); i++ {
		if entries[i].Timestamp.Before(entries[i-1].Timestamp) {
			t.Error("entries not sorted ascending")
		}
	}
}

func TestLatest_Empty(t *testing.T) {
	s, err := history.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	e, err := s.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if e != nil {
		t.Error("expected nil for empty store")
	}
}

func TestLatest_ReturnsNewest(t *testing.T) {
	s, err := history.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Record(makeResults(false)); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Millisecond)
	if err := s.Record(makeResults(true)); err != nil {
		t.Fatal(err)
	}
	e, err := s.Latest()
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil entry")
	}
	if !e.Results[0].Drifted {
		t.Error("expected latest entry to have Drifted=true")
	}
}
