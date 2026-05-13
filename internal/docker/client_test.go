package docker

import (
	"testing"
)

func TestParseEnv_Normal(t *testing.T) {
	input := []string{"FOO=bar", "BAZ=qux", "EMPTY="}
	got := parseEnv(input)

	cases := map[string]string{
		"FOO":   "bar",
		"BAZ":   "qux",
		"EMPTY": "",
	}
	for k, want := range cases {
		if v, ok := got[k]; !ok {
			t.Errorf("parseEnv: key %q missing", k)
		} else if v != want {
			t.Errorf("parseEnv: key %q = %q, want %q", k, v, want)
		}
	}
	if len(got) != len(cases) {
		t.Errorf("parseEnv: got %d entries, want %d", len(got), len(cases))
	}
}

func TestParseEnv_Empty(t *testing.T) {
	got := parseEnv(nil)
	if len(got) != 0 {
		t.Errorf("parseEnv(nil): expected empty map, got %v", got)
	}
}

func TestParseEnv_ValueContainsEquals(t *testing.T) {
	input := []string{"URL=http://example.com?a=1&b=2"}
	got := parseEnv(input)
	want := "http://example.com?a=1&b=2"
	if got["URL"] != want {
		t.Errorf("parseEnv value with '=': got %q, want %q", got["URL"], want)
	}
}

func TestParseEnv_NoEquals(t *testing.T) {
	// Entries without '=' should be silently skipped.
	input := []string{"NOEQUALS", "VALID=yes"}
	got := parseEnv(input)
	if _, ok := got["NOEQUALS"]; ok {
		t.Error("parseEnv: entry without '=' should not appear in result")
	}
	if got["VALID"] != "yes" {
		t.Errorf("parseEnv: VALID = %q, want \"yes\"", got["VALID"])
	}
}

func TestNewClient_EnvFallback(t *testing.T) {
	// NewClient should not panic even when DOCKER_HOST is unset;
	// it may return an error if Docker is unavailable, which is acceptable.
	_, err := NewClient()
	// We only assert no panic; a missing daemon is fine in CI.
	_ = err
}
