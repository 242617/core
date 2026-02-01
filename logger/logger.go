package logger

import (
	"context"
	"fmt"
	"log/slog"
)

// New creates a new logger with the provided options.
// If no options are provided, uses ProductionConfig by default.
//
// Example:
//
//	logger, err := log.New()
//	if err != nil {
//	    panic(err)
//	}
//
//	logger, err := log.New(log.WithConfig(log.DevelopmentConfig))
func New(options ...Option) (*Logger, error) {
	var logger Logger

	// Apply default config first, then user options
	for _, option := range append([]Option{withDefaultConfig()}, options...) {
		if err := option(&logger); err != nil {
			return nil, err
		}
	}

	if logger.handler == nil {
		return nil, ErrHandlerNotInitialized
	}

	logger.logger = slog.New(logger.handler)

	return &logger, nil
}

// Logger is a wrapper around slog.Logger with additional functionality.
// It should be created via New() and configured via options.
type Logger struct {
	name    string
	handler slog.Handler
	logger  *slog.Logger
	config  Config // Store config for level manipulation
}

// applyLabels prepends logger name to the log attributes if it exists.
func (l *Logger) applyLabels(args ...any) []any {
	var labels []any
	if l.name != "" {
		labels = append(labels, "name", l.name)
	}
	return append(labels, args...)
}

// Debug logs a message at debug level.
// Use for detailed diagnostic information.
func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, l.applyLabels(args...)...)
}

// Info logs a message at info level.
// Use for general informational messages about normal operation.
func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, l.applyLabels(args...)...)
}

// Warn logs a message at warn level.
// Use for potentially harmful situations that don't prevent the application from running.
func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, l.applyLabels(args...)...)
}

// Error logs a message at error level.
// Use for error events that might still allow the application to continue running.
func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, l.applyLabels(args...)...)
}

// New creates a new named child logger.
// The child logger inherits the parent's handler and configuration,
// but includes its name in all log entries.
//
// Example:
//
//	serviceLogger := logger.New("auth-service")
//	handlerLogger := serviceLogger.New("http-handler")
func (l *Logger) New(name string) *Logger {
	return &Logger{
		name:    name,
		logger:  slog.New(l.handler),
		handler: l.handler,
		config:  l.config,
	}
}

// SetLevel dynamically changes the log level at runtime.
// This recreates the handler with the new level.
//
// Example:
//
//	err := logger.SetLevel(log.LevelDebug)
func (l *Logger) SetLevel(level string) error {
	// Validate the new level
	if level != LevelDebug && level != LevelInfo && level != LevelWarn && level != LevelError {
		return fmt.Errorf("invalid log level: %s, must be one of debug, info, warn, error", level)
	}

	// Create a new config with the updated level
	newConfig := l.config.Clone()
	newConfig.Level = level

	// Create a new handler with the new level
	handler, err := newConfig.handler()
	if err != nil {
		return fmt.Errorf("failed to create handler with new level: %w", err)
	}

	// Update the logger with the new handler
	l.handler = handler
	l.logger = slog.New(handler)
	l.config = newConfig

	return nil
}

// LogLevel returns the current log level as a string.
func (l *Logger) LogLevel() string {
	return l.config.Level
}

// Enabled returns true if logging at given level is enabled.
// This can be used to avoid expensive operations when the log level is disabled.
//
// Example:
//
//	if logger.Enabled(ctx, log.LevelDebug) {
//	    logger.Debug(ctx, "expensive operation result", "data", computeExpensiveData())
//	}
func (l *Logger) Enabled(ctx context.Context, level slog.Level) bool {
	return l.logger.Enabled(ctx, level)
}

// SetConfig allows updating the logger's configuration.
// This recreates the handler with the new configuration.
func (l *Logger) SetConfig(cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	handler, err := cfg.handler()
	if err != nil {
		return fmt.Errorf("failed to create handler: %w", err)
	}

	l.handler = handler
	l.logger = slog.New(handler)
	l.config = cfg

	return nil
}
