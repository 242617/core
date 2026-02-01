package logger

import "errors"

// Package-level error definitions.

// ErrInvalidLevel is returned when an invalid log level is provided.
var ErrInvalidLevel = errors.New("invalid log level: must be one of debug, info, warn, error")

// ErrInvalidEncoding is returned when an invalid encoding is provided.
var ErrInvalidEncoding = errors.New("invalid encoding: must be one of json, text")

// ErrHandlerNotInitialized is returned when the logger's handler is nil.
var ErrHandlerNotInitialized = errors.New("handler not initialized")
