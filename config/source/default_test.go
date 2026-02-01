package source_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/242617/core/config/source"
)

type testConfig struct {
	StringVal  string        `default:"hello"`
	IntVal     int           `default:"42"`
	UintVal    uint          `default:"99"`
	FloatVal   float64       `default:"3.14"`
	BoolVal    bool          `default:"true"`
	Duration   time.Duration `default:"5s"`
	Int8Val    int8          `default:"1"`
	Int16Val   int16         `default:"2"`
	Int32Val   int32         `default:"3"`
	Int64Val   int64         `default:"4"`
	Uint8Val   uint8         `default:"5"`
	Uint16Val  uint16        `default:"6"`
	Uint32Val  uint32        `default:"7"`
	Uint64Val  uint64        `default:"8"`
	Float32Val float32       `default:"2.5"`
	Nested     nestedConfig  `default:""`
	NoDefault  string
}

type nestedConfig struct {
	ChildField string `default:"child"`
}

func TestDefault_Success(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "string",
			test: func(t *testing.T) {
				var cfg struct {
					Value string `default:"test"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "test", cfg.Value)
			},
		},
		{
			name: "int",
			test: func(t *testing.T) {
				var cfg struct {
					Value int `default:"123"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 123, cfg.Value)
			},
		},
		{
			name: "int8",
			test: func(t *testing.T) {
				var cfg struct {
					Value int8 `default:"127"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int8(127), cfg.Value)
			},
		},
		{
			name: "int16",
			test: func(t *testing.T) {
				var cfg struct {
					Value int16 `default:"32767"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int16(32767), cfg.Value)
			},
		},
		{
			name: "int32",
			test: func(t *testing.T) {
				var cfg struct {
					Value int32 `default:"2147483647"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int32(2147483647), cfg.Value)
			},
		},
		{
			name: "int64",
			test: func(t *testing.T) {
				var cfg struct {
					Value int64 `default:"9223372036854775807"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, int64(9223372036854775807), cfg.Value)
			},
		},
		{
			name: "uint",
			test: func(t *testing.T) {
				var cfg struct {
					Value uint `default:"456"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint(456), cfg.Value)
			},
		},
		{
			name: "uint8",
			test: func(t *testing.T) {
				var cfg struct {
					Value uint8 `default:"255"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint8(255), cfg.Value)
			},
		},
		{
			name: "uint16",
			test: func(t *testing.T) {
				var cfg struct {
					Value uint16 `default:"65535"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint16(65535), cfg.Value)
			},
		},
		{
			name: "uint32",
			test: func(t *testing.T) {
				var cfg struct {
					Value uint32 `default:"4294967295"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint32(4294967295), cfg.Value)
			},
		},
		{
			name: "uint64",
			test: func(t *testing.T) {
				var cfg struct {
					Value uint64 `default:"18446744073709551615"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint64(18446744073709551615), cfg.Value)
			},
		},
		{
			name: "float32",
			test: func(t *testing.T) {
				var cfg struct {
					Value float32 `default:"1.5"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, float32(1.5), cfg.Value)
			},
		},
		{
			name: "float64",
			test: func(t *testing.T) {
				var cfg struct {
					Value float64 `default:"2.75"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 2.75, cfg.Value)
			},
		},
		{
			name: "bool_true",
			test: func(t *testing.T) {
				var cfg struct {
					Value bool `default:"true"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.True(t, cfg.Value)
			},
		},
		{
			name: "bool_false",
			test: func(t *testing.T) {
				var cfg struct {
					Value bool `default:"false"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.False(t, cfg.Value)
			},
		},
		{
			name: "bool_case_insensitive",
			test: func(t *testing.T) {
				var cfg struct {
					Value bool `default:"TRUE"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.True(t, cfg.Value)
			},
		},
		{
			name: "duration",
			test: func(t *testing.T) {
				var cfg struct {
					Value time.Duration `default:"10s"`
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 10*time.Second, cfg.Value)
			},
		},
		{
			name: "nested_struct",
			test: func(t *testing.T) {
				var cfg struct {
					Nested struct {
						Value string `default:"nested_value"`
					}
				}
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "nested_value", cfg.Nested.Value)
			},
		},
		{
			name: "no_default_tag",
			test: func(t *testing.T) {
				var cfg struct {
					Value string
				}
				cfg.Value = "initial"
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "initial", cfg.Value)
			},
		},
		{
			name: "empty_default_tag",
			test: func(t *testing.T) {
				var cfg struct {
					Value string `default:""`
				}
				cfg.Value = "initial"
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "initial", cfg.Value)
			},
		},
		{
			name: "complete_config",
			test: func(t *testing.T) {
				var cfg testConfig
				err := source.Default().Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "hello", cfg.StringVal)
				require.Equal(t, 42, cfg.IntVal)
				require.Equal(t, uint(99), cfg.UintVal)
				require.Equal(t, 3.14, cfg.FloatVal)
				require.True(t, cfg.BoolVal)
				require.Equal(t, 5*time.Second, cfg.Duration)
				require.Equal(t, int8(1), cfg.Int8Val)
				require.Equal(t, int16(2), cfg.Int16Val)
				require.Equal(t, int32(3), cfg.Int32Val)
				require.Equal(t, int64(4), cfg.Int64Val)
				require.Equal(t, uint8(5), cfg.Uint8Val)
				require.Equal(t, uint16(6), cfg.Uint16Val)
				require.Equal(t, uint32(7), cfg.Uint32Val)
				require.Equal(t, uint64(8), cfg.Uint64Val)
				require.Equal(t, float32(2.5), cfg.Float32Val)
				require.Equal(t, "child", cfg.Nested.ChildField)
				require.Empty(t, cfg.NoDefault)
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

func TestDefault_Error(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "nil_pointer",
			test: func(t *testing.T) {
				var cfg *struct{}
				err := source.Default().Scan(cfg)
				require.Error(t, err)
				require.Contains(t, err.Error(), "unexpected kind")
			},
		},
		{
			name: "non_pointer",
			test: func(t *testing.T) {
				var cfg struct{}
				err := source.Default().Scan(cfg)
				require.Error(t, err)
				require.Contains(t, err.Error(), "unexpected kind")
			},
		},
		{
			name: "invalid_int",
			test: func(t *testing.T) {
				var cfg struct {
					Value int `default:"not_a_number"`
				}
				err := source.Default().Scan(&cfg)
				require.Error(t, err)
			},
		},
		{
			name: "invalid_float",
			test: func(t *testing.T) {
				var cfg struct {
					Value float64 `default:"not_a_float"`
				}
				err := source.Default().Scan(&cfg)
				require.Error(t, err)
			},
		},
		{
			name: "invalid_duration",
			test: func(t *testing.T) {
				var cfg struct {
					Value time.Duration `default:"invalid_duration"`
				}
				err := source.Default().Scan(&cfg)
				require.Error(t, err)
			},
		},
		{
			name: "unsupported_type",
			test: func(t *testing.T) {
				var cfg struct {
					Value []string `default:"test"`
				}
				err := source.Default().Scan(&cfg)
				require.Error(t, err)
				require.Contains(t, err.Error(), "unsupported type")
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
