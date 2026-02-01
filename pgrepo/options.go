package pgrepo

import (
	"github.com/242617/core/protocol"
	"github.com/pkg/errors"
)

// Option is a function that modifies the DB configuration.
type Option func(*DB) error

// defaults returns an option that provides default values for DB
func defaults() []Option {
	return []Option{
		WithLogger(protocol.NopLogger{}),
	}
}

func WithConfig(cfg Config) Option {
	return func(db *DB) error {
		if err := cfg.Validate(); err != nil {
			return errors.Wrap(err, "invalid config")
		}
		db.cfg = cfg
		return nil
	}
}

func WithLogger(log protocol.Logger) Option {
	return func(db *DB) error {
		if log == nil {
			return errors.New("logger cannot be nil")
		}
		db.log = log
		return nil
	}
}
