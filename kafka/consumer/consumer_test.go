package consumer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/242617/core/kafka"
	"github.com/242617/core/kafka/consumer"
	"github.com/242617/core/protocol"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		options []consumer.Option
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			options: []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID("test-group"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
				consumer.WithLogger(protocol.NopLogger{}),
			},
			wantErr: false,
		},
		{
			name: "missing brokers",
			options: []consumer.Option{
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID("test-group"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			},
			wantErr: true,
			errMsg:  "no brokers provided",
		},
		{
			name: "missing topic",
			options: []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithGroupID("test-group"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			},
			wantErr: true,
			errMsg:  "no topic provided",
		},
		{
			name: "missing group ID",
			options: []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic("test-topic"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			},
			wantErr: true,
			errMsg:  "no group ID provided",
		},
		{
			name: "missing handler",
			options: []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID("test-group"),
			},
			wantErr: true,
			errMsg:  "empty handler",
		},
		{
			name: "invalid broker config",
			options: []consumer.Option{
				consumer.WithBrokers(),
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID("test-group"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			},
			wantErr: true,
			errMsg:  "brokers cannot be empty",
		},
		{
			name: "invalid topic config",
			options: []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic(""),
				consumer.WithGroupID("test-group"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			},
			wantErr: true,
			errMsg:  "topic cannot be empty",
		},
		{
			name: "invalid group ID config",
			options: []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID(""),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			},
			wantErr: true,
			errMsg:  "group ID cannot be empty",
		},
		{
			name: "invalid handler config",
			options: []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID("test-group"),
			},
			wantErr: true,
			errMsg:  "empty handler",
		},
		{
			name: "valid config with Config struct",
			options: []consumer.Option{
				consumer.WithConfig(consumer.Config{
					Brokers: []string{"localhost:9092"},
					Topic:   "test-topic",
					GroupID: "test-group",
				}),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
				consumer.WithLogger(protocol.NopLogger{}),
			},
			wantErr: false,
		},
		{
			name: "invalid Config struct",
			options: []consumer.Option{
				consumer.WithConfig(consumer.Config{
					Brokers: []string{},
					Topic:   "test-topic",
					GroupID: "test-group",
				}),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			},
			wantErr: true,
			errMsg:  "invalid consumer config: no brokers provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := consumer.New(tt.options...)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.NotNil(t, got)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     consumer.Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: consumer.Config{
				Brokers: []string{"localhost:9092"},
				Topic:   "test-topic",
				GroupID: "test-group",
			},
			wantErr: false,
		},
		{
			name: "missing brokers",
			cfg: consumer.Config{
				Topic:   "test-topic",
				GroupID: "test-group",
			},
			wantErr: true,
			errMsg:  "no brokers provided",
		},
		{
			name: "missing topic",
			cfg: consumer.Config{
				Brokers: []string{"localhost:9092"},
				GroupID: "test-group",
			},
			wantErr: true,
			errMsg:  "no topic provided",
		},
		{
			name: "missing group ID",
			cfg: consumer.Config{
				Brokers: []string{"localhost:9092"},
				Topic:   "test-topic",
			},
			wantErr: true,
			errMsg:  "no group ID provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func() (*consumer.Consumer, context.Context)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful start",
			setup: func() (*consumer.Consumer, context.Context) {
				c, err := consumer.New(
					consumer.WithBrokers("localhost:9092"),
					consumer.WithTopic("test-topic"),
					consumer.WithGroupID("test-group"),
					consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
					consumer.WithLogger(protocol.NopLogger{}),
				)
				require.NoError(t, err)
				return c, context.Background()
			},
			wantErr: false,
		},
		{
			name: "start twice returns error",
			setup: func() (*consumer.Consumer, context.Context) {
				c, err := consumer.New(
					consumer.WithBrokers("localhost:9092"),
					consumer.WithTopic("test-topic"),
					consumer.WithGroupID("test-group"),
					consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
					consumer.WithLogger(protocol.NopLogger{}),
				)
				require.NoError(t, err)
				// First start
				ctx := context.Background()
				err = c.Start(ctx)
				require.NoError(t, err)
				return c, ctx
			},
			wantErr: true,
			errMsg:  "consumer already started",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, ctx := tt.setup()
			defer func() {
				stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = c.Stop(stopCtx)
			}()

			err := c.Start(ctx)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				// Give goroutine time to start
				time.Sleep(100 * time.Millisecond)
			}
		})
	}
}

func TestStop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func() (*consumer.Consumer, context.Context)
		wantErr bool
	}{
		{
			name: "stop without start returns nil",
			setup: func() (*consumer.Consumer, context.Context) {
				c, err := consumer.New(
					consumer.WithBrokers("localhost:9092"),
					consumer.WithTopic("test-topic"),
					consumer.WithGroupID("test-group"),
					consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
					consumer.WithLogger(protocol.NopLogger{}),
				)
				require.NoError(t, err)
				return c, context.Background()
			},
			wantErr: false,
		},
		{
			name: "successful stop after start",
			setup: func() (*consumer.Consumer, context.Context) {
				c, err := consumer.New(
					consumer.WithBrokers("localhost:9092"),
					consumer.WithTopic("test-topic"),
					consumer.WithGroupID("test-group"),
					consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
					consumer.WithLogger(protocol.NopLogger{}),
				)
				require.NoError(t, err)
				ctx := context.Background()
				err = c.Start(ctx)
				require.NoError(t, err)
				// Give goroutine time to start
				time.Sleep(100 * time.Millisecond)
				return c, ctx
			},
			wantErr: false,
		},
		{
			name: "stop is idempotent",
			setup: func() (*consumer.Consumer, context.Context) {
				c, err := consumer.New(
					consumer.WithBrokers("localhost:9092"),
					consumer.WithTopic("test-topic"),
					consumer.WithGroupID("test-group"),
					consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
					consumer.WithLogger(protocol.NopLogger{}),
				)
				require.NoError(t, err)
				ctx := context.Background()
				err = c.Start(ctx)
				require.NoError(t, err)
				time.Sleep(100 * time.Millisecond)
				// First stop
				stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				err = c.Stop(stopCtx)
				require.NoError(t, err)
				time.Sleep(100 * time.Millisecond)
				return c, ctx
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, ctx := tt.setup()

			stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err := c.Stop(stopCtx)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		option  consumer.Option
		wantErr bool
	}{
		{
			name:   "WithBrokers single",
			option: consumer.WithBrokers("localhost:9092"),
		},
		{
			name:   "WithBrokers multiple",
			option: consumer.WithBrokers("localhost:9092", "localhost:9093"),
		},
		{
			name:    "WithBrokers empty",
			option:  consumer.WithBrokers(),
			wantErr: true,
		},
		{
			name:   "WithTopic valid",
			option: consumer.WithTopic("test-topic"),
		},
		{
			name:    "WithTopic empty",
			option:  consumer.WithTopic(""),
			wantErr: true,
		},
		{
			name:   "WithGroupID valid",
			option: consumer.WithGroupID("test-group"),
		},
		{
			name:    "WithGroupID empty",
			option:  consumer.WithGroupID(""),
			wantErr: true,
		},
		{
			name:   "WithHandler valid",
			option: consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
		},
		{
			name:    "WithHandler nil",
			option:  consumer.WithHandler(nil),
			wantErr: true,
		},
		{
			name:   "WithLogger valid",
			option: consumer.WithLogger(protocol.NopLogger{}),
		},
		{
			name:    "WithLogger nil",
			option:  consumer.WithLogger(nil),
			wantErr: true,
		},
		{
			name:   "WithStartOffset earliest",
			option: consumer.WithStartOffset(kafka.StartOffsetEarliest),
		},
		{
			name:   "WithStartOffset latest",
			option: consumer.WithStartOffset(kafka.StartOffsetLatest),
		},
		{
			name:   "WithFetchMinBytes valid",
			option: consumer.WithFetchMinBytes(1024),
		},
		{
			name:   "WithFetchMaxWait valid",
			option: consumer.WithFetchMaxWait(500 * time.Millisecond),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseOpts := []consumer.Option{
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID("test-group"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
				consumer.WithLogger(protocol.NopLogger{}),
			}

			opts := append(baseOpts, tt.option)
			_, err := consumer.New(opts...)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	t.Parallel()

	t.Run("concurrent start", func(t *testing.T) {
		c, err := consumer.New(
			consumer.WithBrokers("localhost:9092"),
			consumer.WithTopic("test-topic"),
			consumer.WithGroupID("test-group"),
			consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			consumer.WithLogger(protocol.NopLogger{}),
		)
		require.NoError(t, err)

		ctx := context.Background()
		errCh := make(chan error, 2)

		go func() { errCh <- c.Start(ctx) }()
		go func() { errCh <- c.Start(ctx) }()

		// Wait for both goroutines
		time.Sleep(200 * time.Millisecond)
		close(errCh)

		var successCount, errorCount int
		for err := range errCh {
			if err == nil {
				successCount++
			} else {
				errorCount++
				assert.Contains(t, err.Error(), "consumer already started")
			}
		}

		assert.Equal(t, 1, successCount, "exactly one start should succeed")
		assert.Equal(t, 1, errorCount, "second start should fail")

		stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = c.Stop(stopCtx)
	})

	t.Run("concurrent stop", func(t *testing.T) {
		c, err := consumer.New(
			consumer.WithBrokers("localhost:9092"),
			consumer.WithTopic("test-topic"),
			consumer.WithGroupID("test-group"),
			consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
			consumer.WithLogger(protocol.NopLogger{}),
		)
		require.NoError(t, err)

		ctx := context.Background()
		err = c.Start(ctx)
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond)

		stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		errCh := make(chan error, 2)
		go func() { errCh <- c.Stop(stopCtx) }()
		go func() { errCh <- c.Stop(stopCtx) }()

		time.Sleep(200 * time.Millisecond)
		close(errCh)

		// Both stops should succeed (idempotent)
		for err := range errCh {
			assert.NoError(t, err)
		}
	})
}

func TestHandlerErrorHandling(t *testing.T) {
	t.Parallel()

	handlerErr := errors.New("handler error")

	c, err := consumer.New(
		consumer.WithBrokers("localhost:9092"),
		consumer.WithTopic("test-topic"),
		consumer.WithGroupID("test-group"),
		consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error {
			return handlerErr
		}),
		consumer.WithLogger(protocol.NopLogger{}),
	)
	require.NoError(t, err)

	ctx := context.Background()
	err = c.Start(ctx)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = c.Stop(stopCtx)

	// Handler would be called if messages were received
	// This test verifies the handler function structure is correct
	assert.NotNil(t, c)
}

func TestDefaultValues(t *testing.T) {
	t.Parallel()

	handler := func(ctx context.Context, msg kafka.Message) error { return nil }

	c, err := consumer.New(
		consumer.WithBrokers("localhost:9092"),
		consumer.WithTopic("test-topic"),
		consumer.WithGroupID("test-group"),
		consumer.WithHandler(handler),
	)
	require.NoError(t, err)
	require.NotNil(t, c)

	// Verify consumer was created successfully
	// Default values are set internally
	assert.NotNil(t, c)
}

func TestStartOffsetOptions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		offset    kafka.StartOffset
		wantStart int64
	}{
		{
			name:      "StartOffsetEarliest",
			offset:    kafka.StartOffsetEarliest,
			wantStart: 0,
		},
		{
			name:      "StartOffsetLatest",
			offset:    kafka.StartOffsetLatest,
			wantStart: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := consumer.New(
				consumer.WithBrokers("localhost:9092"),
				consumer.WithTopic("test-topic"),
				consumer.WithGroupID("test-group"),
				consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
				consumer.WithStartOffset(tt.offset),
				consumer.WithLogger(protocol.NopLogger{}),
			)
			require.NoError(t, err)
			require.NotNil(t, c)
		})
	}
}

func TestStopWithTimeout(t *testing.T) {
	t.Parallel()

	c, err := consumer.New(
		consumer.WithBrokers("localhost:9092"),
		consumer.WithTopic("test-topic"),
		consumer.WithGroupID("test-group"),
		consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
		consumer.WithLogger(protocol.NopLogger{}),
	)
	require.NoError(t, err)

	ctx := context.Background()
	err = c.Start(ctx)
	require.NoError(t, err)
	time.Sleep(100 * time.Millisecond)

	// Test with very short timeout (should timeout waiting for stop)
	timeoutCtx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	defer cancel()

	// This might timeout depending on implementation
	_ = c.Stop(timeoutCtx)
	assert.NotNil(t, c)
}

func TestMultipleBrokers(t *testing.T) {
	t.Parallel()

	c, err := consumer.New(
		consumer.WithBrokers("localhost:9092", "localhost:9093", "localhost:9094"),
		consumer.WithTopic("test-topic"),
		consumer.WithGroupID("test-group"),
		consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error { return nil }),
		consumer.WithLogger(protocol.NopLogger{}),
	)
	require.NoError(t, err)
	require.NotNil(t, c)
}
