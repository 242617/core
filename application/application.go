package application

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	l "github.com/rs/zerolog/log"

	"github.com/242617/core/protocol"
)

type option = func(a *Application) error

func withDefaultTimeouts() option {
	return func(a *Application) error {
		a.startTimeout, a.stopTimeout = time.Second, time.Second
		return nil
	}
}
func WithStartTimeout(timeout time.Duration) option {
	return func(a *Application) error {
		a.startTimeout = timeout
		return nil
	}
}

func withDefaultLogger() option {
	return func(a *Application) error {
		a.log = l.With().Str("component", "application").Logger()
		return nil
	}
}

func WithComponents(components ...Component) option {
	return func(a *Application) error {
		a.components = components
		return nil
	}
}

func New(options ...option) (*Application, error) {
	var a Application
	options = append([]option{
		withDefaultTimeouts(),
		withDefaultLogger(),
	}, options...)
	for _, option := range options {
		if err := option(&a); err != nil {
			return nil, errors.New("apply option")
		}
	}
	return &a, nil
}

type Application struct {
	startTimeout, stopTimeout time.Duration
	log                       zerolog.Logger
	components                []Component
}

type Component interface {
	fmt.Stringer
	protocol.Lifecycle
}

func NewLifecycleComponent(name string, cmp protocol.Lifecycle) *LifecycleComponent {
	return &LifecycleComponent{name, cmp}
}

type LifecycleComponent struct {
	string
	protocol.Lifecycle
}

func (s *LifecycleComponent) String() string { return s.string }

type ContextFunc = func(context.Context) error

type MethodsComponent struct {
	name        string
	start, stop ContextFunc
}

func NewMethodsComponent(name string, start, stop ContextFunc) MethodsComponent {
	return MethodsComponent{
		name:  name,
		start: start,
		stop:  stop,
	}
}

func (c MethodsComponent) Start(ctx context.Context) error { return c.call(ctx, c.start) }
func (c MethodsComponent) Stop(ctx context.Context) error  { return c.call(ctx, c.stop) }

func (c MethodsComponent) call(ctx context.Context, f ContextFunc) error {
	if f == nil {
		return nil
	}
	return f(ctx)
}

func (c MethodsComponent) String() string { return c.name }

func PlainToContextFunc(f func()) ContextFunc {
	return func(context.Context) error {
		f()
		return nil
	}
}
