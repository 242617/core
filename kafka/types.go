package kafka

import "context"

// Message represents a Kafka message with key, value, headers, and routing info.
type Message struct {
	Key       []byte
	Value     []byte
	Headers   []Header
	Topic     string // Overrides default topic if set
	Partition int32  // Specific partition (optional)
}

// Header represents a Kafka message header key-value pair.
type Header struct {
	Key   string
	Value []byte
}

// Handler processes consumed messages. Returns nil on success, or error if message
// should be retried. Handler errors are logged but don't stop the consumer loop.
type Handler func(ctx context.Context, msg Message) error

// StartOffset defines the initial offset for new consumer groups.
type StartOffset int

const (
	// StartOffsetLatest begins consuming from newest messages.
	StartOffsetLatest StartOffset = iota

	// StartOffsetEarliest begins consuming from oldest messages.
	StartOffsetEarliest
)
