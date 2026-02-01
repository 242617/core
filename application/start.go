package application

import (
	"context"
	"fmt"
	"time"
)

// start initializes all components.
func (a *Application) start(ctx context.Context) error {
	a.log.Info(ctx, "starting application",
		"name", a.name,
		"hostname", a.hostname,
	)

	var startedComponents []string

	rollback := func(startErr *ComponentError) error {
		a.log.Warn(ctx, "startup failed, rolling back started components")

		for i := len(startedComponents) - 1; i >= 0; i-- {
			c := a.components.ByName(startedComponents[i])
			if c != nil {
				stopCtx, stopCancel := context.WithTimeout(ctx, 5*time.Second)
				if stopErr := c.Stop(stopCtx); stopErr != nil {
					a.log.Error(ctx, "error during rollback stop", stopErr,
						"component", c.String(),
					)
				}
				stopCancel()
			}
		}

		return fmt.Errorf("start failed: %w", startErr)
	}

	for i := range a.components {
		c := a.components[i]

		a.log.Debug(ctx, "starting component", "component", c.String())

		startTime := time.Now()
		err := c.Start(ctx)
		duration := time.Since(startTime)

		if err != nil {
			a.log.Error(ctx, "cannot start component", err,
				"component", c.String(),
				"duration", duration,
			)
			return rollback(&ComponentError{
				Component: c.String(),
				Phase:     ComponentPhaseStart,
				Err:       err,
			})
		}

		a.log.Debug(ctx, "component started",
			"component", c.String(),
			"duration", duration,
		)

		startedComponents = append(startedComponents, c.String())
	}

	a.log.Info(ctx, "application started", "components", len(a.components))

	return nil
}
