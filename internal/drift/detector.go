package drift

import (
	"fmt"

	"github.com/user/driftwatch/internal/docker"
	"github.com/user/driftwatch/internal/manifest"
)

// DriftResult holds the result of comparing a container against its manifest spec.
type DriftResult struct {
	ContainerName string
	Drifted       bool
	Reasons       []string
}

// Detector compares running containers against manifest definitions.
type Detector struct {
	client *docker.Client
}

// NewDetector creates a new Detector using the provided Docker client.
func NewDetector(client *docker.Client) *Detector {
	return &Detector{client: client}
}

// Detect checks all containers defined in the manifest for drift.
func (d *Detector) Detect(m *manifest.Manifest) ([]DriftResult, error) {
	results := make([]DriftResult, 0, len(m.Containers))

	for _, spec := range m.Containers {
		info, err := d.client.InspectContainer(spec.Name)
		if err != nil {
			return nil, fmt.Errorf("inspect container %q: %w", spec.Name, err)
		}

		result := DriftResult{ContainerName: spec.Name}

		// Check image drift
		if info.Image != spec.Image {
			result.Drifted = true
			result.Reasons = append(result.Reasons,
				fmt.Sprintf("image mismatch: running=%q manifest=%q", info.Image, spec.Image),
			)
		}

		// Check environment variable drift
		envDrifts := compareEnv(spec.Env, info.Env)
		if len(envDrifts) > 0 {
			result.Drifted = true
			result.Reasons = append(result.Reasons, envDrifts...)
		}

		results = append(results, result)
	}

	return results, nil
}

// compareEnv returns a list of drift reasons between expected and actual env maps.
func compareEnv(expected, actual map[string]string) []string {
	var reasons []string

	for k, want := range expected {
		got, ok := actual[k]
		if !ok {
			reasons = append(reasons, fmt.Sprintf("env %q missing from container", k))
			continue
		}
		if got != want {
			reasons = append(reasons,
				fmt.Sprintf("env %q mismatch: running=%q manifest=%q", k, got, want),
			)
		}
	}

	return reasons
}
