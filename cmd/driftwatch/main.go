// Package main is the entry point for the driftwatch daemon.
// It wires together manifest loading, Docker inspection, drift detection,
// and reporting into a single runnable command.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yourorg/driftwatch/internal/docker"
	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/manifest"
	"github.com/yourorg/driftwatch/internal/reporter"
)

func main() {
	var (
		manifestPath = flag.String("manifest", "manifest.yaml", "Path to the container manifest file")
		outputFormat = flag.String("format", "text", "Output format: text or json")
		watchMode    = flag.Bool("watch", false, "Run continuously and re-check on each interval")
		interval     = flag.Duration("interval", 30*time.Second, "Polling interval when running in watch mode")
		verbose      = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	// Validate format flag early so we fail fast.
	if *outputFormat != "text" && *outputFormat != "json" {
		fmt.Fprintf(os.Stderr, "error: unsupported format %q — must be \"text\" or \"json\"\n", *outputFormat)
		os.Exit(1)
	}

	ctx := context.Background()

	// Initialise Docker client (honours DOCKER_HOST / DOCKER_TLS_VERIFY etc.).
	dockerClient, err := docker.NewClient()
	if err != nil {
		log.Fatalf("failed to create Docker client: %v", err)
	}
	defer dockerClient.Close()

	rep, err := reporter.New(*outputFormat, os.Stdout)
	if err != nil {
		log.Fatalf("failed to create reporter: %v", err)
	}

	if *watchMode {
		log.Printf("watch mode enabled — checking every %s", *interval)
		for {
			if err := run(ctx, *manifestPath, dockerClient, rep, *verbose); err != nil {
				log.Printf("check error: %v", err)
			}
			time.Sleep(*interval)
		}
	}

	// Single-shot mode — exit with a non-zero code when drift is found so the
	// binary can be used directly in CI pipelines.
	if err := run(ctx, *manifestPath, dockerClient, rep, *verbose); err != nil {
		log.Fatalf("error: %v", err)
	}
}

// run executes one full drift-detection cycle and writes the report.
func run(ctx context.Context, manifestPath string, client *docker.Client, rep *reporter.Reporter, verbose bool) error {
	manifests, err := manifest.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	if verbose {
		log.Printf("loaded %d container spec(s) from %s", len(manifests.Containers), manifestPath)
	}

	detector := drift.NewDetector(client)

	results, err := detector.Detect(ctx, manifests)
	if err != nil {
		return fmt.Errorf("detecting drift: %w", err)
	}

	if err := rep.Report(results); err != nil {
		return fmt.Errorf("writing report: %w", err)
	}

	// Signal drift to the caller via exit code in single-shot mode.
	for _, r := range results {
		if r.HasDrift() {
			os.Exit(2)
		}
	}

	return nil
}
