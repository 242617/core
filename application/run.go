package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

// Run starts components, waits for shutdown signal, stops components.
func (a *Application) Run(ctx context.Context) error {
	a.startedMu.Lock()
	if a.started {
		a.startedMu.Unlock()
		return ErrApplicationAlreadyStarted
	}
	a.started = true
	a.startedMu.Unlock()

	shutdownCh := a.setupSignalHandling()

	startCtx, startCancel := context.WithTimeout(ctx, a.startTimeout)
	defer startCancel()

	if err := a.start(startCtx); err != nil {
		return fmt.Errorf("start application: %w", err)
	}

	select {
	case <-shutdownCh:
	case <-ctx.Done():
		a.log.Info(ctx, "parent context canceled, initiating shutdown")
	}

	stopCtx, stopCancel := context.WithTimeout(ctx, a.stopTimeout)
	defer stopCancel()

	if err := a.stop(stopCtx); err != nil {
		return fmt.Errorf("stop application: %w", err)
	}

	return nil
}

// setupSignalHandling returns channel closed on shutdown signal.
func (a *Application) setupSignalHandling() <-chan struct{} {
	shutdownCh := make(chan struct{})

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, os.Kill)
		defer signal.Stop(sigCh)

		select {
		case sig := <-sigCh:
			a.log.Info(context.Background(), "received shutdown signal", "signal", sig.String())
			close(shutdownCh)
		case <-a.stopCh:
			close(shutdownCh)
		case <-a.ctx.Done():
			close(shutdownCh)
		}
	}()

	return shutdownCh
}
