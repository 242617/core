package file_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/242617/core/config/source/file"
)

type yamlConfig struct {
	StringVal   string        `yaml:"string_val"`
	IntVal      int           `yaml:"int_val"`
	UintVal     uint          `yaml:"uint_val"`
	FloatVal    float64       `yaml:"float_val"`
	BoolVal     bool          `yaml:"bool_val"`
	DurationVal time.Duration `yaml:"duration_val"`
	Nested      struct {
		ChildString string `yaml:"child_string"`
		ChildInt    int    `yaml:"child_int"`
	} `yaml:"nested"`
	NoTag string
}

func TestYAML_Success(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "basic_string",
			test: func(t *testing.T) {
				content := `string_val: "hello"`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "hello", cfg.StringVal)
			},
		},
		{
			name: "int",
			test: func(t *testing.T) {
				content := `int_val: 42`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 42, cfg.IntVal)
			},
		},
		{
			name: "uint",
			test: func(t *testing.T) {
				content := `uint_val: 100`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, uint(100), cfg.UintVal)
			},
		},
		{
			name: "float64",
			test: func(t *testing.T) {
				content := `float_val: 3.14159`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 3.14159, cfg.FloatVal)
			},
		},
		{
			name: "bool_true",
			test: func(t *testing.T) {
				content := `bool_val: true`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.True(t, cfg.BoolVal)
			},
		},
		{
			name: "bool_false",
			test: func(t *testing.T) {
				content := `bool_val: false`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.False(t, cfg.BoolVal)
			},
		},
		{
			name: "duration",
			test: func(t *testing.T) {
				content := `duration_val: "5s"`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, 5*time.Second, cfg.DurationVal)
			},
		},
		{
			name: "nested_struct",
			test: func(t *testing.T) {
				content := `nested:
  child_string: "child_value"
  child_int: 42`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "child_value", cfg.Nested.ChildString)
				require.Equal(t, 42, cfg.Nested.ChildInt)
			},
		},
		{
			name: "complete_config",
			test: func(t *testing.T) {
				content := `string_val: "complete"
int_val: 123
uint_val: 456
float_val: 2.718
bool_val: true
duration_val: "10m"
nested:
  child_string: "nested_child"
  child_int: 789`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "complete", cfg.StringVal)
				require.Equal(t, 123, cfg.IntVal)
				require.Equal(t, uint(456), cfg.UintVal)
				require.Equal(t, 2.718, cfg.FloatVal)
				require.True(t, cfg.BoolVal)
				require.Equal(t, 10*time.Minute, cfg.DurationVal)
				require.Equal(t, "nested_child", cfg.Nested.ChildString)
				require.Equal(t, 789, cfg.Nested.ChildInt)
			},
		},
		{
			name: "empty_values",
			test: func(t *testing.T) {
				content := `string_val: ""
int_val: 0
bool_val: false`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Empty(t, cfg.StringVal)
				require.Equal(t, 0, cfg.IntVal)
				require.False(t, cfg.BoolVal)
			},
		},
		{
			name: "multiline_string",
			test: func(t *testing.T) {
				content := `string_val: |
  line1
  line2
  line3`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
				require.Equal(t, "line1\nline2\nline3", cfg.StringVal)
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

func TestYAML_Error(t *testing.T) {
	for _, tc := range []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "file_not_found",
			test: func(t *testing.T) {
				filePath := filepath.Join(t.TempDir(), "nonexistent.yaml")
				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.Error(t, err)
				require.Contains(t, err.Error(), "no such file")
			},
		},
		{
			name: "invalid_yaml_syntax",
			test: func(t *testing.T) {
				content := `string_val: [unclosed array`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.Error(t, err)
			},
		},
		{
			name: "nil_pointer",
			test: func(t *testing.T) {
				content := `string_val: "test"`
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg *yamlConfig
				err := file.YAML(filePath).Scan(cfg)
				require.Error(t, err)
			},
		},
		{
			name: "empty_file",
			test: func(t *testing.T) {
				content := ``
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "config.yaml")
				require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

				var cfg yamlConfig
				err := file.YAML(filePath).Scan(&cfg)
				require.NoError(t, err)
			},
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

func TestYAML_ArrayFields(t *testing.T) {
	type arrayConfig struct {
		Strings []string `yaml:"strings"`
		Ints    []int    `yaml:"ints"`
	}

	content := `strings:
  - "a"
  - "b"
  - "c"
ints:
  - 1
  - 2
  - 3`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yaml")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	var cfg arrayConfig
	err := file.YAML(filePath).Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, []string{"a", "b", "c"}, cfg.Strings)
	require.Equal(t, []int{1, 2, 3}, cfg.Ints)
}

func TestYAML_MapFields(t *testing.T) {
	type mapConfig struct {
		StringMap map[string]string `yaml:"string_map"`
	}

	content := `string_map:
  key1: "value1"
  key2: "value2"`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "config.yaml")
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	var cfg mapConfig
	err := file.YAML(filePath).Scan(&cfg)
	require.NoError(t, err)
	require.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, cfg.StringMap)
}
