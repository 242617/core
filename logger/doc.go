// Package logger provides a lightweight, production-ready structured logging library
// built on top of Go's standard library slog (Go 1.21+).
//
// Features:
//   - Structured logging with key-value attributes
//   - Automatic context propagation (request IDs, trace IDs)
//   - Support for JSON and text encodings with colored output
//   - Named loggers for module-level separation
//   - Flexible configuration via modifiers pattern
//   - Built-in timer for measuring operation duration
//   - Dynamic log level adjustment (runtime)
//
// Basic usage:
//
//	import "github.com/242617/core/logger"
//
//	// Create logger with default production config
//	logger, err := logger.New()
//	if err != nil {
//	    panic(err)
//	}
//
//	// Create named child logger
//	serviceLogger := logger.New("my-service")
//
//	// Log with context
//	serviceLogger.Info(ctx, "handling request", "method", "GET", "path", "/api/v1")
//
//	// Log errors
//	serviceLogger.Error(ctx, "failed to process", "err", err, "user_id", userID)
//
// Advanced usage with modifiers:
//
//	logger, err := logger.New(
//	    logger.WithConfig(logger.Config{
//	        Level:    logger.LevelDebug,
//	        Encoding: logger.EncodingText,
//	        Colorize: true,
//	    }),
//	)
//
// Configuration:
//
// The logger supports YAML configuration:
//
//	log:
//	  level: debug        # debug, info, warn, error
//	  encoding: json      # json, text
//	  colorize: true      # colored text output (text encoding only)
//
// For more examples, see examples/main.go.
package logger
