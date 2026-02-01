package file

import (
	"fmt"
	"os"
	"reflect"

	yamlv3 "gopkg.in/yaml.v3"

	"github.com/242617/core/config/source"
)

// YAML loads struct fields from a YAML file.
func YAML(file string) source.ConfigSource {
	return &yaml{file}
}

type yaml struct{ file string }

func (y *yaml) Scan(p interface{}) error {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("unexpected kind: %q", v.Kind())
	}

	data, err := os.ReadFile(y.file)
	if err != nil {
		return err
	}

	if err = yamlv3.Unmarshal(data, p); err != nil {
		return err
	}

	return nil
}
