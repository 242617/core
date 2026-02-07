// Package kafka provides producer and consumer implementations for Apache Kafka.
//
// Both producer and consumer implement the protocol.Lifecycle interface.
// The consumer supports consumer groups for scalable consumption.
//
// Producer example:
//
//	producer, err := producer.New(producer.WithConfig(cfg.Producer))
//	msg := kafka.Message{Key: []byte("key"), Value: []byte("value")}
//	producer.Produce(ctx, msg)
//
// Consumer example:
//
//	consumer, err := consumer.New(
//	    consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error {
//	        return nil // commits offset
//	    }),
//	)
package kafka
