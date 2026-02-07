package producer

import (
	"errors"
	"fmt"

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
type Option func(*Producer) error

// defaults returns default producer configuration.
func defaults() []Option {
	return []Option{
		WithLogger(protocol.NopLogger{}),
	}
}

func WithConfig(cfg Config) Option {
	return func(p *Producer) error {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}
		p.brokers = cfg.Brokers
		p.topic = cfg.Topic
		return nil
	}
}

// WithBrokers sets the Kafka broker addresses for producer.
func WithBrokers(brokers ...string) Option {
	return func(p *Producer) error {
		if len(brokers) == 0 {
			return errors.New("brokers cannot be empty")
		}
		p.brokers = brokers
		return nil
	}
}

// WithTopic sets the default topic for producer messages.
func WithTopic(topic string) Option {
	return func(p *Producer) error {
		if topic == "" {
			return errors.New("topic cannot be empty")
		}
		p.topic = topic
		return nil
	}
}

// WithLogger sets the logger for producer.
func WithLogger(logger protocol.Logger) Option {
	return func(p *Producer) error {
		if logger == nil {
			return errors.New("logger cannot be nil")
		}
		p.log = logger
		return nil
	}
}
