// Package config provides multi-source configuration loading with override capability.
//
// Sources are applied in order, with later sources overriding earlier ones:
// defaults → environment variables → files.
//
// Example:
//
//	var cfg struct {
//	    DB      pgrepo.Config `yaml:"db"`
//	    Timeout time.Duration `yaml:"timeout" default:"30s"`
//	}
//	config.New().
//	    With(source.Env()).
//	    With(file.YAML("config.yaml")).
//	    Scan(&cfg)
package config
