package logger

import "fmt"

// Option is a function that modifies a Logger during creation.
// This pattern allows flexible configuration without complex builder logic.
//
// Example:
//
//	logger, err := log.New(
//	    log.WithConfig(log.Config{...}),
//	    log.WithDevelopmentConfig(),
//	)
type Option func(*Logger) error

// WithConfig creates a modifier that sets the logger's handler from the given config.
func WithConfig(cfg Config) Option {
	return func(l *Logger) error {
		handler, err := cfg.handler()
		if err != nil {
			return fmt.Errorf("failed to create handler: %w", err)
		}
		l.handler = handler
		l.config = cfg
		return nil
	}
}

// WithDevelopmentConfig creates a modifier that applies the development configuration.
func WithDevelopmentConfig() Option {
	return WithConfig(DevelopmentConfig)
}

// WithProductionConfig creates a modifier that applies the production configuration.
func WithProductionConfig() Option {
	return WithConfig(ProductionConfig)
}

// WithLevel creates a modifier that sets the logger's level.
// Valid levels are: debug, info, warn, error.
//
// Example:
//
//	logger, err := log.New(
//	    log.WithLevel(log.LevelDebug),
//	)
func WithLevel(level string) Option {
	return func(l *Logger) error {
		// Clone the current config (default config is already applied by the time options run)
		cfg := l.config.Clone()
		cfg.Level = level

		handler, err := cfg.handler()
		if err != nil {
			return fmt.Errorf("failed to create handler with level: %w", err)
		}

		l.handler = handler
		l.config = cfg
		return nil
	}
}

// withDefaultConfig creates a modifier that applies the default configuration.
// This is used internally when no modifiers are provided.
func withDefaultConfig() Option {
	return WithConfig(ProductionConfig)
}
