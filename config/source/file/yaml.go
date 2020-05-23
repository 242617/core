package file

import (
	"io/ioutil"

	yaml2 "gopkg.in/yaml.v2"

	"github.com/242617/core/config/source"
)

// YAML creates config source that fills config with values from yaml-file
func YAML(file string) source.ConfigSource {
	return &yaml{file}
}

type yaml struct{ file string }

func (y *yaml) Scan(p interface{}) error {
	barr, err := ioutil.ReadFile(y.file)
	if err != nil {
		return err
	}

	if err = yaml2.Unmarshal(barr, p); err != nil {
		return err
	}

	return nil
}
