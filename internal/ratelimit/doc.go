// Package ratelimit provides per-container rate limiting for drift notifications.
//
// It prevents alert storms by enforcing a minimum interval between successive
// notifications for the same container. Each container is tracked independently,
// so a noisy container does not affect the notification cadence of others.
//
// # Usage
//
//	rl := ratelimit.New(5 * time.Minute)
//
//	if rl.Allow(containerName) {
//		// send notification
//	}
//
// The window duration is configured once at construction time and applies
// uniformly to all containers. Call Reset to clear the state for a specific
// container, for example after a suppression rule is lifted.
package ratelimit
