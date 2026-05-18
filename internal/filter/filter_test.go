package filter_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/filter"
)

func makeContainer(name string, labels map[string]string) filter.ContainerMeta {
	return filter.ContainerMeta{Name: name, Labels: labels}
}

func TestMatch_NoFilter(t *testing.T) {
	c := makeContainer("web-1", map[string]string{"env": "prod"})
	if !filter.Match(c, filter.Options{}) {
		t.Fatal("empty options should match any container")
	}
}

func TestMatch_NamePrefix_Match(t *testing.T) {
	c := makeContainer("web-1", nil)
	if !filter.Match(c, filter.Options{NamePrefix: "web"}) {
		t.Fatal("expected match on name prefix")
	}
}

func TestMatch_NamePrefix_NoMatch(t *testing.T) {
	c := makeContainer("worker-1", nil)
	if filter.Match(c, filter.Options{NamePrefix: "web"}) {
		t.Fatal("expected no match on name prefix")
	}
}

func TestMatch_Label_Match(t *testing.T) {
	c := makeContainer("api", map[string]string{"env": "prod", "tier": "backend"})
	opts := filter.Options{Labels: map[string]string{"env": "prod"}}
	if !filter.Match(c, opts) {
		t.Fatal("expected label match")
	}
}

func TestMatch_Label_WrongValue(t *testing.T) {
	c := makeContainer("api", map[string]string{"env": "staging"})
	opts := filter.Options{Labels: map[string]string{"env": "prod"}}
	if filter.Match(c, opts) {
		t.Fatal("expected no match due to wrong label value")
	}
}

func TestMatch_Label_Missing(t *testing.T) {
	c := makeContainer("api", map[string]string{})
	opts := filter.Options{Labels: map[string]string{"env": "prod"}}
	if filter.Match(c, opts) {
		t.Fatal("expected no match due to missing label")
	}
}

func TestApply_FiltersCorrectly(t *testing.T) {
	containers := []filter.ContainerMeta{
		makeContainer("web-1", map[string]string{"env": "prod"}),
		makeContainer("worker-1", map[string]string{"env": "prod"}),
		makeContainer("web-2", map[string]string{"env": "staging"}),
	}
	opts := filter.Options{NamePrefix: "web", Labels: map[string]string{"env": "prod"}}
	result := filter.Apply(containers, opts)
	if len(result) != 1 || result[0].Name != "web-1" {
		t.Fatalf("expected [web-1], got %v", result)
	}
}

func TestApply_EmptyOptions_ReturnsAll(t *testing.T) {
	containers := []filter.ContainerMeta{
		makeContainer("a", nil),
		makeContainer("b", nil),
	}
	result := filter.Apply(containers, filter.Options{})
	if len(result) != 2 {
		t.Fatalf("expected all containers returned, got %d", len(result))
	}
}
