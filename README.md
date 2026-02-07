# Core Library Documentation

Production-grade Go infrastructure components for microservices and distributed systems.

---

## Logger

Structured logging with context propagation.

```go
import "github.com/242617/core/logger"

log, _ := logger.New(logger.WithConfig(logger.DevelopmentConfig))
serviceLog := log.New("auth-service")

serviceLog.Info(ctx, "request", "method", "GET")
serviceLog.Error(ctx, "failed", "err", err)
log.SetLevel(logger.LevelDebug)
```

### Config

```go
type Config struct {
    Level, Encoding string
    Colorize        bool
}
// Predefined: ProductionConfig, DevelopmentConfig
```

---

## Config

Multi-source loading: defaults → env → file.

```go
import "github.com/242617/core/config"
import "github.com/242617/core/config/source/file"

type Config struct {
    DB       pgrepo.Config   `yaml:"db"`
    Messages consumer.Config `yaml:"messages"`
    Timeout time.Duration    `yaml:"timeout" default:"30s"`
}

var cfg Config
config.New().With(file.YAML("config.yaml")).Scan(&cfg)
```

---

## Protocol

Common interfaces.

```go
type Logger interface {
    Debug, Info, Warn, Error(ctx context.Context, msg string, args ...any)
}

type Lifecycle interface {
    Start(context.Context) error
    Stop(context.Context) error
}

protocol.NopLogger{} // implements protocol.Logger interface
```

---

## Application

Component lifecycle management with graceful shutdown.

```go
import "github.com/242617/core/application"

app, _ := application.New(
    application.WithName("my-service"),
    application.WithLogger(log.New("my-application")),
    application.WithComponents(
        application.NewLifecycleComponent("db", db),
        application.NewLifecycleComponent("kafka-consumer", consumer),
    ),
)

app.Run(ctx)

go func() {
    <-shutdownSignal
    app.Exit()
}()
```

---

## PostgreSQL Repository

Connection pool management with master/replica support.

```go
import "github.com/242617/core/pgrepo"

db, _ := pgrepo.New(pgrepo.WithConfig(cfg.DB), pgrepo.WithLogger(log.New("pgrepo")))
db.Start(ctx)

pgxscan.Select(ctx, db.Master(), &items, "SELECT * FROM items")
pgxscan.Select(ctx, db.Replica(ctx), &items, "SELECT * FROM items")

db.Stop(ctx)
```

### Config

```go
type Config struct {
    Host, Port, Schema, User, Password, Name string/int
    SSL bool
    ConnMaxLifeTime, ConnMaxIdleTime, ShutdownTimeout time.Duration
    MinConns, MaxConns int
    Replicas []Config
}
```

### Transaction

```go
tx := db.BeginTx(ctx)
defer tx.Rollback(ctx)
pgxscan.Select(ctx, tx, &items, query)
tx.Commit(ctx)
```

---

## Kafka

Producer and consumer with consumer group support.

### Producer

```go
import "github.com/242617/core/kafka/producer"

producer, _ := producer.New(producer.WithConfig(cfg.Producer))

msg := kafka.Message{Key: []byte("key"), Value: []byte("value")}
producer.Produce(ctx, msg)
```

### Consumer

```go
import "github.com/242617/core/kafka/consumer"
import "github.com/242617/core/kafka"

consumer, _ := consumer.New(
    consumer.WithConfig(cfg.Consumer),
    consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error {
        log.Info(ctx, "processing", "key", string(msg.Key))
        return nil
    }),
)

consumer.Start(ctx)
consumer.Stop(ctx)
```

### Config

```go
type Config struct {
    Brokers []string
    Topic   string
    GroupID string // consumer only
}

kafka.StartOffsetLatest, kafka.StartOffsetEarliest
```

---

## Mock Generation

Generate mocks with [mockery](https://github.com/vektra/mockery).

```bash
mockery --all
```

### Example

```go
import "github.com/242617/core/mocks"

mockLogger := mocks.Logger{}
mockLogger.On("Info", mock.Anything, "test", mock.Anything).Return()
svc := NewService(mockLogger)
svc.DoSomething(ctx)
mockLogger.AssertExpectations(t)
```

### Available

Logger, Lifecycle, Component, ConfigEngine, ConfigSource, Handler, KafkaProducer

---

## Complete Example

```go
package main

import (
    "github.com/242617/core/application"
    "github.com/242617/core/config"
    "github.com/242617/core/config/source/file"
    "github.com/242617/core/kafka"
    "github.com/242617/core/kafka/consumer"
    "github.com/242617/core/kafka/producer"
    "github.com/242617/core/logger"
    "github.com/242617/core/pgrepo"
)

func main() {
    log, _ := logger.New()

    var cfg struct {
        DB       pgrepo.Config   `yaml:"db"`
        Producer producer.Config `yaml:"producer"`
        Consumer consumer.Config `yaml:"consumer"`
    }
    config.New().With(file.YAML("config.yaml")).Scan(&cfg)

    db, _ := pgrepo.New(pgrepo.WithConfig(cfg.DB), pgrepo.WithLogger(log.New("db")))
    producer, _ := producer.New(producer.WithConfig(cfg.Producer), producer.WithLogger(log.New("producer")))

    consumer, _ := consumer.New(
        consumer.WithConfig(cfg.Consumer),
        consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error {
            return producer.Produce(ctx, msg)
        }),
    )

    app, _ := application.New(
        application.WithName("my-service"),
        application.WithLogger(log.New("application")),
        application.WithComponents(
            application.NewLifecycleComponent("db", db),
            application.NewLifecycleComponent("producer", producer),
            application.NewLifecycleComponent("consumer", consumer),
        ),
    )

    app.Run(context.Background())
}
```

---

## Pipeline

Fluent workflow execution with error handling and cleanup.

```go
import "github.com/242617/core/pipeline"

errCh := make(chan error)
go pipeline.New(ctx).
    Before(func() { fmt.Println("setup") }).
    Then(func(ctx context.Context) error {
        // main logic
        return nil
    }).
    Else(func(ctx context.Context) error {
        // fallback on error
        return nil
    }).
    After(func() { fmt.Println("cleanup") }).
    Run(func(err error) { errCh <- err })
```

**Execution order:** `Before` → `Then`/`Else` → `Error`/`NoError` → `After`

---

## Request ID

Context-based request ID propagation for tracing.

```go
import "github.com/242617/core/request_id"

// Set in middleware
ctx = request_id.ContextWithRequestID(ctx, requestID)

// Retrieve anywhere
requestID := request_id.RequestIDFromContext(ctx)
```

---

## Best Practices

1. Always pass context through all layers
2. Use named loggers (`log.New("component")`) for filtering
3. Handle all errors - never ignore returns
4. Use defaults with `default:"value"` tags
5. Implement protocol.Lifecycle for long-running components
6. Generate mocks with `mockery --all` before testing
7. Validate configurations early
8. Use structured logging with key-value pairs
