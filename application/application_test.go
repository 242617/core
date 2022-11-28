package application_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/242617/core/application"
)

func TestBasic(t *testing.T) {
	period := 10 * time.Millisecond
	a, err := application.New()
	assert.NoError(t, err, "new application")
	go func() {
		time.Sleep(period)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	assert.NoError(t, a.Run(), "run application")
}

func TestWithComponent(t *testing.T) {
	period := 10 * time.Millisecond
	a, err := application.New(
		application.WithComponents(
			application.NewMethodsComponent("test",
				func(context.Context) error { return nil },
				func(context.Context) error { return nil },
			),
		),
	)
	assert.NoError(t, err, "new application")
	go func() {
		time.Sleep(period)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	assert.NoError(t, a.Run(), "run application")
}

func TestStartError(t *testing.T) {
	startErr := errors.New("start error")
	cmp := application.NewMethodsComponent("test",
		func(context.Context) error { return startErr },
		func(context.Context) error { return nil },
	)

	a, err := application.New(application.WithComponents(cmp))
	assert.NoError(t, err, "new application")
	assert.ErrorIs(t, a.Run(), startErr, "start error")
}

func TestStopError(t *testing.T) {
	period := 100 * time.Millisecond
	stopErr := errors.New("stop error")
	cmp := application.NewMethodsComponent("test",
		func(context.Context) error { return nil },
		func(context.Context) error { return stopErr },
	)

	a, err := application.New(application.WithComponents(cmp))
	assert.NoError(t, err, "new application")
	go func() {
		time.Sleep(period)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	assert.ErrorIs(t, a.Run(), stopErr, "stop error")
}
