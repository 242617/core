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
	"github.com/stretchr/testify/require"

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
				require.NoError(t, err, "no error")
			})

		assert.Equal(t, 1, first.Called(), "first called once")
		assert.Equal(t, 1, second.Called(), "second called once")
		assert.Equal(t, 0, third.Called(), "third never called")
	}

	{
		sampleErr := errors.New("sample error")
		errFunc := withError{sampleErr}
		pipeline.New(context.Background(), errFunc.Call).
			Run(func(err error) {
				require.ErrorIs(t, err, sampleErr, "sample error")
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
				require.ErrorIs(t, err, sampleErr, "sample error")
			})
		assert.Equal(t, 0, first.Called(), "first never called")
		assert.Equal(t, 0, second.Called(), "second never called")
		assert.Equal(t, 0, third.Called(), "third never called")
		assert.Equal(t, 0, fourth.Called(), "fourth never called")
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
				require.ErrorIs(t, err, sampleErr, "sample error")
			})

		assert.Equal(t, 0, first.Called(), "first never called")
		assert.Equal(t, 0, second.Called(), "second never called")
		assert.Equal(t, 0, third.Called(), "third never called")
		assert.Equal(t, 0, fourth.Called(), "fourth never called")
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
				require.NoError(t, err, sampleErr, "no error")
			})
		assert.Equal(t, 1, first.Called(), "first called once")
		assert.Equal(t, 1, second.Called(), "second called once")
		assert.Equal(t, 1, third.Called(), "third called once")
		assert.Equal(t, 0, fourth.Called(), "fourth never called")
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

		assert.Equal(t, 0, next.Called(), "next never called")
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
				require.Equal(t, "context canceled", err.Error(), "context canceled")
			})
		assert.Equal(t, 0, third.Called(), "third never called")
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

		assert.Equal(t, 0, next.Called(), "next never called")
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
		assert.Equal(t, 0, third.Called(), "third never called")
	}
}

func TestAll(t *testing.T) {
	{ // successful
		var first, second, third withCallCounter
		pipeline.New(context.Background(), first.Call).
			Then(second.Call, third.Call).
			Run(func(err error) {
				require.NoError(t, err, "no error")
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
				require.True(t, strings.Contains(err.Error(), "sample"), "sample")
			})

		assert.Equal(t, 0, second.Called(), "second never called")
		assert.Equal(t, 0, third.Called(), "third never called")
	}

	{ // waiting for errors
		firstErr := errors.New("first")
		first, second := withErrorAfter{period, firstErr}, withErrorAfter{2 * period, nil}
		pipeline.New(context.Background(), new(withCallCounter).Call).
			Then(first.Call, second.Call).
			Run(func(err error) {
				require.ErrorIs(t, err, firstErr, "first error")
			})
	}
}

func TestBeforeAndAfter(t *testing.T) {
	first, second, third := withCallCounter{}, withCallCounter{}, withCallCounter{}
	pipeline.New(context.Background()).
		Before(func() { _ = first.Call(context.Background()) }).
		Before(func() {
			assert.Equal(t, 0, first.Called(), "first never called")
			assert.Equal(t, 0, second.Called(), "second never called")
			assert.Equal(t, 0, third.Called(), "third never called")
		}).
		Then(second.Call).
		After(func() {
			assert.Equal(t, 0, first.Called(), "first never called")
			assert.Equal(t, 1, second.Called(), "second called once")
			assert.Equal(t, 0, third.Called(), "third never called")
		}).
		Before(func() {
			assert.Equal(t, 0, first.Called(), "first never called")
			assert.Equal(t, 1, second.Called(), "second called once")
			assert.Equal(t, 0, third.Called(), "third never called")
		}).
		After(func() {
			assert.Equal(t, 0, first.Called(), "first never called")
			assert.Equal(t, 1, second.Called(), "second called once")
			assert.Equal(t, 1, third.Called(), "third called once")
		}).
		Then(third.Call).
		Run(func(err error) {
			require.NoError(t, err, "expect no error")
		})
}

func TestErrorNoError(t *testing.T) {
	{
		first := withEmpty{}
		second, third := withCallCounter{}, withCallCounter{}
		fourthErr := errors.New("fourth")
		fourth, fifth := withError{fourthErr}, withCallCounter{}
		pipeline.New(context.Background()).
			Then(first.Call).
			Error(func(err error) error { return second.Call(context.Background()) }).
			NoError(func() error { return third.Call(context.Background()) }).
			Then(fourth.Call).
			Then(fifth.Call).
			Run(func(err error) {
				require.ErrorIs(t, err, fourthErr, "fourth error")
			})
		assert.Equal(t, 0, second.Called(), "second never called")
		assert.Equal(t, 1, third.Called(), "third called once")
		assert.Equal(t, 0, fifth.Called(), "fifth never called")
	}
	{
		firstErr := errors.New("first")
		first := withError{firstErr}
		second, third := withCallCounter{}, withCallCounter{}
		pipeline.New(context.Background()).
			Then(first.Call).
			Error(func(err error) error {
				require.ErrorIs(t, err, firstErr, "first error")
				return nil
			}).
			NoError(func() error { return second.Call(context.Background()) }).
			Then(third.Call).
			Run(func(err error) {
				require.NoError(t, err, "expect no error")
			})
		assert.Equal(t, 1, second.Called(), "second called once")
		assert.Equal(t, 1, third.Called(), "third called once")
	}
}

func TestCatches(t *testing.T) {
	{
		var one, two bool
		firstErr, secondErr := errors.New("first"), errors.New("second")
		first, second := withError{firstErr}, withError{secondErr}
		third := withCallCounter{}
		pipeline.New(context.Background()).
			Then(first.Call).
			ThenCatch(func(err error) error {
				one = true
				return err
			}).
			Else(second.Call).
			ElseCatch(func(err error) error {
				two = true
				return err
			}).
			Then(third.Call).
			Run(func(err error) {
				require.ErrorIs(t, err, secondErr, "second error")
			})
		assert.Equal(t, 0, third.Called(), "third never called")
		assert.True(t, one, "unexpected one value")
		assert.True(t, two, "unexpected two value")
	}

	{
		firstBefore, firstAfter := withCallCounter{}, withCallCounter{}
		firstErr := errors.New("first")
		first := withError{firstErr}
		firstElseCatch := withCallCounter{}
		secondBefore, secondAfter := withCallCounter{}, withCallCounter{}
		second := withCallCounter{}
		secondThenCatch := withCallCounter{}
		pipeline.New(context.Background()).
			Before(func() { _ = firstBefore.Call(context.Background()) }).
			Then(first.Call).
			ThenCatch(func(err error) error {
				require.ErrorIs(t, err, firstErr, "unexpected first error")
				return err
			}).
			ElseCatch(func(error) error { return firstElseCatch.Call(context.Background()) }).
			After(func() { _ = firstAfter.Call(context.Background()) }).
			Before(func() { _ = secondBefore.Call(context.Background()) }).
			Then(second.Call).
			ThenCatch(func(err error) error { return secondThenCatch.Call(context.Background()) }).
			After(func() { _ = secondAfter.Call(context.Background()) }).
			Run(func(err error) {
				require.ErrorIs(t, err, firstErr, "first error")
			})
		assert.Equal(t, 1, firstBefore.Called(), "firstBefore called once")
		assert.Equal(t, 0, firstElseCatch.Called(), "firstAfter never called")
		assert.Equal(t, 1, firstAfter.Called(), "firstAfter called once")
		assert.Equal(t, 0, secondBefore.Called(), "secondBefore never called")
		assert.Equal(t, 0, second.Called(), "second never called")
		assert.Equal(t, 0, secondThenCatch.Called(), "secondThenCatch never called")
		assert.Equal(t, 0, secondAfter.Called(), "secondAfter never called")
	}

	{ // Fall-through after then catches
		sampleErr := errors.New("sample error")
		noError, sampleError := withEmpty{}, withError{sampleErr}
		pipeline.New(context.Background()).
			Then(noError.Call).
			ThenCatch(func(err error) error {
				assert.NoError(t, err, "expect no error")
				return err
			}).
			Then(sampleError.Call).
			Run(func(err error) {
				require.ErrorIs(t, err, sampleErr, "sample error")
			})
	}
}

func TestAppend(t *testing.T) {
	var numbers []string
	p1 := pipeline.New(context.Background()).
		Then(func(ctx context.Context) error {
			numbers = append(numbers, "one")
			return nil
		})
	p2 := pipeline.New(context.Background()).
		Then(func(ctx context.Context) error {
			numbers = append(numbers, "two")
			return nil
		})
	p3 := pipeline.New(context.Background()).
		Then(func(ctx context.Context) error {
			numbers = append(numbers, "three")
			return nil
		})
	fourthErr := errors.New("fourth")
	p4 := pipeline.New(context.Background()).
		Then(func(ctx context.Context) error {
			numbers = append(numbers, "four")
			return fourthErr
		})

	p1.Append(p2, p3).
		Run(func(err error) {
			require.NoError(t, err, "no error")
			require.Equal(t, "one,two,three", strings.Join(numbers, ","), "unexpected after appending p2 and p3")
		})
	numbers = []string{}

	p1.Append(p4).
		Run(func(err error) {
			require.ErrorIs(t, err, fourthErr, "unexpected error after appending p4")
			require.Equal(t, "one,four", strings.Join(numbers, ","), "unexpected after appending p4")
		})
	numbers = []string{}

	p1.Run(func(err error) {
		require.NoError(t, err, "no error")
		require.Equal(t, "one", strings.Join(numbers, ","), "unexpected p1")
	})
	numbers = []string{}
}

func TestMerge(t *testing.T) {
	var numbers []string
	p1 := pipeline.New(context.Background()).
		Then(func(ctx context.Context) error {
			numbers = append(numbers, "one")
			return nil
		}).Name("one")

	tail := func(name string) *pipeline.Pipeline {
		return pipeline.New(context.Background()).
			Then(func(ctx context.Context) error {
				numbers = append(numbers, name)
				return nil
			}).Name(name)
	}

	p1.
		Then(func(ctx context.Context) error {
			numbers = append(numbers, "two")
			return nil
		}).Name("two").
		Merge(func() *pipeline.Pipeline { return tail("three") }).
		Run(func(err error) {
			require.NoError(t, err, "no error")
			assert.Equal(t, "one,two,three", strings.Join(numbers, ","), "unexpected")
		})
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
