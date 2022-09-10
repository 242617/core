package pipeline_test

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/242617/core/pipeline"
)

var period = 10 * time.Millisecond

func init() { log.SetFlags(log.Lshortfile) }
func TestBasic(t *testing.T) {
	{
		var first, second, third withCallCounter
		pipeline.New(context.Background(), first.Call).
			Then(second.Call).
			Else(third.Call).
			Run(func(err error) {
				assert.NoError(t, err, "no error")
			})

		assert.Equal(t, 1, first.Called(), "first called once")
		assert.Equal(t, 1, second.Called(), "second called once")
		assert.Equal(t, 0, third.Called(), "third not called")
	}

	{
		sampleErr := errors.New("sample error")
		errFunc := withError{sampleErr}
		pipeline.New(context.Background(), errFunc.Call).
			Run(func(err error) {
				assert.ErrorIs(t, err, sampleErr, "sample error")
			})
	}

	{
		sampleErr := errors.New("sample error")
		var first, second, third, fourth withCallCounter
		errFunc := withError{sampleErr}
		pipeline.New(context.Background(), errFunc.Call).
			Then(first.Call).
			Then(second.Call).
			Else(third.Call).
			Then(fourth.Call).
			Run(func(err error) {
				assert.ErrorIs(t, err, sampleErr, "sample error")
			})
		assert.Equal(t, 0, first.Called(), "first not called")
		assert.Equal(t, 0, second.Called(), "second not called")
		assert.Equal(t, 0, third.Called(), "third not called")
		assert.Equal(t, 0, fourth.Called(), "fourth not called")
	}

	{
		sampleErr := errors.New("sample error")
		var first, second, third withCallCounter
		errFunc := withError{sampleErr}
		pipeline.New(context.Background(), errFunc.Call).
			Then(first.Call).
			Reset().
			Then(second.Call).
			Else(third.Call).
			Run(func(err error) {
				assert.NoError(t, err, "no error")
			})
		assert.Equal(t, 0, first.Called(), "first not called")
		assert.Equal(t, 1, second.Called(), "second not called")
		assert.Equal(t, 0, third.Called(), "third not called")
	}

	{
		sampleErr := errors.New("sample error")
		var first, second, third, fourth withCallCounter
		pipeline.New(context.Background()).
			Then(func(context.Context) error { return sampleErr }).
			Then(first.Call).
			Else(second.Call).
			Else(third.Call).
			Then(fourth.Call).
			Run(func(err error) {
				assert.ErrorIs(t, err, sampleErr, "sample error")
			})

		assert.Equal(t, 0, first.Called(), "first not called")
		assert.Equal(t, 0, second.Called(), "second not called")
		assert.Equal(t, 0, third.Called(), "third not called")
		assert.Equal(t, 0, fourth.Called(), "fourth not called")
	}

	{
		sampleErr := errors.New("sample error")
		var first, second, third, fourth withCallCounter
		pipeline.New(context.Background(), new(withEmpty).Call).
			Then(func(context.Context) error { return sampleErr }).
			Else(func(context.Context) error { return nil }).
			Then(first.Call).
			Then(second.Call).
			Then(func(context.Context) error { return sampleErr }).
			Else(third.Call).
			Else(fourth.Call).
			Run(func(err error) {
				assert.NoError(t, err, sampleErr, "no error")
			})
		assert.Equal(t, 1, first.Called(), "first called once")
		assert.Equal(t, 1, second.Called(), "second called once")
		assert.Equal(t, 1, third.Called(), "third called once")
		assert.Equal(t, 0, fourth.Called(), "fourth not called")
	}
}

func TestContextCancel(t *testing.T) {
	{
		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(period, cancel)

		wait := withTimeout{2 * period}
		var next withCallCounter
		var summary string
		pipeline.New(ctx, wait.Call).
			Then(next.Call).
			Run(func(err error) { summary = err.Error() })

		assert.Equal(t, 0, next.Called(), "next not called")
		assert.Equal(t, "context canceled", summary, "context canceled")
	}

	{
		ctx, cancel := context.WithCancel(context.Background())
		time.AfterFunc(period, cancel)

		first, second := withTimeout{2 * period}, withTimeout{3 * period}
		var third withCallCounter
		pipeline.New(ctx, new(withEmpty).Call).
			Then(first.Call, second.Call).
			Then(third.Call).
			Run(func(err error) {
				assert.Equal(t, "context canceled", err.Error(), "context canceled")
			})
		assert.Equal(t, 0, third.Called(), "third not called")
	}
}

func TestContextTimeout(t *testing.T) {
	{
		ctx, cancel := context.WithTimeout(context.Background(), period)
		defer cancel()

		wait := withTimeout{2 * period}
		var next withCallCounter
		var summary string
		pipeline.New(ctx, wait.Call).
			Then(next.Call).
			Run(func(err error) { summary = err.Error() })

		assert.Equal(t, 0, next.Called(), "next not called")
		assert.Equal(t, "context deadline exceeded", summary, "context deadline exceeded")
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), period)
		defer cancel()

		first, second := withTimeout{2 * period}, withTimeout{3 * period}
		var third withCallCounter
		pipeline.New(ctx, new(withEmpty).Call).
			Then(first.Call, second.Call).
			Then(third.Call).
			Run(func(err error) {
				assert.Equal(t, "context deadline exceeded", err.Error(), "context deadline exceeded")
			})
		assert.Equal(t, 0, third.Called(), "third not called")
	}
}

func TestAll(t *testing.T) {
	{ // successful
		var first, second, third withCallCounter
		pipeline.New(context.Background(), first.Call).
			Then(second.Call, third.Call).
			Run(func(err error) {
				assert.NoError(t, err, "no error")
			})

		assert.Equal(t, 1, first.Called(), "first called once")
		assert.Equal(t, 1, second.Called(), "second called once")
		assert.Equal(t, 1, third.Called(), "third called once")
	}

	{ // successful
		first := withErrorMessage{"sample"}
		var second, third withCallCounter
		pipeline.New(context.Background(), first.Call).
			Then(second.Call, third.Call).
			Run(func(err error) {
				assert.True(t, strings.Contains(err.Error(), "sample"), "sample")
			})

		assert.Equal(t, 0, second.Called(), "second not called")
		assert.Equal(t, 0, third.Called(), "third not called")
	}

	{ // waiting for errors
		firstErr := errors.New("first")
		first, second := withErrorAfter{period, firstErr}, withErrorAfter{2 * period, nil}
		pipeline.New(context.Background(), new(withCallCounter).Call).
			Then(first.Call, second.Call).
			Run(func(err error) {
				assert.ErrorIs(t, err, firstErr, "first error")
			})
	}
}

type (
	withEmpty        struct{}
	withError        struct{ err error }
	withErrorMessage struct{ msg string }
	withCallCounter  struct {
		sync.RWMutex
		n int
	}
	withTimeout    struct{ d time.Duration }
	withErrorAfter struct {
		d   time.Duration
		err error
	}
)

func (w *withEmpty) Call(context.Context) error { return nil }

func (w *withError) Call(context.Context) error { return w.err }

func (w *withErrorMessage) Call(context.Context) error { return errors.New(w.msg) }

func (w *withCallCounter) Call(context.Context) error {
	w.Lock()
	defer w.Unlock()
	w.n += 1
	return nil
}
func (w *withCallCounter) Called() int {
	w.RLock()
	defer w.RUnlock()
	return w.n
}

func (w *withTimeout) Call(context.Context) error {
	time.Sleep(w.d)
	return nil
}

func (a *withErrorAfter) Call(context.Context) error {
	time.Sleep(a.d)
	return a.err
}
