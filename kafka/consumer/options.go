package consumer

import (
	stderrors "errors"
	"fmt"
	"time"

	"github.com/242617/core/kafka"
	"github.com/242617/core/protocol"
)

var (
	// ErrNoBrokers is returned when no brokers are provided.
	ErrNoBrokers = kafka.ErrNoBrokers

	// ErrNoTopic is returned when no topic is provided.
	ErrNoTopic = kafka.ErrNoTopic

	// ErrNoGroupID is returned when no group ID is provided.
	ErrNoGroupID = kafka.ErrNoGroupID
)

// Config contains configuration for the Kafka consumer.
type Config struct {
	// Brokers is the list of Kafka broker addresses (required).
	Brokers []string `yaml:"brokers"`

	// Topic is the Kafka topic to consume from (required).
	Topic string `yaml:"topic"`

	// GroupID is the consumer group ID (required for consumer groups).
	GroupID string `yaml:"group_id"`
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	if len(c.Brokers) == 0 {
		return ErrNoBrokers
	}
	if c.Topic == "" {
		return ErrNoTopic
	}
	if c.GroupID == "" {
		return ErrNoGroupID
	}
	return nil
}

// Option is a function that configures the Consumer.
type Option func(*Consumer) error

// defaults returns default consumer configuration.
func defaults() []Option {
	return []Option{
		WithLogger(protocol.NopLogger{}),
		WithFetchMinBytes(1),
		WithFetchMaxWait(500 * time.Millisecond),
		WithStartOffset(kafka.StartOffsetLatest),
	}
}

// WithConfig sets all consumer configuration from a Config struct.
// This is recommended way to configure the consumer when loading from YAML or similar sources.
func WithConfig(cfg Config) Option {
	return func(consumer *Consumer) error {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid consumer config: %w", err)
		}

		consumer.brokers = cfg.Brokers
		consumer.topic = cfg.Topic
		consumer.groupID = cfg.GroupID

		return nil
	}
}

// WithBrokers sets the Kafka broker addresses for consumer.
func WithBrokers(brokers ...string) Option {
	return func(consumer *Consumer) error {
		if len(brokers) == 0 {
			return stderrors.New("brokers cannot be empty")
		}
		consumer.brokers = brokers
		return nil
	}
}

// WithTopic sets the topic for consumer.
func WithTopic(topic string) Option {
	return func(consumer *Consumer) error {
		if topic == "" {
			return stderrors.New("topic cannot be empty")
		}
		consumer.topic = topic
		return nil
	}
}

// WithGroupID sets the consumer group ID.
func WithGroupID(groupID string) Option {
	return func(consumer *Consumer) error {
		if groupID == "" {
			return stderrors.New("group ID cannot be empty")
		}
		consumer.groupID = groupID
		return nil
	}
}

// WithHandler sets the message handler function for consumer.
func WithHandler(handler kafka.Handler) Option {
	return func(consumer *Consumer) error {
		if handler == nil {
			return stderrors.New("handler cannot be nil")
		}
		consumer.handler = handler
		return nil
	}
}

// WithLogger sets the logger for consumer.
func WithLogger(logger protocol.Logger) Option {
	return func(consumer *Consumer) error {
		if logger == nil {
			return stderrors.New("logger cannot be nil")
		}
		consumer.log = logger
		return nil
	}
}

// WithStartOffset sets the initial offset for new consumer groups.
func WithStartOffset(offset kafka.StartOffset) Option {
	return func(consumer *Consumer) error {
		consumer.startOffset = offset
		return nil
	}
}

// WithFetchMinBytes sets the minimum number of bytes to wait for in a fetch.
func WithFetchMinBytes(minBytes int32) Option {
	return func(consumer *Consumer) error {
		consumer.fetchMinBytes = minBytes
		return nil
	}
}

// WithFetchMaxWait sets the maximum time to wait for FetchMinBytes.
func WithFetchMaxWait(maxWait time.Duration) Option {
	return func(consumer *Consumer) error {
		consumer.fetchMaxWait = maxWait
		return nil
	}
}
