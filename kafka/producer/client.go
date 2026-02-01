package producer

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

// KafkaProducer is the interface for Kafka producer operations.
// This allows mocking for unit tests.
type KafkaProducer interface {
	Produce(ctx context.Context, r *kgo.Record, ack func(*kgo.Record, error))
	ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults
	Close()
}

// RealProducer wraps the actual kgo.Client for production use.
type RealProducer struct {
	client *kgo.Client
}

// NewRealProducer creates a new RealProducer wrapping the kgo.Client.
func NewRealProducer(client *kgo.Client) KafkaProducer {
	return &RealProducer{client: client}
}

// Produce implements KafkaProducer.
func (rp *RealProducer) Produce(ctx context.Context, r *kgo.Record, ack func(*kgo.Record, error)) {
	rp.client.Produce(ctx, r, ack)
}

// ProduceSync implements KafkaProducer.
func (rp *RealProducer) ProduceSync(ctx context.Context, rs ...*kgo.Record) kgo.ProduceResults {
	return rp.client.ProduceSync(ctx, rs...)
}

// Close implements KafkaProducer.
func (rp *RealProducer) Close() {
	rp.client.Close()
}
