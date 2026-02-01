package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/242617/core/config"
	"github.com/242617/core/config/source"
	"github.com/242617/core/config/source/file"
)

type testItem struct {
	User struct {
		Name struct {
			First  string `env:"USER_FIRST_NAME"`
			Second string
		}
		Age     uint    `env:"USER_AGE"`
		Balance float64 `env:"USER_BALANCE" default:"10.25"`
		Active  bool    `env:"USER_ACTIVE" default:"true"`
	}
	Status  string        `yaml:"status_string" default:"ok"`
	Timeout time.Duration `env:"TIMEOUT" default:"10s"`
}

func TestDefault(t *testing.T) {
	t.Parallel()

	var cfg testItem

	engine := config.New()
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.True(t, cfg.User.Active)
	require.Equal(t, "ok", cfg.Status)
	require.Equal(t, 10.25, cfg.User.Balance)
	require.Equal(t, 10*time.Second, cfg.Timeout)
}

func TestEnvBasic(t *testing.T) {
	t.Parallel()

	envs := map[string]string{
		"USER_FIRST_NAME": "Vasily",
		"USER_ACTIVE":     "true",
		"USER_AGE":        "30",
		"USER_BALANCE":    "2.5",
		"TIMEOUT":         "20s",
	}

	for k, v := range envs {
		require.NoError(t, os.Setenv(k, v))
		defer os.Unsetenv(k)
	}

	var cfg testItem

	engine := config.New().With(source.Env())
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, "Vasily", cfg.User.Name.First)
	require.True(t, cfg.User.Active)
	require.Equal(t, uint(30), cfg.User.Age)
	require.Equal(t, 2.5, cfg.User.Balance)
	require.Equal(t, 20*time.Second, cfg.Timeout)
}

func TestYAMLBasic(t *testing.T) {
	t.Parallel()

	content := strings.Join([]string{
		"user:",
		"   name:",
		"       first: Ivan",
		"   active: true",
		"status_string: idle",
	}, "\n")

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yaml")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0666))

	var cfg testItem

	engine := config.New().With(file.YAML(filePath))
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, "Ivan", cfg.User.Name.First)
	require.True(t, cfg.User.Active)
	require.Equal(t, 10.25, cfg.User.Balance)
	require.Equal(t, "idle", cfg.Status)
}

func TestConfigEngine_With(t *testing.T) {
	t.Parallel()

	engine := config.New()
	withEngine := engine.With(source.Env())

	require.NotNil(t, withEngine)
	require.Equal(t, engine, withEngine)
}

func TestConfigEngine_Scan_NilPointer(t *testing.T) {
	t.Parallel()

	engine := config.New()
	var cfg *testItem

	err := engine.Scan(cfg)
	require.Error(t, err)
}

func TestConfigEngine_MultipleSources_OverrideOrder(t *testing.T) {
	t.Parallel()

	content := `duration: "5s"
status: "yaml_status"`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yaml")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0666))

	require.NoError(t, os.Setenv("TIMEOUT", "15s"))
	defer os.Unsetenv("TIMEOUT")

	type configType struct {
		Timeout time.Duration `env:"TIMEOUT" default:"10s" yaml:"duration"`
		Status  string        `default:"ok" yaml:"status"`
	}

	var cfg configType

	engine := config.New().With(file.YAML(filePath)).With(source.Env())
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, 15*time.Second, cfg.Timeout)
	require.Equal(t, "yaml_status", cfg.Status)
}

func TestConfigEngine_Scan_NonPointer(t *testing.T) {
	t.Parallel()

	engine := config.New()
	var cfg testItem

	err := engine.Scan(cfg)
	require.Error(t, err)
}

func TestConfigEngine_EmptyConfig(t *testing.T) {
	t.Parallel()

	engine := config.New()
	var cfg struct{}

	err := engine.Scan(&cfg)
	require.NoError(t, err)
}

func TestConfigEngine_EnvOverridesDefault(t *testing.T) {
	t.Parallel()

	type configType struct {
		Value string `env:"TEST_VALUE" default:"default_value"`
	}

	require.NoError(t, os.Setenv("TEST_VALUE", "env_value"))
	defer os.Unsetenv("TEST_VALUE")

	var cfg configType

	engine := config.New().With(source.Env())
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, "env_value", cfg.Value)
}

func TestConfigEngine_YAMLOverridesDefault(t *testing.T) {
	t.Parallel()

	type configType struct {
		Value string `default:"default_value" yaml:"value"`
	}

	content := `value: "yaml_value"`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yaml")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0666))

	var cfg configType

	engine := config.New().With(file.YAML(filePath))
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, "yaml_value", cfg.Value)
}

func TestConfigEngine_ComplexNested(t *testing.T) {
	t.Parallel()

	type nestedConfig struct {
		Level2 struct {
			Level3 struct {
				Value string `default:"deep"`
			}
		}
	}

	var cfg nestedConfig

	engine := config.New()
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, "deep", cfg.Level2.Level3.Value)
}

func TestConfigEngine_AllNumericTypes(t *testing.T) {
	t.Parallel()

	type numericConfig struct {
		Int8     int8          `default:"1"`
		Int16    int16         `default:"2"`
		Int32    int32         `default:"3"`
		Int64    int64         `default:"4"`
		Int      int           `default:"5"`
		Uint8    uint8         `default:"6"`
		Uint16   uint16        `default:"7"`
		Uint32   uint32        `default:"8"`
		Uint64   uint64        `default:"9"`
		Uint     uint          `default:"10"`
		Float32  float32       `default:"1.5"`
		Float64  float64       `default:"2.5"`
		Duration time.Duration `default:"30s"`
	}

	var cfg numericConfig

	engine := config.New()
	err := engine.Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, int8(1), cfg.Int8)
	require.Equal(t, int16(2), cfg.Int16)
	require.Equal(t, int32(3), cfg.Int32)
	require.Equal(t, int64(4), cfg.Int64)
	require.Equal(t, 5, cfg.Int)
	require.Equal(t, uint8(6), cfg.Uint8)
	require.Equal(t, uint16(7), cfg.Uint16)
	require.Equal(t, uint32(8), cfg.Uint32)
	require.Equal(t, uint64(9), cfg.Uint64)
	require.Equal(t, uint(10), cfg.Uint)
	require.Equal(t, float32(1.5), cfg.Float32)
	require.Equal(t, 2.5, cfg.Float64)
	require.Equal(t, 30*time.Second, cfg.Duration)
}

func TestConfigEngine_BoolVariations(t *testing.T) {
	for _, tc := range []struct {
		name     string
		envValue string
		expected bool
	}{
		{"true_lowercase", "true", true},
		{"true_uppercase", "TRUE", true},
		{"true_mixedcase", "True", true},
		{"false_lowercase", "false", false},
		{"false_uppercase", "FALSE", false},
		{"false_mixedcase", "False", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			type boolConfig struct {
				Value bool `env:"TEST_BOOL"`
			}

			require.NoError(t, os.Setenv("TEST_BOOL", tc.envValue))
			defer os.Unsetenv("TEST_BOOL")

			var cfg boolConfig
			engine := config.New().With(source.Env())
			err := engine.Scan(&cfg)
			require.NoError(t, err)
			require.Equal(t, tc.expected, cfg.Value)
		})
	}
}

func TestConfigEngine_SourceOrder(t *testing.T) {
	t.Parallel()

	type orderedConfig struct {
		Value string `env:"TEST_VALUE" default:"default" yaml:"value"`
	}

	t.Run("default_then_yaml_then_env", func(t *testing.T) {
		t.Parallel()

		content := `value: "yaml_value"`
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "config.yaml")
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0666))

		require.NoError(t, os.Setenv("TEST_VALUE", "env_value"))
		defer os.Unsetenv("TEST_VALUE")

		var cfg orderedConfig
		engine := config.New().With(file.YAML(filePath)).With(source.Env())
		err := engine.Scan(&cfg)
		require.NoError(t, err)
		require.Equal(t, "env_value", cfg.Value)
	})

	t.Run("default_then_env_then_yaml", func(t *testing.T) {
		t.Parallel()

		content := `value: "yaml_value"`
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "config.yaml")
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0666))

		require.NoError(t, os.Setenv("TEST_VALUE", "env_value"))
		defer os.Unsetenv("TEST_VALUE")

		var cfg orderedConfig
		engine := config.New().With(source.Env()).With(file.YAML(filePath))
		err := engine.Scan(&cfg)
		require.NoError(t, err)
		require.Equal(t, "yaml_value", cfg.Value)
	})
}
