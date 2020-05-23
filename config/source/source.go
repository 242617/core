package source

// ConfigSource is an interface for config source
type ConfigSource interface {
	Scan(p interface{}) error
}
