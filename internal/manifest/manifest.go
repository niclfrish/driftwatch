// Package manifest loads and validates driftwatch manifest files.
package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Container describes the desired state of a single container.
type Container struct {
	Name   string            `yaml:"name"`
	Image  string            `yaml:"image"`
	Env    map[string]string `yaml:"env"`
	Labels map[string]string `yaml:"labels"`
}

// Manifest is the top-level structure of a driftwatch manifest file.
type Manifest struct {
	Version    string      `yaml:"version"`
	Containers []Container `yaml:"containers"`
}

// Load reads and parses a manifest YAML file from the given path.
func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("manifest: read %q: %w", path, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("manifest: parse %q: %w", path, err)
	}

	if len(m.Containers) == 0 {
		return nil, fmt.Errorf("manifest: %q defines no containers", path)
	}

	for i, c := range m.Containers {
		if c.Name == "" {
			return nil, fmt.Errorf("manifest: container[%d] missing name", i)
		}
		if c.Image == "" {
			return nil, fmt.Errorf("manifest: container %q missing image", c.Name)
		}
	}

	return &m, nil
}
