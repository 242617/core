package application

import (
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/242617/core/protocol"
)

type Option func(*Application) error

func defaults() []Option {
	return []Option{
		WithName("application"),
		WithLogger(protocol.NopLogger{}),
		WithStartTimeout(30 * time.Second),
		WithStopTimeout(30 * time.Second),
		withDefaultHostname(),
	}
}

// WithName sets application name.
func WithName(name string) Option {
	return func(app *Application) error {
		app.name = name
		return nil
	}
}

// WithStartTimeout sets startup timeout.
func WithStartTimeout(timeout time.Duration) Option {
	return func(app *Application) error {
		if timeout <= 0 {
			return errors.New("invalid start timeout")
		}
		app.startTimeout = timeout
		return nil
	}
}

// WithStopTimeout sets shutdown timeout.
func WithStopTimeout(timeout time.Duration) Option {
	return func(app *Application) error {
		if timeout <= 0 {
			return errors.New("invalid stop timeout")
		}
		app.stopTimeout = timeout
		return nil
	}
}

// WithLogger sets logger instance.
func WithLogger(logger protocol.Logger) Option {
	return func(app *Application) error {
		if logger == nil {
			return errors.New("empty logger")
		}
		app.log = logger
		return nil
	}
}

// WithComponents sets components to manage.
func WithComponents(components ...Component) Option {
	return func(app *Application) error {
		app.components = components
		return nil
	}
}

func withDefaultHostname() Option {
	return func(app *Application) error {
		hostname, err := os.Hostname()
		if err != nil {
			return errors.Wrap(err, "os hostname")
		}
		return WithHostname(hostname)(app)
	}
}
func WithHostname(hostname string) Option {
	return func(app *Application) error {
		if hostname == "" {
			return errors.New("empty hostname")
		}
		app.hostname = hostname
		return nil
	}
}
