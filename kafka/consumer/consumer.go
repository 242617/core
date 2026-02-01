package consumer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/242617/core/kafka"
	"github.com/242617/core/protocol"
)

// Consumer consumes messages from Kafka with consumer group support,
// handling rebalancing, offset commits, and graceful shutdown.
// Safe for concurrent use. Implements protocol.Lifecycle.
type Consumer struct {
	client  *kgo.Client
	handler kafka.Handler
	log     protocol.Logger

	brokers       []string
	topic         string
	groupID       string
	startOffset   kafka.StartOffset
	fetchMinBytes int32
	fetchMaxWait  time.Duration

	mu         sync.Mutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	started    bool
}

// New creates a new Kafka consumer with the provided options.
func New(options ...Option) (*Consumer, error) {
	consumer := &Consumer{}

	for _, option := range append(defaults(), options...) {
		if err := option(consumer); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	if consumer.log == nil {
		return nil, errors.New("empty logger")
	}
	if len(consumer.brokers) == 0 {
		return nil, ErrNoBrokers
	}
	if consumer.topic == "" {
		return nil, ErrNoTopic
	}
	if consumer.groupID == "" {
		return nil, ErrNoGroupID
	}
	if consumer.handler == nil {
		return nil, errors.New("empty handler")
	}

	// Build franz-go client options
	opts := []kgo.Opt{
		kgo.SeedBrokers(consumer.brokers...),
		kgo.ConsumerGroup(consumer.groupID),
		kgo.ConsumeTopics(consumer.topic),
	}

	// Configure start offset
	switch consumer.startOffset {
	case kafka.StartOffsetEarliest:
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	case kafka.StartOffsetLatest:
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	}

	// Configure fetch options
	if consumer.fetchMinBytes > 0 {
		opts = append(opts, kgo.FetchMinBytes(consumer.fetchMinBytes))
	}
	if consumer.fetchMaxWait > 0 {
		opts = append(opts, kgo.FetchMaxWait(consumer.fetchMaxWait))
	}

	// Add rebalancing hooks
	opts = append(opts,
		kgo.OnPartitionsRevoked(func(ctx context.Context, cl *kgo.Client, rev map[string][]int32) {
			for topic, partitions := range rev {
				consumer.log.Info(ctx, "partitions revoked", "topic", topic, "partitions", partitions, "group_id", consumer.groupID)
			}
			if err := cl.CommitUncommittedOffsets(ctx); err != nil {
				consumer.log.Error(ctx, "failed to commit on revoke", "err", err)
			}
		}),
		kgo.OnPartitionsAssigned(func(ctx context.Context, cl *kgo.Client, assigned map[string][]int32) {
			for topic, partitions := range assigned {
				consumer.log.Info(ctx, "partitions assigned", "topic", topic, "partitions", partitions, "group_id", consumer.groupID)
			}
		}),
		kgo.OnPartitionsLost(func(ctx context.Context, cl *kgo.Client, lost map[string][]int32) {
			for topic, partitions := range lost {
				consumer.log.Warn(ctx, "partitions lost", "topic", topic, "partitions", partitions, "group_id", consumer.groupID)
			}
		}),
	)

	consumer.log.Info(context.Background(), "attempting to create kafka client",
		"brokers", consumer.brokers,
		"topic", consumer.topic,
		"group_id", consumer.groupID)

	client, err := kgo.NewClient(opts...)
	if err != nil {
		consumer.log.Error(context.Background(), "failed to create kafka client",
			"brokers", consumer.brokers,
			"err", err)
		return nil, fmt.Errorf("create kafka client: %w", err)
	}

	consumer.log.Info(context.Background(), "kafka client created successfully",
		"brokers", consumer.brokers,
		"topic", consumer.topic,
		"group_id", consumer.groupID)

	consumer.client = client
	return consumer, nil
}

// Start begins consuming messages in the background. Implements protocol.Lifecycle.
// Returns immediately after starting the consumer goroutine.
func (c *Consumer) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return errors.New("consumer already started")
	}

	// Create a cancellable context for the consumer goroutine
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	c.started = true

	c.log.Info(ctx, "starting consumer",
		"brokers", c.brokers,
		"topic", c.topic,
		"group_id", c.groupID)

	// Start consumer goroutine directly
	c.wg.Add(1)
	go c.run()

	return nil
}

// Stop gracefully stops the consumer and waits for shutdown. Implements protocol.Lifecycle.
// Idempotent and safe to call multiple times.
func (c *Consumer) Stop(ctx context.Context) error {
	c.mu.Lock()
	if !c.started {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	c.log.Info(ctx, "stopping consumer")

	// Cancel the consumer context to signal shutdown
	if c.cancelFunc != nil {
		c.cancelFunc()
	}

	// Wait for consumer loop to finish or context timeout
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		c.log.Info(ctx, "consumer stopped")
	case <-ctx.Done():
		c.log.Warn(ctx, "consumer stop timeout")
		return ctx.Err()
	}

	// Clean shutdown
	c.client.LeaveGroup()
	c.client.Close()

	return nil
}

// run is the main consumer loop, executed in a goroutine.
func (c *Consumer) run() {
	defer c.wg.Done()
	defer func() {
		c.mu.Lock()
		c.started = false
		c.mu.Unlock()
	}()

	c.log.Info(c.ctx, "consumer loop started",
		"brokers", c.brokers,
		"topic", c.topic,
		"group_id", c.groupID)

	// Counter for logging empty polls (to avoid spamming logs)
	emptyPollCount := 0
	const maxEmptyPollLogs = 5

	for {
		// Check for cancellation before each poll
		select {
		case <-c.ctx.Done():
			c.log.Info(c.ctx, "consumer loop exiting due to context cancellation")
			return
		default:
		}

		// Poll for records
		fetches := c.client.PollFetches(c.ctx)

		// Check for context cancellation after poll
		if c.ctx.Err() != nil {
			c.log.Info(c.ctx, "consumer loop exiting due to context cancellation")
			return
		}

		// Handle fetch-level errors
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.log.Error(c.ctx, "fetch error",
					"topic", err.Topic,
					"partition", err.Partition,
					"err", err.Err)
			}
			// Continue polling even after errors (they are retriable)
			continue
		}

		// Check if we got any records
		recordCount := 0
		partitionCount := 0
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			recordCount += len(p.Records)
			if len(p.Records) > 0 {
				partitionCount++
			}
		})

		if recordCount == 0 {
			emptyPollCount++
			if emptyPollCount <= maxEmptyPollLogs {
				c.log.Debug(c.ctx, "poll returned no records",
					"empty_poll_count", emptyPollCount)
			}
			continue
		}

		// Log when we receive records (reset empty poll counter)
		emptyPollCount = 0
		c.log.Info(c.ctx, "received fetch",
			"records", recordCount,
			"partitions_with_data", partitionCount)

		// Process records per partition
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			if p.Err != nil {
				c.log.Error(c.ctx, "partition error",
					"topic", p.Topic,
					"partition", p.Partition,
					"err", p.Err)
				return
			}

			// Process each record
			for _, record := range p.Records {
				c.handleMessage(c.ctx, record)
			}

			// Commit offsets after processing all partition records
			if err := c.client.CommitRecords(c.ctx, p.Records...); err != nil {
				c.log.Error(c.ctx, "failed to commit offsets",
					"topic", p.Topic,
					"partition", p.Partition,
					"err", err)
			}
		})
	}
}

// handleMessage processes a single message.
func (c *Consumer) handleMessage(ctx context.Context, record *kgo.Record) {
	start := time.Now()

	// Convert to kafka.Message
	headers := make([]kafka.Header, len(record.Headers))
	for i, h := range record.Headers {
		headers[i] = kafka.Header{Key: h.Key, Value: h.Value}
	}

	msg := kafka.Message{
		Key:       record.Key,
		Value:     record.Value,
		Headers:   headers,
		Topic:     record.Topic,
		Partition: record.Partition,
	}

	// Call handler
	err := c.handler(ctx, msg)
	latency := time.Since(start)

	if err != nil {
		c.log.Error(ctx, "handler failed",
			"topic", record.Topic,
			"partition", record.Partition,
			"offset", record.Offset,
			"latency_ms", latency.Milliseconds(),
			"err", err)
		return
	}

	c.log.Debug(ctx, "message processed",
		"topic", record.Topic,
		"partition", record.Partition,
		"offset", record.Offset,
		"latency_ms", latency.Milliseconds())
}
