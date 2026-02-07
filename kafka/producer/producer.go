package producer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/242617/core/kafka"
	"github.com/242617/core/protocol"
)

// Producer publishes messages to Kafka with async and sync modes.
// Safe for concurrent use. Implements protocol.Lifecycle.
type Producer struct {
	log     protocol.Logger
	brokers []string
	topic   string

	client *kgo.Client
	mu     sync.RWMutex
	closed bool
}

// New creates a new Kafka producer with the provided options.
func New(options ...Option) (*Producer, error) {
	var p Producer

	// Apply defaults then user options
	for _, option := range append(defaults(), options...) {
		if err := option(&p); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	// Use custom producer if provided (for testing), otherwise create real one
	client, err := kgo.NewClient(kgo.SeedBrokers(p.brokers...))
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}
	p.client = client

	return &p, nil
}

// Produce sends messages synchronously, waiting for all to complete.
// Returns the first error encountered, if any.
func (p *Producer) Produce(ctx context.Context, msgs ...kafka.Message) error {
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
		records[i] = &kgo.Record{
			Topic: p.topic,
			Key:   msg.Key,
			Value: msg.Value,
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
		p.log.Error(ctx, "sync produce failed",
			"count", len(msgs),
			"latency_ms", latency.Milliseconds(),
			"err", err)
		return fmt.Errorf("produce messages: %w", err)
	}

	p.log.Debug(ctx, "sync produce completed",
		"count", len(msgs),
		"latency_ms", latency.Milliseconds())

	return nil
}

// Start is a no-op (producer is ready immediately). Implements protocol.Lifecycle.
func (p *Producer) Start(ctx context.Context) error { return nil }

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
	p.log.Info(ctx, "producer stopped")

	return nil
}
