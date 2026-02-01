package config

import "github.com/242617/core/config/source"

// ConfigEngine scans configuration from multiple sources in order.
type ConfigEngine interface {
	With(...source.ConfigSource) ConfigEngine
	Scan(any) error
}

// New creates config engine with Default() source.
func New() ConfigEngine {
	return &config{[]source.ConfigSource{source.Default()}}
}

type config struct{ sources []source.ConfigSource }

// With adds sources to scan order (later sources override earlier ones).
func (c *config) With(sources ...source.ConfigSource) ConfigEngine {
	c.sources = append(c.sources, sources...)
	return c
}

// Scan applies all sources to target, stopping on first error.
func (c *config) Scan(p any) error {
	for _, source := range c.sources {
		if err := source.Scan(p); err != nil {
			return err
		}
	}
	return nil
}
