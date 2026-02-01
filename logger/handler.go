package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/242617/core/request_id"
)

const (
	FieldNameRequestID = "request_id"
	FieldNameDuration  = "duration"
)

// contextHandler is a custom slog.Handler that adds context-aware fields
// to log records, such as request IDs from the context.
type contextHandler struct {
	slog.Handler
}

// Handle implements slog.Handler, adding context fields before delegating to the wrapped handler.
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add request_id from context if available
	if requestID := request_id.RequestIDFromContext(ctx); requestID != "" {
		r.Add(slog.String(FieldNameRequestID, requestID))
	}

	// Add other context fields here as needed
	// For example: trace IDs, user IDs, etc.

	return h.Handler.Handle(ctx, r)
}

// Timer is a utility for measuring and logging operation duration.
//
// Example:
//
//	timer := log.NewTimer(logger, ctx, "database query")
//	// ... perform operation ...
//	timer.Stop("query completed", "rows", result.RowsAffected())
type Timer struct {
	logger    *Logger
	ctx       context.Context
	message   string
	startTime time.Time
}

// NewTimer creates a new timer that starts immediately.
// Call Stop() to log the elapsed time.
func NewTimer(logger *Logger, ctx context.Context, message string) *Timer {
	return &Timer{
		logger:    logger,
		ctx:       ctx,
		message:   message,
		startTime: time.Now(),
	}
}

// Stop logs the elapsed time with the given message and optional additional fields.
// It returns the elapsed duration for further processing if needed.
//
// Example:
//
//	timer := log.NewTimer(logger, ctx, "database query")
//	defer func() {
//	    timer.Stop("query completed", "rows", rows)
//	}()
func (t *Timer) Stop(msg string, args ...any) time.Duration {
	duration := time.Since(t.startTime)
	logArgs := append(args, FieldNameDuration, duration)
	t.logger.Info(t.ctx, msg, logArgs...)
	return duration
}

// Debug stops the timer and logs at debug level.
func (t *Timer) Debug(msg string, args ...any) time.Duration {
	duration := time.Since(t.startTime)
	logArgs := append(args, FieldNameDuration, duration)
	t.logger.Debug(t.ctx, msg, logArgs...)
	return duration
}

// Error stops the timer and logs at error level.
// This is useful when an operation fails.
func (t *Timer) Error(msg string, args ...any) time.Duration {
	duration := time.Since(t.startTime)
	logArgs := append(args, FieldNameDuration, duration)
	t.logger.Error(t.ctx, msg, logArgs...)
	return duration
}
