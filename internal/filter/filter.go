package filter

import "strings"

// Options holds filtering criteria for containers.
type Options struct {
	// Labels filters containers that have all specified label key=value pairs.
	Labels map[string]string
	// NamePrefix filters containers whose names start with the given prefix.
	NamePrefix string
}

// ContainerMeta is a minimal view of a container used for filtering.
type ContainerMeta struct {
	Name   string
	Labels map[string]string
}

// Match returns true when the container satisfies all criteria in opts.
// An empty Options matches every container.
func Match(c ContainerMeta, opts Options) bool {
	if opts.NamePrefix != "" && !strings.HasPrefix(c.Name, opts.NamePrefix) {
		return false
	}
	for k, v := range opts.Labels {
		got, ok := c.Labels[k]
		if !ok || got != v {
			return false
		}
	}
	return true
}

// Apply returns the subset of containers that satisfy opts.
func Apply(containers []ContainerMeta, opts Options) []ContainerMeta {
	if opts.NamePrefix == "" && len(opts.Labels) == 0 {
		return containers
	}
	out := make([]ContainerMeta, 0, len(containers))
	for _, c := range containers {
		if Match(c, opts) {
			out = append(out, c)
		}
	}
	return out
}
