// Package healthcheck provides a lightweight HTTP liveness/readiness probe
// for the driftwatch daemon.
//
// Usage:
//
//	check := healthcheck.New()
//
//	// call after every scheduled drift scan:
//	check.RecordRun(err)
//
//	// mount the handler:
//	http.Handle("/healthz", check.Handler())
//
// The endpoint returns HTTP 200 when the last run succeeded and
// HTTP 503 when it failed, with a JSON body containing run statistics.
package healthcheck
