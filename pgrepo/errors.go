package pgrepo

import (
	"errors"
)

var (
	// ErrDatabaseNotStarted indicates that the database has not been started.
	ErrDatabaseNotStarted = errors.New("database not started")

	// ErrInvalidConfig indicates that the configuration is invalid.
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrShutdownTimeout indicates that the graceful shutdown timed out.
	ErrShutdownTimeout = errors.New("shutdown timeout")
)
