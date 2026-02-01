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

// Producer publishes messages to Kafka with async and sync modes.
// Safe for concurrent use. Implements protocol.Lifecycle.
type Producer struct {
	client KafkaProducer
	topic  string
	logger protocol.Logger

	mu     sync.RWMutex
	closed bool
}

// New creates a new Kafka producer with the provided options.
func New(options ...Option) (*Producer, error) {
	cfg := &Config{}

	// Apply defaults then user options
	for _, option := range append(defaults(), options...) {
		if err := option(cfg); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Use custom producer if provided (for testing), otherwise create real one
	var client KafkaProducer
	if cfg.producer != nil {
		client = cfg.producer
	} else {
		kgoClient, err := kgo.NewClient(kgo.SeedBrokers(cfg.Brokers...))
		if err != nil {
			return nil, fmt.Errorf("create kafka client: %w", err)
		}
		client = NewRealProducer(kgoClient)
	}

	return &Producer{
		client: client,
		topic:  cfg.Topic,
		logger: cfg.Logger,
	}, nil
}

// Produce sends a message asynchronously. Callback is invoked on completion or error.
// If producer is closed, callback receives kafka.ErrClosed immediately.
func (p *Producer) Produce(ctx context.Context, msg kafka.Message, callback func(*kafka.Message, error)) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		if callback != nil {
			callback(&msg, kafka.ErrClosed)
		}
		return
	}
	p.mu.RUnlock()

	start := time.Now()
	topic := msg.Topic
	if topic == "" {
		topic = p.topic
	}

	record := &kgo.Record{
		Topic:     topic,
		Partition: msg.Partition,
		Key:       msg.Key,
		Value:     msg.Value,
	}

	// Convert headers
	if len(msg.Headers) > 0 {
		record.Headers = make([]kgo.RecordHeader, len(msg.Headers))
		for i, h := range msg.Headers {
			record.Headers[i] = kgo.RecordHeader{Key: h.Key, Value: h.Value}
		}
	}

	p.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		latency := time.Since(start)

		if err == nil {
			p.logger.Debug(ctx, "message produced",
				"topic", r.Topic,
				"partition", r.Partition,
				"offset", r.Offset,
				"latency_ms", latency.Milliseconds())
		} else {
			p.logger.Error(ctx, "produce failed",
				"topic", topic,
				"latency_ms", latency.Milliseconds(),
				"err", err)
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
// Returns the first error encountered, if any.
func (p *Producer) ProduceSync(ctx context.Context, msgs ...kafka.Message) error {
	p.mu.RLock()
	closed := p.closed
	p.mu.RUnlock()

	if closed {
		return kafka.ErrClosed
	}

	if len(msgs) == 0 {
		return nil
	}

	records := make([]*kgo.Record, len(msgs))
	for i, msg := range msgs {
		topic := msg.Topic
		if topic == "" {
			topic = p.topic
		}

		records[i] = &kgo.Record{
			Topic:     topic,
			Partition: msg.Partition,
			Key:       msg.Key,
			Value:     msg.Value,
		}

		// Convert headers
		if len(msg.Headers) > 0 {
			records[i].Headers = make([]kgo.RecordHeader, len(msg.Headers))
			for j, h := range msg.Headers {
				records[i].Headers[j] = kgo.RecordHeader{Key: h.Key, Value: h.Value}
			}
		}
	}

	start := time.Now()
	results := p.client.ProduceSync(ctx, records...)
	latency := time.Since(start)

	if err := results.FirstErr(); err != nil {
		p.logger.Error(ctx, "sync produce failed",
			"count", len(msgs),
			"latency_ms", latency.Milliseconds(),
			"err", err)
		return fmt.Errorf("produce messages: %w", err)
	}

	p.logger.Debug(ctx, "sync produce completed",
		"count", len(msgs),
		"latency_ms", latency.Milliseconds())

	return nil
}

// Start is a no-op (producer is ready immediately). Implements protocol.Lifecycle.
func (p *Producer) Start(ctx context.Context) error {
	p.logger.Debug(ctx, "producer ready")
	return nil
}

// Stop closes the producer gracefully. Implements protocol.Lifecycle.
// Idempotent and safe to call multiple times.
func (p *Producer) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	p.client.Close()
	p.logger.Info(ctx, "producer stopped")

	return nil
}
