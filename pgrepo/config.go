package pgrepo

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Config holds the database configuration for a single database instance.
// It can represent either a master or a replica database.
type Config struct {
	// Connection settings
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Schema   string `yaml:"schema"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSL      bool   `yaml:"ssl" default:"false"`

	// Pool settings
	ConnMaxLifeTime time.Duration `yaml:"conn_max_life_time" default:"1h"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" default:"30m"`
	MinConns        int           `yaml:"min_conns" default:"2"`
	MaxConns        int           `yaml:"max_conns" default:"25"`

	// Graceful shutdown settings
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" default:"30s"`

	// Replicas - for master configuration only
	Replicas []Config `yaml:"replicas"`
}

// Validate checks if the configuration has all required fields and valid values.
func (cfg *Config) Validate() error {
	switch {
	case cfg.Host == "":
		return errors.New("host is required")
	case cfg.Port <= 0 || cfg.Port > 65535:
		return errors.New("port must be between 1 and 65535")
	case cfg.Schema == "":
		return errors.New("schema is required")
	case cfg.User == "":
		return errors.New("user is required")
	case cfg.Password == "":
		return errors.New("password is required")
	case cfg.Name == "":
		return errors.New("database name is required")
	case cfg.MinConns < 0:
		return errors.New("min_conns must be non-negative")
	case cfg.MaxConns <= 0:
		return errors.New("max_conns must be positive")
	case cfg.MinConns > cfg.MaxConns:
		return errors.New("min_conns cannot be greater than max_conns")
	case cfg.ConnMaxLifeTime <= 0:
		return errors.New("conn_max_life_time must be positive")
	case cfg.ConnMaxIdleTime <= 0:
		return errors.New("conn_max_idle_time must be positive")
	case cfg.ShutdownTimeout <= 0:
		return errors.New("shutdown_timeout must be positive")
	}

	// Validate replicas recursively
	for i := range cfg.Replicas {
		if err := cfg.Replicas[i].Validate(); err != nil {
			return errors.Wrapf(err, "replica[%d]", i)
		}
	}

	return nil
}

// String returns the PostgreSQL DSN (Data Source Name) for this configuration.
// The password is included in the returned string (redact before logging).
func (cfg *Config) String() string {
	format := "postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s"
	if cfg.SSL {
		format = "postgres://%s:%s@%s:%d/%s?sslmode=require&search_path=%s"
	}
	return fmt.Sprintf(
		format,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Schema,
	)
}

// RedactedDSN returns the PostgreSQL DSN with the password redacted for logging.
func (cfg *Config) RedactedDSN() string {
	format := "postgres://%s:***@%s:%d/%s?sslmode=%s&search_path=%s"
	sslMode := "disable"
	if cfg.SSL {
		sslMode = "require"
	}
	return fmt.Sprintf(
		format,
		cfg.User,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		sslMode,
		cfg.Schema,
	)
}
