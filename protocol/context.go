package protocol

import (
	"context"
	"time"
)

// ContextFunc is a function type for callbacks that receive a context.
// It's commonly used for lifecycle hooks like PreStart, PostStart, PreStop, PostStop.
type ContextFunc func(context.Context) error

func ContextRunWithTimeout(ctx context.Context, d time.Duration, f ContextFunc) error {
	errCh := make(chan error)
	go func() {
		if err := f(ctx); err != nil {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		return context.DeadlineExceeded
	case err := <-errCh:
		return err
	case <-time.After(d):
		return nil
	}
}
