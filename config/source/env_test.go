package source_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/242617/core/config/source"
)

type envTestConfig struct {
	StringVal  string        `env:"TEST_STRING"`
	IntVal     int           `env:"TEST_INT"`
	UintVal    uint          `env:"TEST_UINT"`
	FloatVal   float64       `env:"TEST_FLOAT"`
	BoolVal    bool          `env:"TEST_BOOL"`
	Duration   time.Duration `env:"TEST_DURATION"`
	Int8Val    int8          `env:"TEST_INT8"`
	Int16Val   int16         `env:"TEST_INT16"`
	Int32Val   int32         `env:"TEST_INT32"`
	Int64Val   int64         `env:"TEST_INT64"`
	Uint8Val   uint8         `env:"TEST_UINT8"`
	Uint16Val  uint16        `env:"TEST_UINT16"`
	Uint32Val  uint32        `env:"TEST_UINT32"`
	Uint64Val  uint64        `env:"TEST_UINT64"`
	Float32Val float32       `env:"TEST_FLOAT32"`
	Nested     nestedEnv     `env:""`
	NoTag      string
}

type nestedEnv struct {
	ChildField string `env:"TEST_CHILD"`
}

func TestEnv_Success(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "string",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "test_string"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value string `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "test_string", cfg.Value)
			},
		},
		{
			name: "int",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "123"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value int `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 123, cfg.Value)
			},
		},
		{
			name: "int8",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "127"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value int8 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int8(127), cfg.Value)
			},
		},
		{
			name: "int16",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "32767"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value int16 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int16(32767), cfg.Value)
			},
		},
		{
			name: "int32",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "2147483647"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value int32 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int32(2147483647), cfg.Value)
			},
		},
		{
			name: "int64",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "9223372036854775807"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value int64 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int64(9223372036854775807), cfg.Value)
			},
		},
		{
			name: "uint",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "456"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value uint `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint(456), cfg.Value)
			},
		},
		{
			name: "uint8",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "255"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value uint8 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint8(255), cfg.Value)
			},
		},
		{
			name: "uint16",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "65535"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value uint16 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint16(65535), cfg.Value)
			},
		},
		{
			name: "uint32",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "4294967295"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value uint32 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint32(4294967295), cfg.Value)
			},
		},
		{
			name: "uint64",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "18446744073709551615"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value uint64 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint64(18446744073709551615), cfg.Value)
			},
		},
		{
			name: "float32",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "1.5"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value float32 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, float32(1.5), cfg.Value)
			},
		},
		{
			name: "float64",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "2.75"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value float64 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 2.75, cfg.Value)
			},
		},
		{
			name: "bool_true",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "true"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value bool `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.True(t, cfg.Value)
			},
		},
		{
			name: "bool_false",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "false"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value bool `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.False(t, cfg.Value)
			},
		},
		{
			name: "bool_case_insensitive",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "TRUE"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value bool `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.True(t, cfg.Value)
			},
		},
		{
			name: "duration",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "30s"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value time.Duration `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 30*time.Second, cfg.Value)
			},
		},
		{
			name: "nested_struct",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_CHILD", "nested_value"))
				defer os.Unsetenv("TEST_CHILD")

				var cfg struct {
					Nested struct {
						Value string `env:"TEST_CHILD"`
					}
				}
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "nested_value", cfg.Nested.Value)
			},
		},
		{
			name: "no_env_tag",
			test: func(t *testing.T) {
				var cfg struct {
					Value string
				}
				cfg.Value = "initial"
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "initial", cfg.Value)
			},
		},
		{
			name: "empty_env_value",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", ""))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value string `env:"TEST_VAL"`
				}
				cfg.Value = "initial"
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "initial", cfg.Value)
			},
		},
		{
			name: "env_not_set",
			test: func(t *testing.T) {
				var cfg struct {
					Value string `env:"NONEXISTENT_VAR"`
				}
				cfg.Value = "initial"
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "initial", cfg.Value)
			},
		},
		{
			name: "complete_config",
			test: func(t *testing.T) {
				envs := map[string]string{
					"TEST_STRING":   "world",
					"TEST_INT":      "100",
					"TEST_UINT":     "200",
					"TEST_FLOAT":    "6.28",
					"TEST_BOOL":     "false",
					"TEST_DURATION": "15m",
					"TEST_INT8":     "5",
					"TEST_INT16":    "6",
					"TEST_INT32":    "7",
					"TEST_INT64":    "8",
					"TEST_UINT8":    "9",
					"TEST_UINT16":   "10",
					"TEST_UINT32":   "11",
					"TEST_UINT64":   "12",
					"TEST_FLOAT32":  "3.5",
					"TEST_CHILD":    "nested_env",
				}

				for k, v := range envs {
					require.NoError(t, os.Setenv(k, v))
					defer os.Unsetenv(k)
				}

				var cfg envTestConfig
				err := source.Env().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "world", cfg.StringVal)
				require.Equal(t, 100, cfg.IntVal)
				require.Equal(t, uint(200), cfg.UintVal)
				require.Equal(t, 6.28, cfg.FloatVal)
				require.False(t, cfg.BoolVal)
				require.Equal(t, 15*time.Minute, cfg.Duration)
				require.Equal(t, int8(5), cfg.Int8Val)
				require.Equal(t, int16(6), cfg.Int16Val)
				require.Equal(t, int32(7), cfg.Int32Val)
				require.Equal(t, int64(8), cfg.Int64Val)
				require.Equal(t, uint8(9), cfg.Uint8Val)
				require.Equal(t, uint16(10), cfg.Uint16Val)
				require.Equal(t, uint32(11), cfg.Uint32Val)
				require.Equal(t, uint64(12), cfg.Uint64Val)
				require.Equal(t, float32(3.5), cfg.Float32Val)
				require.Equal(t, "nested_env", cfg.Nested.ChildField)
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

func TestEnv_Error(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "nil_pointer",
			test: func(t *testing.T) {
				var cfg *struct{}
				err := source.Env().Scan(cfg)
				require.Error(t, err)
				require.Contains(t, err.Error(), "unexpected kind")
			},
		},
		{
			name: "non_pointer",
			test: func(t *testing.T) {
				var cfg struct{}
				err := source.Env().Scan(cfg)
				require.Error(t, err)
				require.Contains(t, err.Error(), "unexpected kind")
			},
		},
		{
			name: "invalid_int",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "not_a_number"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value int `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.Error(t, err)
			},
		},
		{
			name: "invalid_float",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "not_a_float"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value float64 `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.Error(t, err)
			},
		},
		{
			name: "invalid_duration",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "invalid_duration"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value time.Duration `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.Error(t, err)
			},
		},
		{
			name: "unsupported_type",
			test: func(t *testing.T) {
				require.NoError(t, os.Setenv("TEST_VAL", "test"))
				defer os.Unsetenv("TEST_VAL")

				var cfg struct {
					Value []string `env:"TEST_VAL"`
				}
				err := source.Env().Scan(&cfg)
				require.Error(t, err)
				require.Contains(t, err.Error(), "unsupported type")
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
