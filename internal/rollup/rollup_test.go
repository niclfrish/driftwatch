package rollup_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/rollup"
)

func makeResult(name string, imageDrift bool, envKeys []string, labelKeys []string) drift.Result {
	env := make(map[string]drift.FieldDiff, len(envKeys))
	for _, k := range envKeys {
		env[k] = drift.FieldDiff{Expected: "a", Actual: "b"}
	}
	labels := make(map[string]drift.FieldDiff, len(labelKeys))
	for _, k := range labelKeys {
		labels[k] = drift.FieldDiff{Expected: "x", Actual: "y"}
	}
	return drift.Result{
		ContainerName: name,
		ImageDrift:    imageDrift,
		EnvDrift:      env,
		LabelDrift:    labels,
	}
}

func TestCompute_Empty(t *testing.T) {
	a := rollup.New()
	s := a.Compute(nil)
	if s.Total != 0 || s.Drifted != 0 || s.Clean != 0 {
		t.Fatalf("expected zero summary, got %+v", s)
	}
	if s.DriftRate != 0.0 {
		t.Fatalf("expected 0.0 drift rate, got %f", s.DriftRate)
	}
}

func TestCompute_AllClean(t *testing.T) {
	a := rollup.New()
	results := []drift.Result{
		makeResult("web", false, nil, nil),
		makeResult("api", false, nil, nil),
	}
	s := a.Compute(results)
	if s.Total != 2 || s.Clean != 2 || s.Drifted != 0 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	if s.DriftRate != 0.0 {
		t.Fatalf("expected 0.0 drift rate")
	}
}

func TestCompute_AllDrifted(t *testing.T) {
	a := rollup.New()
	results := []drift.Result{
		makeResult("web", true, nil, nil),
		makeResult("api", false, []string{"PORT"}, nil),
	}
	s := a.Compute(results)
	if s.Total != 2 || s.Drifted != 2 || s.Clean != 0 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	if s.DriftRate != 1.0 {
		t.Fatalf("expected 1.0 drift rate, got %f", s.DriftRate)
	}
}

func TestCompute_DriftRate_Partial(t *testing.T) {
	a := rollup.New()
	results := []drift.Result{
		makeResult("a", true, nil, nil),
		makeResult("b", false, nil, nil),
		makeResult("c", false, nil, nil),
		makeResult("d", false, nil, nil),
	}
	s := a.Compute(results)
	if s.DriftRate != 0.25 {
		t.Fatalf("expected 0.25, got %f", s.DriftRate)
	}
}

func TestCompute_FieldNames(t *testing.T) {
	a := rollup.New()
	results := []drift.Result{
		makeResult("svc", true, []string{"DB_URL"}, []string{"team"}),
	}
	s := a.Compute(results)
	cs := s.ByContainer["svc"]
	if !cs.Drifted {
		t.Fatal("expected container to be marked drifted")
	}
	want := map[string]bool{"image": true, "env": true, "labels": true}
	for _, f := range cs.Fields {
		if !want[f] {
			t.Fatalf("unexpected field %q", f)
		}
	}
	if len(cs.Fields) != 3 {
		t.Fatalf("expected 3 drifted fields, got %d", len(cs.Fields))
	}
}
