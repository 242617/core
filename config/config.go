package config

import "github.com/242617/core/config/source"

// ConfigEngine is an interface for config scanner
type ConfigEngine interface {
	With(...source.ConfigSource) ConfigEngine
	Scan(interface{}) error
}

// New creates a new config engine with default scanner
func New() ConfigEngine {
	return &config{[]source.ConfigSource{source.Default()}}
}

type config struct{ sources []source.ConfigSource }

// With adds source(s) for engine. Make sure you are adding sources in desired order.
func (c *config) With(sources ...source.ConfigSource) ConfigEngine {
	c.sources = append(c.sources, sources...)
	return c
}

// Scan returns error of scanning sources into config
func (c *config) Scan(p interface{}) error {
	for _, source := range c.sources {
		if err := source.Scan(p); err != nil {
			return err
		}
	}
	return nil
}
