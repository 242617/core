package application

import (
	"context"
	"errors"
	"testing"

	"github.com/242617/core/mocks"
	"github.com/stretchr/testify/mock"
)

func TestStart_Success(t *testing.T) {
	ctx := context.Background()

	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Start", mock.Anything).Return(nil)
	ml1.On("Stop", mock.Anything).Maybe().Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Start", mock.Anything).Return(nil)
	ml2.On("Stop", mock.Anything).Maybe().Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Start", mock.Anything).Return(nil)
	ml3.On("Stop", mock.Anything).Maybe().Return(nil)

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

	err = app.start(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ml1.AssertCalled(t, "Start", mock.Anything)
	ml2.AssertCalled(t, "Start", mock.Anything)
	ml3.AssertCalled(t, "Start", mock.Anything)
}

func TestStart_EmptyComponents(t *testing.T) {
	ctx := context.Background()

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = app.start(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStart_SingleComponent(t *testing.T) {
	ctx := context.Background()

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Return(nil)
	ml.On("Stop", mock.Anything).Maybe().Return(nil)

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

	err = app.start(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ml.AssertCalled(t, "Start", mock.Anything)
}

func TestStart_FirstComponentError(t *testing.T) {
	ctx := context.Background()

	startErr := errors.New("first component failed")
	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Start", mock.Anything).Return(startErr)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")

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

	err = app.start(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml1.AssertCalled(t, "Start", mock.Anything)
	ml2.AssertNotCalled(t, "Start", mock.Anything)
	ml3.AssertNotCalled(t, "Start", mock.Anything)
}

func TestStart_MiddleComponentError(t *testing.T) {
	ctx := context.Background()

	startErr := errors.New("middle component failed")
	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Start", mock.Anything).Return(nil)
	ml1.On("Stop", mock.Anything).Maybe().Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Start", mock.Anything).Return(startErr)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")

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

	err = app.start(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml1.AssertCalled(t, "Start", mock.Anything)
	ml2.AssertCalled(t, "Start", mock.Anything)
	ml3.AssertNotCalled(t, "Start", mock.Anything)

	ml1.AssertCalled(t, "Stop", mock.Anything)
}

func TestStart_LastComponentError(t *testing.T) {
	ctx := context.Background()

	startErr := errors.New("last component failed")
	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Start", mock.Anything).Return(nil)
	ml1.On("Stop", mock.Anything).Maybe().Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Start", mock.Anything).Return(nil)
	ml2.On("Stop", mock.Anything).Maybe().Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Start", mock.Anything).Return(startErr)

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

	err = app.start(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml1.AssertCalled(t, "Start", mock.Anything)
	ml2.AssertCalled(t, "Start", mock.Anything)
	ml3.AssertCalled(t, "Start", mock.Anything)

	ml1.AssertCalled(t, "Stop", mock.Anything)
	ml2.AssertCalled(t, "Stop", mock.Anything)
}

func TestStart_RollbackErrors(t *testing.T) {
	ctx := context.Background()

	startErr := errors.New("component failed")
	stopErr := errors.New("stop failed")
	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Start", mock.Anything).Return(nil)
	ml1.On("Stop", mock.Anything).Return(stopErr)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Start", mock.Anything).Return(startErr)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")

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

	err = app.start(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml1.AssertCalled(t, "Start", mock.Anything)
	ml2.AssertCalled(t, "Start", mock.Anything)

	ml1.AssertCalled(t, "Stop", mock.Anything)
}

func TestStart_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	called := false

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Start", mock.Anything).Run(func(args mock.Arguments) {
		cancel()
		called = true
		<-args.Get(0).(context.Context).Done()
	}).Return(context.Canceled)

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

	err = app.start(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !called {
		t.Error("Start was not called")
	}
}

func TestStart_RollbackOrder(t *testing.T) {
	ctx := context.Background()

	var order []string

	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Start", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "start1")
	}).Return(nil)
	ml1.On("Stop", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "stop1")
	}).Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Start", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "start2")
	}).Return(nil)
	ml2.On("Stop", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "stop2")
	}).Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Start", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "start3")
	}).Return(errors.New("error"))

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

	_ = app.start(ctx)

	expected := []string{"start1", "start2", "start3", "stop2", "stop1"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(order), order)
	}

	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("call %d: got %q, want %q", i, order[i], exp)
		}
	}
}

func TestStart_ComponentErrorFormat(t *testing.T) {
	ctx := context.Background()

	startErr := errors.New("specific error")
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

	err = app.start(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var compErr *ComponentError
	if !errors.As(err, &compErr) {
		t.Errorf("expected ComponentError, got %T", err)
	} else {
		if compErr.Component != "comp1" {
			t.Errorf("component name = %q, want 'comp1'", compErr.Component)
		}
		if compErr.Phase != ComponentPhaseStart {
			t.Errorf("phase = %q, want %q", compErr.Phase, ComponentPhaseStart)
		}
		if !errors.Is(compErr, startErr) {
			t.Errorf("error should contain original error: %v", compErr.Err)
		}
	}
}
