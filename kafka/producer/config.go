package producer

import (
	stderrors "errors"

	"github.com/242617/core/kafka"
	"github.com/242617/core/protocol"
)

var (
	// ErrNoBrokers indicates missing broker configuration.
	ErrNoBrokers = kafka.ErrNoBrokers

	// ErrNoTopic indicates missing topic configuration.
	ErrNoTopic = kafka.ErrNoTopic
)

// Config contains configuration for the Kafka producer.
type Config struct {
	// Brokers is the list of Kafka broker addresses (required).
	Brokers []string `yaml:"brokers"`

	// Topic is the default Kafka topic for messages (required unless specified in Produce calls).
	Topic string `yaml:"topic"`

	// Logger is the logger for structured logging (optional, uses DefaultLogger).
	Logger protocol.Logger `yaml:"-"`

	// producer is the Kafka producer implementation (optional, primarily for testing).
	producer KafkaProducer `yaml:"-"`
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	if len(c.Brokers) == 0 {
		return ErrNoBrokers
	}
	if c.Topic == "" {
		return ErrNoTopic
	}
	return nil
}

// Option is a function that configures the Producer.
type Option func(*Config) error

// defaults returns default producer configuration.
func defaults() []Option {
	return []Option{
		WithLogger(protocol.NopLogger{}),
	}
}

// WithBrokers sets the Kafka broker addresses for producer.
func WithBrokers(brokers ...string) Option {
	return func(cfg *Config) error {
		if len(brokers) == 0 {
			return stderrors.New("brokers cannot be empty")
		}
		cfg.Brokers = brokers
		return nil
	}
}

// WithTopic sets the default topic for producer messages.
func WithTopic(topic string) Option {
	return func(cfg *Config) error {
		if topic == "" {
			return stderrors.New("topic cannot be empty")
		}
		cfg.Topic = topic
		return nil
	}
}

// WithLogger sets the logger for producer.
func WithLogger(logger protocol.Logger) Option {
	return func(cfg *Config) error {
		if logger == nil {
			return stderrors.New("logger cannot be nil")
		}
		cfg.Logger = logger
		return nil
	}
}

// WithKafkaProducer sets a custom Kafka producer implementation.
// This is primarily used for testing with mock implementations.
func WithKafkaProducer(p KafkaProducer) Option {
	return func(cfg *Config) error {
		cfg.producer = p
		return nil
	}
}
