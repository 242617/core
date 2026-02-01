package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"

	EncodingText = "text"
	EncodingJSON = "json"
)

var (
	// DevelopmentConfig provides sensible defaults for development environments
	// with debug level, text encoding, and colored output.
	DevelopmentConfig = Config{
		Level:    LevelDebug,
		Encoding: EncodingText,
		Colorize: true,
	}

	// ProductionConfig provides sensible defaults for production environments
	// with info level, JSON encoding, and no colors.
	ProductionConfig = Config{
		Level:    LevelInfo,
		Encoding: EncodingJSON,
	}
)

// Config represents the logger configuration.
type Config struct {
	// Level is the minimum log level to capture.
	// Accepted values: debug, info, warn, error.
	// Default: info.
	Level string `yaml:"level" default:"info"`

	// Encoding specifies the log output format.
	// Accepted values: json, text.
	// Default: json.
	Encoding string `yaml:"encoding" default:"json"`

	// Colorize enables colored output for text encoding.
	// Has no effect for JSON encoding.
	Colorize bool `yaml:"colorize"`
}

// Clone creates a deep copy of the configuration.
func (cfg Config) Clone() Config {
	return Config{
		Level:    cfg.Level,
		Encoding: cfg.Encoding,
		Colorize: cfg.Colorize,
	}
}

// Validate checks if the configuration is valid.
// Returns an error if:
//   - level is not one of: debug, info, warn, error
//   - encoding is not one of: json, text
func (cfg Config) Validate() error {
	switch cfg.Level {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		// valid
	default:
		return ErrInvalidLevel
	}

	switch cfg.Encoding {
	case EncodingText, EncodingJSON:
		// valid
	default:
		return ErrInvalidEncoding
	}

	return nil
}

// ParseLevel converts the string level to slog.Level.
func ParseLevel(level string) slog.Level {
	switch level {
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// handler creates the slog.Handler based on the configuration.
func (cfg Config) handler() (slog.Handler, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	level := ParseLevel(cfg.Level)

	switch cfg.Encoding {
	case EncodingText:
		return &contextHandler{
			Handler: tint.NewHandler(
				os.Stderr,
				&tint.Options{
					Level: level,
					ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
						// Special handling for errors to show stack traces in dev
						if err, ok := a.Value.Any().(error); ok && a.Key == "err" {
							return tint.Err(err)
						}
						return a
					},
					TimeFormat: "15:04:05.99",
					NoColor:    !cfg.Colorize,
				},
			),
		}, nil

	case EncodingJSON:
		return &contextHandler{
			Handler: slog.NewJSONHandler(
				os.Stderr,
				&slog.HandlerOptions{
					Level: level,
				},
			),
		}, nil

	default:
		// This should never happen due to Validate()
		return nil, ErrInvalidEncoding
	}
}
