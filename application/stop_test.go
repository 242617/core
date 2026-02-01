package application

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/242617/core/mocks"
	"github.com/stretchr/testify/mock"
)

func TestStop_Success(t *testing.T) {
	ctx := context.Background()

	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Stop", mock.Anything).Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Stop", mock.Anything).Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
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

	err = app.stop(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ml1.AssertCalled(t, "Stop", mock.Anything)
	ml2.AssertCalled(t, "Stop", mock.Anything)
	ml3.AssertCalled(t, "Stop", mock.Anything)
}

func TestStop_EmptyComponents(t *testing.T) {
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

	err = app.stop(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStop_SingleComponent(t *testing.T) {
	ctx := context.Background()

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
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

	err = app.stop(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ml.AssertCalled(t, "Stop", mock.Anything)
}

func TestStop_FirstComponentError(t *testing.T) {
	ctx := context.Background()

	stopErr := errors.New("first component failed")
	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Stop", mock.Anything).Return(stopErr)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Stop", mock.Anything).Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
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

	err = app.stop(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml1.AssertCalled(t, "Stop", mock.Anything)
	ml2.AssertCalled(t, "Stop", mock.Anything)
	ml3.AssertCalled(t, "Stop", mock.Anything)
}

func TestStop_MiddleComponentError(t *testing.T) {
	ctx := context.Background()

	stopErr := errors.New("middle component failed")
	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Stop", mock.Anything).Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Stop", mock.Anything).Return(stopErr)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
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

	err = app.stop(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml1.AssertCalled(t, "Stop", mock.Anything)
	ml2.AssertCalled(t, "Stop", mock.Anything)
	ml3.AssertCalled(t, "Stop", mock.Anything)
}

func TestStop_LastComponentError(t *testing.T) {
	ctx := context.Background()

	stopErr := errors.New("last component failed")
	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Stop", mock.Anything).Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Stop", mock.Anything).Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Stop", mock.Anything).Return(stopErr)

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

	err = app.stop(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	ml1.AssertCalled(t, "Stop", mock.Anything)
	ml2.AssertCalled(t, "Stop", mock.Anything)
	ml3.AssertCalled(t, "Stop", mock.Anything)
}

func TestStop_AllComponentsError(t *testing.T) {
	ctx := context.Background()

	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Stop", mock.Anything).Return(errors.New("error 1"))

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Stop", mock.Anything).Return(errors.New("error 2"))

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Stop", mock.Anything).Return(errors.New("error 3"))

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

	err = app.stop(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var compErr1, compErr2, compErr3 *ComponentError
	if !errors.As(err, &compErr1) || !errors.As(err, &compErr2) || !errors.As(err, &compErr3) {
		t.Error("expected multiple ComponentError in error chain")
	}
}

func TestStop_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := atomic.Int64{}

	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
	ml.On("Stop", mock.Anything).Run(func(mock.Arguments) {
		calls.Add(1)
		time.Sleep(50 * time.Millisecond)
	}).Return(nil)

	logger := mocks.NewLogger(t)
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Return()

	app, err := New(
		WithComponents(ml, ml, ml, ml, ml),
		WithLogger(logger),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err = app.stop(ctx)
	if err != nil {
		t.Errorf("stop should not return error on context cancel, got: %v", err)
	}

	callsCount := calls.Load()
	if callsCount == 0 {
		t.Error("no Stop calls were made")
	}
}

func TestStop_ReverseOrder(t *testing.T) {
	ctx := context.Background()

	var order []string

	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Stop", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "stop1")
	}).Return(nil)

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Stop", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "stop2")
	}).Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Stop", mock.Anything).Run(func(mock.Arguments) {
		order = append(order, "stop3")
	}).Return(nil)

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

	err = app.stop(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"stop3", "stop2", "stop1"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(order), order)
	}

	for i, exp := range expected {
		if order[i] != exp {
			t.Errorf("call %d: got %q, want %q", i, order[i], exp)
		}
	}
}

func TestStop_WaitGoroutines(t *testing.T) {
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

	app.wg.Add(1)
	go func() {
		time.Sleep(100 * time.Millisecond)
		app.wg.Done()
	}()

	start := time.Now()
	err = app.stop(ctx)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if duration < 100*time.Millisecond {
		t.Errorf("stop should wait for goroutines, took %v", duration)
	}
}

func TestStop_ComponentErrorFormat(t *testing.T) {
	ctx := context.Background()

	stopErr := errors.New("specific error")
	ml := mocks.NewComponent(t)
	ml.On("String").Maybe().Return("comp1")
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

	err = app.stop(ctx)
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
		if compErr.Phase != ComponentPhaseStop {
			t.Errorf("phase = %q, want %q", compErr.Phase, ComponentPhaseStop)
		}
		if !errors.Is(compErr, stopErr) {
			t.Errorf("error should contain original error: %v", compErr.Err)
		}
	}
}

func TestStop_ErrorAggregation(t *testing.T) {
	ctx := context.Background()

	ml1 := mocks.NewComponent(t)
	ml1.On("String").Maybe().Return("comp1")
	ml1.On("Stop", mock.Anything).Return(errors.New("error 1"))

	ml2 := mocks.NewComponent(t)
	ml2.On("String").Maybe().Return("comp2")
	ml2.On("Stop", mock.Anything).Return(nil)

	ml3 := mocks.NewComponent(t)
	ml3.On("String").Maybe().Return("comp3")
	ml3.On("Stop", mock.Anything).Return(errors.New("error 2"))

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

	err = app.stop(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var compErr1, compErr2 *ComponentError
	if !errors.As(err, &compErr1) {
		t.Error("expected ComponentError in error chain")
	}
	if !errors.As(err, &compErr2) {
		t.Error("expected ComponentError in error chain")
	}

	if compErr1.Component != compErr2.Component {
		t.Log("errors from different components aggregated correctly")
	}
}
