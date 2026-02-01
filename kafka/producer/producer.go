package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/242617/core/kafka"
	"github.com/242617/core/protocol"
	"github.com/twmb/franz-go/pkg/kgo"
)

// clientWrapper wraps franz-go client with additional functionality.
type clientWrapper struct {
	client  *kgo.Client
	brokers []string
	topic   string
	logger  protocol.Logger

	mu     sync.RWMutex
	closed bool
}

// isClosed returns true if the client is closed.
func (cw *clientWrapper) isClosed() bool {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	return cw.closed
}

// close closes the client gracefully.
// The method is idempotent and safe to call multiple times.
func (cw *clientWrapper) close(ctx context.Context) error {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	if cw.closed {
		return nil
	}

	cw.closed = true
	cw.client.Close()
	cw.logger.Info(ctx, "kafka client closed", "brokers", cw.brokers)
	return nil
}

// Producer provides methods for publishing messages to Kafka.
// It supports both async (Produce) and sync (ProduceSync) modes.
//
// The Producer is safe for concurrent use from multiple goroutines.
//
// Example usage with options:
//
//	producer, err := producer.New(
//	    producer.WithBrokers("localhost:9092"),
//	    producer.WithTopic("my-topic"),
//	    producer.WithLogger(log.New("kafka")),
//	)
//	if err != nil {
//	    return err
//	}
//	defer producer.Stop(ctx)
//
//	// Async produce
//	producer.Produce(ctx, kafka.Message{
//	    Key:   []byte("key"),
//	    Value: []byte("value"),
//	}, func(msg *kafka.Message, err error) {
//	    if err != nil {
//	        log.Printf("produce failed: %v", err)
//	        return
//	    }
//	    log.Printf("message produced: topic=%s partition=%d offset=%d",
//	        msg.Topic, msg.Partition, msg.Offset)
//	})
type Producer struct {
	client *clientWrapper
	logger protocol.Logger
}

// New creates a new Kafka producer client with options.
func New(options ...Option) (*Producer, error) {
	// Create default config
	cfg := &Config{}

	// Apply defaults
	for _, option := range defaults() {
		if err := option(cfg); err != nil {
			return nil, fmt.Errorf("apply default option: %w", err)
		}
	}

	// Apply user options
	for _, option := range options {
		if err := option(cfg); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	// Validate final configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid producer config: %w", err)
	}

	// Build franz-go options
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}

	cw := &clientWrapper{
		client:  client,
		brokers: cfg.Brokers,
		topic:   cfg.Topic,
		logger:  cfg.Logger,
	}

	return &Producer{
		client: cw,
		logger: cfg.Logger,
	}, nil
}

// Produce sends a message asynchronously.
// The callback is invoked when the message is successfully produced or fails.
//
// If the producer is closed, the callback will be invoked with kafka.ErrClosed immediately.
func (p *Producer) Produce(ctx context.Context, msg kafka.Message, callback func(*kafka.Message, error)) {
	if p.client.isClosed() {
		if callback != nil {
			callback(&msg, kafka.ErrClosed)
		}
		return
	}

	start := time.Now()
	topic := msg.Topic
	if topic == "" {
		topic = p.client.topic
	}

	record := &kgo.Record{
		Topic:     topic,
		Partition: msg.Partition,
		Key:       msg.Key,
		Value:     msg.Value,
	}

	// Add headers
	if len(msg.Headers) > 0 {
		record.Headers = make([]kgo.RecordHeader, len(msg.Headers))
		for i, h := range msg.Headers {
			record.Headers[i] = kgo.RecordHeader{Key: h.Key, Value: h.Value}
		}
	}

	p.client.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		latency := time.Since(start)
		success := err == nil

		if success {
			p.logger.Debug(ctx, "message produced successfully",
				"topic", r.Topic,
				"partition", r.Partition,
				"offset", r.Offset,
				"latency_ms", latency.Milliseconds())
		} else {
			p.logger.Error(ctx, "message produce failed",
				"topic", topic,
				"err", err,
				"latency_ms", latency.Milliseconds())
		}

		if callback != nil {
			var callbackMsg kafka.Message
			if r != nil {
				callbackMsg = kafka.Message{
					Key:       r.Key,
					Value:     r.Value,
					Topic:     r.Topic,
					Partition: r.Partition,
				}
			}
			callback(&callbackMsg, err)
		}
	})
}

// ProduceSync sends messages synchronously, waiting for all to complete.
// Returns the first error if any messages failed.
func (p *Producer) ProduceSync(ctx context.Context, msgs ...kafka.Message) error {
	if p.client.isClosed() {
		return kafka.ErrClosed
	}

	if len(msgs) == 0 {
		return nil
	}

	records := make([]*kgo.Record, len(msgs))
	for i, msg := range msgs {
		topic := msg.Topic
		if topic == "" {
			topic = p.client.topic
		}

		records[i] = &kgo.Record{
			Topic:     topic,
			Partition: msg.Partition,
			Key:       msg.Key,
			Value:     msg.Value,
		}

		// Add headers
		if len(msg.Headers) > 0 {
			records[i].Headers = make([]kgo.RecordHeader, len(msg.Headers))
			for j, h := range msg.Headers {
				records[i].Headers[j] = kgo.RecordHeader{Key: h.Key, Value: h.Value}
			}
		}
	}

	start := time.Now()
	results := p.client.client.ProduceSync(ctx, records...)
	latency := time.Since(start)

	if err := results.FirstErr(); err != nil {
		p.logger.Error(ctx, "sync produce failed",
			"count", len(msgs),
			"err", err,
			"latency_ms", latency.Milliseconds())
		return fmt.Errorf("produce messages: %w", err)
	}

	for _, result := range results {
		p.logger.Debug(ctx, "message produced successfully",
			"topic", result.Record.Topic,
			"partition", result.Record.Partition,
			"offset", result.Record.Offset)
	}

	p.logger.Debug(ctx, "sync produce completed",
		"count", len(msgs),
		"latency_ms", latency.Milliseconds())

	return nil
}

// Close closes the producer gracefully.
// The method is idempotent and safe to call multiple times.
func (p *Producer) Close(ctx context.Context) error {
	if p.client.isClosed() {
		return nil
	}

	p.logger.Info(ctx, "closing producer", "brokers", p.client.brokers)
	return p.client.close(ctx)
}

// Start is a no-op for the producer as it's ready to produce immediately after creation.
// This method implements protocol.Lifecycle interface.
func (p *Producer) Start(ctx context.Context) error {
	p.logger.Debug(ctx, "producer started (ready to produce)")
	return nil
}

// Stop stops the producer by closing it.
// This method implements protocol.Lifecycle interface.
func (p *Producer) Stop(ctx context.Context) error {
	return p.Close(ctx)
}
