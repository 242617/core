package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/242617/core/mocks"
	"github.com/stretchr/testify/mock"
)

func TestRun_AlreadyStarted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Return(nil)
	ml.On("Stop", mock.Anything).Return(nil)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		time.AfterFunc(100*time.Millisecond, app.Exit)
		done <- app.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for Run to complete")
	}

	ml.AssertCalled(t, "Start", mock.Anything)
	ml.AssertCalled(t, "Stop", mock.Anything)

	err = app.Run(ctx)
	if !errors.Is(err, ErrApplicationAlreadyStarted) {
		t.Errorf("expected ErrApplicationAlreadyStarted, got %v", err)
	}
}

func TestRun_StartError(t *testing.T) {
	ctx := context.Background()

	startErr := errors.New("start failed")
	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Return(startErr)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = app.Run(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml.AssertCalled(t, "Start", mock.Anything)

	if !errors.Is(err, startErr) {
		t.Errorf("expected start error, got %v", err)
	}
}

func TestRun_StopError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stopErr := errors.New("stop failed")
	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Return(nil)
	ml.On("Stop", mock.Anything).Return(stopErr)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		app.Exit()
	}()

	err = app.Run(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml.AssertCalled(t, "Start", mock.Anything)
	ml.AssertCalled(t, "Stop", mock.Anything)

	var compErr *ComponentError
	if !errors.As(err, &compErr) {
		t.Errorf("expected ComponentError, got %T", err)
	} else if compErr.Phase != ComponentPhaseStop {
		t.Errorf("expected stop phase, got %s", compErr.Phase)
	}
}

func TestRun_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Return(nil)
	ml.On("Stop", mock.Anything).Return(nil)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err = app.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ml.AssertCalled(t, "Start", mock.Anything)
	ml.AssertCalled(t, "Stop", mock.Anything)
}

func TestRun_Exit(t *testing.T) {
	ctx := context.Background()

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Return(nil)
	ml.On("Stop", mock.Anything).Return(nil)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- app.Run(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	app.Exit()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for Run to complete")
	}

	ml.AssertCalled(t, "Start", mock.Anything)
	ml.AssertCalled(t, "Stop", mock.Anything)
}

func TestRun_Timeout_Start(t *testing.T) {
	ctx := context.Background()

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Run(func(args mock.Arguments) {
		select {
		case <-args.Get(0).(context.Context).Done():
		case <-time.After(200 * time.Millisecond):
		}
	}).Return(func(ctx context.Context) error {
		return ctx.Err()
	})

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml),
		WithStartTimeout(100*time.Millisecond),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- app.Run(ctx)
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for Run to complete")
	}

	ml.AssertCalled(t, "Start", mock.Anything)
	ml.AssertNotCalled(t, "Stop")
}

func TestRun_Timeout_Stop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Return(nil)
	ml.On("Stop", mock.Anything).Run(func(args mock.Arguments) {
		select {
		case <-args.Get(0).(context.Context).Done():
		case <-time.After(200 * time.Millisecond):
		}
	}).Return(func(ctx context.Context) error {
		return ctx.Err()
	})

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml),
		WithStopTimeout(100*time.Millisecond),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	done := make(chan error, 1)
	go func() {
		done <- app.Run(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	app.Exit()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected timeout error, got nil")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for Run to complete")
	}

	ml.AssertCalled(t, "Start", mock.Anything)
	ml.AssertCalled(t, "Stop", mock.Anything)
}

func TestRun_MultipleComponents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Start", mock.Anything).Return(nil)
	ml1.On("Stop", mock.Anything).Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Start", mock.Anything).Return(nil)
	ml2.On("Stop", mock.Anything).Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Start", mock.Anything).Return(nil)
	ml3.On("Stop", mock.Anything).Return(nil)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml1, ml2, ml3),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		app.Exit()
	}()

	err = app.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ml1.AssertCalled(t, "Start", mock.Anything)
	ml2.AssertCalled(t, "Start", mock.Anything)
	ml3.AssertCalled(t, "Start", mock.Anything)

	ml1.AssertCalled(t, "Stop", mock.Anything)
	ml2.AssertCalled(t, "Stop", mock.Anything)
	ml3.AssertCalled(t, "Stop", mock.Anything)
}

func TestSetupSignalHandling(t *testing.T) {
	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(WithLogger(logger))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	shutdownCh := app.setupSignalHandling()

	if shutdownCh == nil {
		t.Fatal("expected shutdownCh, got nil")
	}

	select {
	case <-shutdownCh:
		t.Error("shutdownCh should not be closed immediately")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestSetupSignalHandling_Exit(t *testing.T) {
	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(WithLogger(logger))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	shutdownCh := app.setupSignalHandling()

	go func() {
		time.Sleep(50 * time.Millisecond)
		app.Exit()
	}()

	select {
	case <-shutdownCh:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for shutdownCh to close")
	}
}

func TestSetupSignalHandling_ContextCancel(t *testing.T) {
	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(WithLogger(logger))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	shutdownCh := app.setupSignalHandling()

	go func() {
		time.Sleep(50 * time.Millisecond)
		app.cancel()
	}()

	select {
	case <-shutdownCh:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for shutdownCh to close")
	}
}

func TestSetupSignalHandling_OsSignal(t *testing.T) {
	t.Skip("signal sending tests are unreliable in test suites - covered by integration tests")

	// Signal handling is tested via:
	// - TestSetupSignalHandling_Exit (programmatic shutdown)
	// - TestSetupSignalHandling_ContextCancel (context cancellation)
	// - TestRun_Exit (full lifecycle with Exit)
	// - TestRun_ContextCanceled (full lifecycle with context)
	// - Integration tests with actual signal handling
}
