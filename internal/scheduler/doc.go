// Package scheduler provides a simple interval-based job scheduler for
// driftwatch. It runs a user-supplied Job function immediately on start and
// then repeatedly at the configured interval until the context is cancelled.
//
// Example usage:
//
//	s := scheduler.New(30*time.Second, func(ctx context.Context) error {
//		// perform drift detection
//		return nil
//	}, logger)
//	if err := s.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
package scheduler
