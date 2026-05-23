// Package rollup provides aggregation of per-container drift results into
// cycle-level summaries.
//
// Usage:
//
//	a := rollup.New()
//	summary := a.Compute(results)
//	fmt.Printf("drift rate: %.0f%%\n", summary.DriftRate*100)
//
// The Summary type carries total, drifted and clean counts, an overall
// drift rate (0.0–1.0), and a per-container breakdown listing which
// fields (image, env, labels) have drifted.
package rollup
