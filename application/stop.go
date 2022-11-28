package application

import (
	"context"

	"github.com/pkg/errors"
)

func (a *Application) stop(ctx context.Context) error {
	a.log.Info().Msgf("stopping %s", Name)

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for i := len(a.components) - 1; i >= 0; i-- {
			c := a.components[i]
			a.log.Info().Msgf("stopping %q...", c)
			if err := c.Stop(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot stop %q", c)
				errCh <- errors.Wrapf(err, "cannot stop %q", c)
				return
			}
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("stop timeout")
	case err := <-errCh:
		return err
	case <-okCh:
	}

	a.log.Info().Msg("application stopped")
	return nil
}
