package manifest

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ContainerSpec represents the desired state of a container as defined in a manifest file.
type ContainerSpec struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Env         map[string]string `yaml:"env"`
	Ports       []string          `yaml:"ports"`
	Labels      map[string]string `yaml:"labels"`
	RestartPolicy string          `yaml:"restartPolicy"`
}

// Manifest holds one or more container specs loaded from a YAML file.
type Manifest struct {
	Version    string          `yaml:"version"`
	Containers []ContainerSpec `yaml:"containers"`
}

// Load reads and parses a YAML manifest file from the given path.
func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest %q: %w", path, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest %q: %w", path, err)
	}

	if err := m.validate(); err != nil {
		return nil, fmt.Errorf("invalid manifest %q: %w", path, err)
	}

	return &m, nil
}

// validate performs basic sanity checks on the manifest.
func (m *Manifest) validate() error {
	if len(m.Containers) == 0 {
		return fmt.Errorf("manifest must define at least one container")
	}
	for i, c := range m.Containers {
		if c.Name == "" {
			return fmt.Errorf("container[%d] missing name", i)
		}
		if c.Image == "" {
			return fmt.Errorf("container %q missing image", c.Name)
		}
	}
	return nil
}

// ByName returns the ContainerSpec with the given name, or an error if not found.
func (m *Manifest) ByName(name string) (*ContainerSpec, error) {
	for i := range m.Containers {
		if m.Containers[i].Name == name {
			return &m.Containers[i], nil
		}
	}
	return nil, fmt.Errorf("container %q not found in manifest", name)
}
