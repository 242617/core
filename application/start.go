package application

import (
	"context"

	"github.com/pkg/errors"
)

func (a *Application) start(ctx context.Context) error {
	a.log.Info().Msgf("starting %s (%s)", Name, Hostname)

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for i := 0; i < len(a.components); i++ {
			c := a.components[i]
			a.log.Info().Msgf("starting %q...", c)
			if err := c.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot start %q", c)
				errCh <- errors.Wrapf(err, "cannot start %q", c)
				return
			}
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return errors.New("start timeout")
	case err := <-errCh:
		return err
	case <-okCh:
	}

	a.log.Info().Msg("application started")
	return nil
}
