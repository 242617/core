package application

import (
	"context"
	"errors"
	"time"
)

// stop shuts down components in reverse order.
func (a *Application) stop(ctx context.Context) error {
	a.log.Info(ctx, "stopping application")

	var errs []error

	for i := len(a.components) - 1; i >= 0; i-- {
		c := a.components[i]

		a.log.Debug(ctx, "stopping component", "component", c.String())

		startTime := time.Now()
		err := c.Stop(ctx)
		duration := time.Since(startTime)

		if err != nil {
			a.log.Error(ctx, "error stopping component", err,
				"component", c.String(),
				"duration", duration,
			)
			errs = append(errs, &ComponentError{
				Component: c.String(),
				Phase:     ComponentPhaseStop,
				Err:       err,
			})
		} else {
			a.log.Debug(ctx, "component stopped",
				"component", c.String(),
				"duration", duration,
			)
		}
	}

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		a.log.Warn(ctx, "shutdown timeout, some goroutines may still be running")
	}

	a.log.Info(ctx, "application stopped")

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
