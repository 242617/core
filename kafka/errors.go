package kafka

import "errors"

var (
	// ErrClosed indicates an operation on a closed client.
	ErrClosed = errors.New("kafka client is closed")

	// ErrNoBrokers indicates missing broker configuration.
	ErrNoBrokers = errors.New("no brokers provided")

	// ErrNoTopic indicates missing topic configuration.
	ErrNoTopic = errors.New("no topic provided")

	// ErrNoGroupID indicates missing consumer group ID.
	ErrNoGroupID = errors.New("no group ID provided")
)
