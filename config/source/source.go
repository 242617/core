package source

// ConfigSource provides values to scan into a config struct.
type ConfigSource interface {
	Scan(p interface{}) error
}
