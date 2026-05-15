package drift

import (
	"testing"

	"github.com/user/driftwatch/internal/manifest"
)

func TestCompareEnv_NoDrift(t *testing.T) {
	expected := map[string]string{"PORT": "8080", "ENV": "prod"}
	actual := map[string]string{"PORT": "8080", "ENV": "prod", "EXTRA": "ignored"}

	reasons := compareEnv(expected, actual)
	if len(reasons) != 0 {
		t.Errorf("expected no drift reasons, got: %v", reasons)
	}
}

func TestCompareEnv_ValueMismatch(t *testing.T) {
	expected := map[string]string{"PORT": "8080"}
	actual := map[string]string{"PORT": "9090"}

	reasons := compareEnv(expected, actual)
	if len(reasons) != 1 {
		t.Fatalf("expected 1 reason, got %d: %v", len(reasons), reasons)
	}
}

func TestCompareEnv_MissingKey(t *testing.T) {
	expected := map[string]string{"SECRET": "abc"}
	actual := map[string]string{}

	reasons := compareEnv(expected, actual)
	if len(reasons) != 1 {
		t.Fatalf("expected 1 reason, got %d: %v", len(reasons), reasons)
	}
}

func TestCompareEnv_Empty(t *testing.T) {
	reasons := compareEnv(map[string]string{}, map[string]string{})
	if len(reasons) != 0 {
		t.Errorf("expected no drift reasons for empty maps, got: %v", reasons)
	}
}

func TestDriftResult_NoDrift(t *testing.T) {
	result := DriftResult{
		ContainerName: "web",
		Drifted:       false,
		Reasons:       nil,
	}

	if result.Drifted {
		t.Error("expected Drifted to be false")
	}
	if len(result.Reasons) != 0 {
		t.Errorf("expected no reasons, got: %v", result.Reasons)
	}
}

func TestNewDetector_NotNil(t *testing.T) {
	d := NewDetector(nil)
	if d == nil {
		t.Fatal("expected non-nil Detector")
	}
}

func TestManifestContainerFields(t *testing.T) {
	// Ensure manifest.Container has the fields Detector depends on.
	c := manifest.Container{
		Name:  "api",
		Image: "myapp:latest",
		Env:   map[string]string{"LOG_LEVEL": "info"},
	}

	if c.Name != "api" {
		t.Errorf("unexpected Name: %q", c.Name)
	}
	if c.Image != "myapp:latest" {
		t.Errorf("unexpected Image: %q", c.Image)
	}
}
