// Package circuitbreaker implements a per-container circuit breaker for
// driftwatch. When a container inspection repeatedly fails (e.g. the Docker
// daemon is unreachable or the container has been removed), the circuit opens
// and further checks are skipped until a configurable cool-down period has
// elapsed. After the cool-down a single probe (half-open) is permitted; a
// successful probe closes the circuit while a failed probe re-opens it.
//
// Usage:
//
//	br := circuitbreaker.New(5, 30*time.Second)
//	if !br.Allow(containerID) {
//		// skip this container
//		return
//	}
//	if err := inspect(containerID); err != nil {
//		br.RecordFailure(containerID)
//	} else {
//		br.RecordSuccess(containerID)
//	}
package circuitbreaker
