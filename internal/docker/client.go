package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ContainerInfo holds the relevant runtime details of a running container.
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Env     map[string]string
	Labels  map[string]string
}

// Client wraps the Docker API client.
type Client struct {
	dc *client.Client
}

// NewClient creates a new Docker client using environment-based configuration.
func NewClient() (*Client, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker: failed to create client: %w", err)
	}
	return &Client{dc: dc}, nil
}

// Close releases resources held by the underlying Docker client.
func (c *Client) Close() error {
	return c.dc.Close()
}

// ListRunning returns ContainerInfo for all currently running containers.
// An optional name filter can be provided; pass an empty string to list all.
func (c *Client) ListRunning(ctx context.Context, nameFilter string) ([]ContainerInfo, error) {
	f := filters.NewArgs(filters.Arg("status", "running"))
	if nameFilter != "" {
		f.Add("name", nameFilter)
	}

	containers, err := c.dc.ContainerList(ctx, types.ContainerListOptions{Filters: f})
	if err != nil {
		return nil, fmt.Errorf("docker: list containers: %w", err)
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, ct := range containers {
		info, err := c.inspect(ctx, ct.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}

// inspect fetches detailed info for a single container.
func (c *Client) inspect(ctx context.Context, id string) (ContainerInfo, error) {
	data, err := c.dc.ContainerInspect(ctx, id)
	if err != nil {
		return ContainerInfo{}, fmt.Errorf("docker: inspect %s: %w", id, err)
	}

	name := data.Name
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}

	env := parseEnv(data.Config.Env)

	return ContainerInfo{
		ID:     data.ID,
		Name:   name,
		Image:  data.Config.Image,
		Env:    env,
		Labels: data.Config.Labels,
	}, nil
}

// parseEnv converts a slice of "KEY=VALUE" strings into a map.
func parseEnv(raw []string) map[string]string {
	m := make(map[string]string, len(raw))
	for _, kv := range raw {
		for i := 0; i < len(kv); i++ {
			if kv[i] == '=' {
				m[kv[:i]] = kv[i+1:]
				break
			}
		}
	}
	return m
}
