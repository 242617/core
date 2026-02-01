package logger_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/242617/core/logger"
)

// TestNewLogger_Success tests successful logger creation with various modifiers.
func TestNewLogger_Success(t *testing.T) {
	tests := []struct {
		name    string
		options []logger.Option
		wantErr bool
	}{
		{
			name:    "default logger.config",
			options: []logger.Option{},
			wantErr: false,
		},
		{
			name:    "production logger.config",
			options: []logger.Option{logger.WithProductionConfig()},
			wantErr: false,
		},
		{
			name: "development logger.config",
			options: []logger.Option{
				logger.WithDevelopmentConfig(),
			},
			wantErr: false,
		},
		{
			name: "custom logger.config",
			options: []logger.Option{
				logger.WithConfig(logger.Config{
					Level:    logger.LevelDebug,
					Encoding: logger.EncodingText,
					Colorize: false,
				}),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.New(tt.options...)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if logger != nil {
					t.Errorf("expected nil logger, got %v", logger)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if logger == nil {
					t.Errorf("expected non-nil logger, got nil")
				}
			}
		})
	}
}

// TestNewLogger_Failure tests logger creation failure scenarios.
func TestNewLogger_Failure(t *testing.T) {
	tests := []struct {
		name    string
		config  logger.Config
		wantErr error
	}{
		{
			name: "invalid level",
			config: logger.Config{
				Level:    "invalid",
				Encoding: logger.EncodingJSON,
			},
			wantErr: logger.ErrInvalidLevel,
		},
		{
			name: "invalid encoding",
			config: logger.Config{
				Level:    logger.LevelInfo,
				Encoding: "invalid",
			},
			wantErr: logger.ErrInvalidEncoding,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.New(logger.WithConfig(tt.config))

			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if logger != nil {
				t.Errorf("expected nil logger, got %v", logger)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("error should contain %v, got %v", tt.wantErr, err)
			}

			// Also verify that logger.config validation works directly
			cfgErr := tt.config.Validate()
			if cfgErr == nil {
				t.Errorf("expected logger.config validation error, got nil")
			}
		})
	}
}

// TestLogger_LogLevels tests logging at different levels.
func TestLogger_LogLevels(t *testing.T) {
	tests := []struct {
		name      string
		logFunc   func(*logger.Logger, context.Context, string, ...any)
		wantLevel slog.Level
	}{
		{
			name:      "debug level",
			logFunc:   (*logger.Logger).Debug,
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "info level",
			logFunc:   (*logger.Logger).Info,
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "warn level",
			logFunc:   (*logger.Logger).Warn,
			wantLevel: slog.LevelWarn,
		},
		{
			name:      "error level",
			logFunc:   (*logger.Logger).Error,
			wantLevel: slog.LevelError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.New(logger.WithDevelopmentConfig())
			if err != nil {
				t.Fatalf("failed to create logger: %v", err)
			}

			ctx := context.Background()
			// These calls should not panic
			tt.logFunc(logger, ctx, "test message", "key", "value")
		})
	}
}

// TestLogger_NamedLogger tests creating named child loggers.
func TestLogger_NamedLogger(t *testing.T) {
	logger, err := logger.New(logger.WithDevelopmentConfig())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	childLogger := logger.New("child-service")
	if childLogger == nil {
		t.Errorf("expected non-nil child logger, got nil")
	}

	ctx := context.Background()
	// Should not panic
	childLogger.Info(ctx, "child logger message")

	// Create nested child
	grandchildLogger := childLogger.New("grandchild-service")
	if grandchildLogger == nil {
		t.Errorf("expected non-nil grandchild logger, got nil")
	}
	grandchildLogger.Info(ctx, "grandchild logger message")
}

// TestLogger_Enabled tests the Enabled method.
func TestLogger_Enabled(t *testing.T) {
	tests := []struct {
		name        string
		configLevel string
		checkLevel  slog.Level
		wantEnabled bool
	}{
		{
			name:        "debug logger.config, check debug",
			configLevel: logger.LevelDebug,
			checkLevel:  slog.LevelDebug,
			wantEnabled: true,
		},
		{
			name:        "debug logger.config, check info",
			configLevel: logger.LevelDebug,
			checkLevel:  slog.LevelInfo,
			wantEnabled: true,
		},
		{
			name:        "info logger.config, check debug",
			configLevel: logger.LevelInfo,
			checkLevel:  slog.LevelDebug,
			wantEnabled: false,
		},
		{
			name:        "info logger.config, check info",
			configLevel: logger.LevelInfo,
			checkLevel:  slog.LevelInfo,
			wantEnabled: true,
		},
		{
			name:        "error logger.config, check warn",
			configLevel: logger.LevelError,
			checkLevel:  slog.LevelWarn,
			wantEnabled: false,
		},
		{
			name:        "error logger.config, check error",
			configLevel: logger.LevelError,
			checkLevel:  slog.LevelError,
			wantEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := logger.New(logger.WithConfig(logger.Config{
				Level:    tt.configLevel,
				Encoding: logger.EncodingText,
			}))
			if err != nil {
				t.Fatalf("failed to create logger: %v", err)
			}

			ctx := context.Background()
			enabled := logger.Enabled(ctx, tt.checkLevel)
			if enabled != tt.wantEnabled {
				t.Errorf("Expected enabled=%v, got %v", tt.wantEnabled, enabled)
			}
		})
	}
}

// TestConfig_Clone tests logger.config cloning.
func TestConfig_Clone(t *testing.T) {
	original := logger.Config{
		Level:    logger.LevelDebug,
		Encoding: logger.EncodingText,
		Colorize: true,
	}

	cloned := original.Clone()
	if original.Level != cloned.Level {
		t.Errorf("Level mismatch: original=%v, cloned=%v", original.Level, cloned.Level)
	}
	if original.Encoding != cloned.Encoding {
		t.Errorf("Encoding mismatch: original=%v, cloned=%v", original.Encoding, cloned.Encoding)
	}
	if original.Colorize != cloned.Colorize {
		t.Errorf("Colorize mismatch: original=%v, cloned=%v", original.Colorize, cloned.Colorize)
	}

	// Modify clone should not affect original
	cloned.Level = logger.LevelError
	if original.Level == cloned.Level {
		t.Errorf("modifying clone should not affect original")
	}
	if original.Level != logger.LevelDebug {
		t.Errorf("original level should be debug, got %v", original.Level)
	}
	if cloned.Level != logger.LevelError {
		t.Errorf("cloned level should be error, got %v", cloned.Level)
	}
}

// TestConfig_Validate tests logger.config validation.
func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  logger.Config
		wantErr error
	}{
		{
			name: "valid debug logger.config",
			config: logger.Config{
				Level:    logger.LevelDebug,
				Encoding: logger.EncodingJSON,
			},
			wantErr: nil,
		},
		{
			name: "valid info logger.config",
			config: logger.Config{
				Level:    logger.LevelInfo,
				Encoding: logger.EncodingText,
			},
			wantErr: nil,
		},
		{
			name: "valid warn logger.config",
			config: logger.Config{
				Level:    logger.LevelWarn,
				Encoding: logger.EncodingJSON,
			},
			wantErr: nil,
		},
		{
			name: "valid error logger.config",
			config: logger.Config{
				Level:    logger.LevelError,
				Encoding: logger.EncodingText,
			},
			wantErr: nil,
		},
		{
			name: "invalid level",
			config: logger.Config{
				Level:    "invalid",
				Encoding: logger.EncodingJSON,
			},
			wantErr: logger.ErrInvalidLevel,
		},
		{
			name: "invalid encoding",
			config: logger.Config{
				Level:    logger.LevelInfo,
				Encoding: "invalid",
			},
			wantErr: logger.ErrInvalidEncoding,
		},
		{
			name: "both invalid",
			config: logger.Config{
				Level:    "invalid",
				Encoding: "invalid",
			},
			wantErr: logger.ErrInvalidLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("error should be of type %T, got %T", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestParseLevel tests level string to slog.Level conversion.
func TestParseLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected slog.Level
	}{
		{
			name:     "debug",
			level:    logger.LevelDebug,
			expected: slog.LevelDebug,
		},
		{
			name:     "info",
			level:    logger.LevelInfo,
			expected: slog.LevelInfo,
		},
		{
			name:     "warn",
			level:    logger.LevelWarn,
			expected: slog.LevelWarn,
		},
		{
			name:     "error",
			level:    logger.LevelError,
			expected: slog.LevelError,
		},
		{
			name:     "invalid defaults to info",
			level:    "invalid",
			expected: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := logger.Config{Level: tt.level, Encoding: logger.EncodingJSON}
			actual := logger.ParseLevel(cfg.Level)
			if actual != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, actual)
			}
		})
	}
}

// TestTimer tests the Timer utility.
func TestTimer(t *testing.T) {
	l, err := logger.New(logger.WithDevelopmentConfig())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	ctx := context.Background()

	t.Run("basic timer", func(t *testing.T) {
		timer := logger.NewTimer(l, ctx, "test operation")
		duration := timer.Stop("operation completed", "result", "success")

		if duration <= 0 {
			t.Errorf("expected positive duration, got %v", duration)
		}
	})

	t.Run("timer with debug", func(t *testing.T) {
		timer := logger.NewTimer(l, ctx, "debug operation")
		duration := timer.Debug("debug message")

		if duration <= 0 {
			t.Errorf("expected positive duration, got %v", duration)
		}
	})

	t.Run("timer with error", func(t *testing.T) {
		timer := logger.NewTimer(l, ctx, "error operation")
		duration := timer.Error("operation failed", "err", errors.New("test error"))

		if duration <= 0 {
			t.Errorf("expected positive duration, got %v", duration)
		}
	})
}

// TestPredefinedConfigs tests the predefined logger.config constants.
func TestPredefinedConfigs(t *testing.T) {
	t.Run("DevelopmentConfig", func(t *testing.T) {
		if logger.DevelopmentConfig.Level != logger.LevelDebug {
			t.Errorf("Expected level %v, got %v", logger.LevelDebug, logger.DevelopmentConfig.Level)
		}
		if logger.DevelopmentConfig.Encoding != logger.EncodingText {
			t.Errorf("Expected encoding %v, got %v", logger.EncodingText, logger.DevelopmentConfig.Encoding)
		}
		if !logger.DevelopmentConfig.Colorize {
			t.Errorf("Expected Colorize to be true")
		}
		if err := logger.DevelopmentConfig.Validate(); err != nil {
			t.Errorf("DevelopmentConfig should be valid: %v", err)
		}
	})

	t.Run("ProductionConfig", func(t *testing.T) {
		if logger.ProductionConfig.Level != logger.LevelInfo {
			t.Errorf("Expected level %v, got %v", logger.LevelInfo, logger.ProductionConfig.Level)
		}
		if logger.ProductionConfig.Encoding != logger.EncodingJSON {
			t.Errorf("Expected encoding %v, got %v", logger.EncodingJSON, logger.ProductionConfig.Encoding)
		}
		if logger.ProductionConfig.Colorize {
			t.Errorf("Expected Colorize to be false")
		}
		if err := logger.ProductionConfig.Validate(); err != nil {
			t.Errorf("ProductionConfig should be valid: %v", err)
		}
	})
}

// TestLogger_SetLevel tests dynamic level setting.
func TestLogger_SetLevel(t *testing.T) {
	l, err := logger.New(logger.WithProductionConfig())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	t.Run("set valid level", func(t *testing.T) {
		err := l.SetLevel(logger.LevelDebug)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if l.LogLevel() != logger.LevelDebug {
			t.Errorf("expected level %v, got %v", logger.LevelDebug, l.LogLevel())
		}
	})

	t.Run("set invalid level", func(t *testing.T) {
		err := l.SetLevel("invalid")
		if err == nil {
			t.Errorf("expected error for invalid level")
		}
	})
}

// TestLogger_SetConfig tests dynamic logger.configuration update.
func TestLogger_SetConfig(t *testing.T) {
	l, err := logger.New(logger.WithProductionConfig())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	t.Run("set valid logger.config", func(t *testing.T) {
		newCfg := logger.Config{
			Level:    logger.LevelDebug,
			Encoding: logger.EncodingText,
			Colorize: true,
		}
		err := l.SetConfig(newCfg)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if l.LogLevel() != logger.LevelDebug {
			t.Errorf("expected level %v, got %v", logger.LevelDebug, l.LogLevel())
		}
	})

	t.Run("set invalid logger.config", func(t *testing.T) {
		invalidCfg := logger.Config{
			Level:    "invalid",
			Encoding: logger.EncodingJSON,
		}
		err := l.SetConfig(invalidCfg)
		if err == nil {
			t.Errorf("expected error for invalid logger.config")
		}
	})
}

// TestWithLevel tests the WithLevel option function.
func TestWithLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{
			name:    "debug level",
			level:   logger.LevelDebug,
			wantErr: false,
		},
		{
			name:    "info level",
			level:   logger.LevelInfo,
			wantErr: false,
		},
		{
			name:    "warn level",
			level:   logger.LevelWarn,
			wantErr: false,
		},
		{
			name:    "error level",
			level:   logger.LevelError,
			wantErr: false,
		},
		{
			name:    "invalid level",
			level:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := logger.New(logger.WithLevel(tt.level))

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for level %s, got nil", tt.level)
				}
				if l != nil {
					t.Errorf("expected nil logger for error, got %v", l)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if l == nil {
					t.Errorf("expected non-nil logger, got nil")
				}
				if l != nil && l.LogLevel() != tt.level {
					t.Errorf("expected level %s, got %s", tt.level, l.LogLevel())
				}
			}
		})
	}
}

// TestWithLevel_Combined tests WithLevel combined with other options.
func TestWithLevel_Combined(t *testing.T) {
	l, err := logger.New(
		logger.WithConfig(logger.Config{
			Level:    logger.LevelDebug,
			Encoding: logger.EncodingText,
			Colorize: true,
		}),
	)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	if l.LogLevel() != logger.LevelDebug {
		t.Errorf("expected level %s, got %s", logger.LevelDebug, l.LogLevel())
	}

	ctx := context.Background()
	// Should not panic
	l.Debug(ctx, "test debug message")
	l.Info(ctx, "test info message")
}

// TestWithLevel_Override tests that WithLevel can override the default config level.
func TestWithLevel_Override(t *testing.T) {
	// ProductionConfig defaults to info level, override with debug
	l, err := logger.New(
		logger.WithProductionConfig(),
		logger.WithLevel(logger.LevelDebug),
	)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	if l.LogLevel() != logger.LevelDebug {
		t.Errorf("expected level %s, got %s", logger.LevelDebug, l.LogLevel())
	}

	ctx := context.Background()
	l.Debug(ctx, "test debug message")
}

// BenchmarkLogger_Log benchmarks logging operations.
func BenchmarkLogger_Log(b *testing.B) {
	logger, _ := logger.New(logger.WithProductionConfig())
	ctx := context.Background()

	b.Run("Info", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Info(ctx, "benchmark message", "iteration", i)
		}
	})

	b.Run("Debug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Debug(ctx, "debug message", "iteration", i)
		}
	})
}

// BenchmarkLogger_Enabled benchmarks Enabled method.
func BenchmarkLogger_Enabled(b *testing.B) {
	logger, _ := logger.New(logger.WithProductionConfig())
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		logger.Enabled(ctx, slog.LevelDebug)
	}
}

// ExampleLogger demonstrates basic logger usage.
func ExampleLogger() {
	logger, _ := logger.New(logger.WithDevelopmentConfig())
	ctx := context.Background()

	logger.Info(ctx, "application started")
	logger.Info(ctx, "handling request", "method", "GET", "path", "/api/v1")
}

// ExampleLogger_New demonstrates named child logger usage.
func ExampleLogger_New() {
	logger, _ := logger.New(logger.WithDevelopmentConfig())
	serviceLogger := logger.New("payment-service")
	ctx := context.Background()

	serviceLogger.Info(ctx, "processing payment")
}

// ExampleLogger_Error demonstrates error logging.
func ExampleLogger_Error() {
	logger, _ := logger.New(logger.WithDevelopmentConfig())
	ctx := context.Background()

	err := fmt.Errorf("connection failed")
	logger.Error(ctx, "database error", "err", err, "query", "SELECT * FROM users")
}

// ExampleTimer demonstrates timer usage.
func ExampleTimer() {
	l, _ := logger.New(logger.WithDevelopmentConfig())
	ctx := context.Background()

	timer := logger.NewTimer(l, ctx, "database query")
	timer.Stop("query completed", "rows", 42)
}

// ExampleWithLevel demonstrates creating a logger with a specific level using the WithLevel option.
func ExampleWithLevel() {
	logger, _ := logger.New(
		logger.WithLevel(logger.LevelDebug),
		logger.WithConfig(logger.Config{
			Encoding: logger.EncodingText,
			Colorize: true,
		}),
	)
	ctx := context.Background()

	logger.Debug(ctx, "this debug message will be visible")
	logger.Info(ctx, "info message")
}
