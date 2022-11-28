package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func (a *Application) Run() error {
	startCtx, startCancel := context.WithTimeout(context.Background(), a.startTimeout)
	defer startCancel()

	if err := a.start(startCtx); err != nil {
		return errors.Wrap(err, "start application")
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	stopCtx, stopCancel := context.WithTimeout(context.Background(), a.stopTimeout)
	defer stopCancel()

	if err := a.stop(stopCtx); err != nil {
		return errors.Wrap(err, "stop application")
	}

	return nil
}
