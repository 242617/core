package protocol

import "context"

// Logger defines the logging interface for protocol-level logging operations.
// It provides structured logging with context propagation and runtime configuration.
type Logger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

// NopLogger is a no-op logger that discards all log messages.
// It implements the protocol.Logger interface and is useful for testing
// or when logging should be completely disabled.
type NopLogger struct{}

func (NopLogger) Debug(_ context.Context, _ string, _ ...any) {}
func (NopLogger) Info(_ context.Context, _ string, _ ...any)  {}
func (NopLogger) Warn(_ context.Context, _ string, _ ...any)  {}
func (NopLogger) Error(_ context.Context, _ string, _ ...any) {}
