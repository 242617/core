package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

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
	log, err := logger.New(
		logger.WithLevel(logger.LevelDebug),
		logger.WithDevelopmentConfig(),
	)
	die(err)

	ctx := context.Background()

	start := time.Now()
	log.Debug(ctx, "start")
	defer func() { log.Debug(ctx, "stop", "in", time.Since(start)) }()

	var cfg struct {
		DB               pgrepo.Config   `yaml:"db"`
		MessagesConsumer consumer.Config `yaml:"messages_consumer"`
		MessagesProducer producer.Config `yaml:"messages_producer"`
	}
	die(config.New().With(file.YAML("config.yaml")).Scan(&cfg))

	// db, err := pgrepo.New(pgrepo.WithLogger(log.New("pgrepo")), pgrepo.WithConfig(cfg.DB))
	// die(err)

	producer, err := producer.New(
		producer.WithLogger(log.New("producer")),
		producer.WithConfig(cfg.MessagesProducer),
	)
	die(err)

	consumer, err := consumer.New(
		consumer.WithLogger(log.New("consumer")),
		consumer.WithConfig(cfg.MessagesConsumer),
		// consumer.WithGroupID(fmt.Sprintf("temporary_group_%v", rand.Int())),
		consumer.WithStartOffset(kafka.StartOffsetEarliest),
		consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error {
			log.Info(ctx, "incoming message",
				"key", string(msg.Key),
				"value", string(msg.Value),
			)
			time.Sleep(time.Second)
			log.Info(ctx, "outcoming message",
				"key", string(msg.Key),
				"value", string(msg.Value),
			)
			return producer.Produce(ctx, msg)
		}),
	)
	die(err)

	app, err := application.New(
		application.WithLogger(log.New("application")),
		application.WithName("main"),
		application.WithComponents(
			// application.NewLifecycleComponent("db", db),
			application.NewLifecycleComponent("consumer", consumer),
			application.NewLifecycleComponent("producer", producer),
		),
	)
	die(err)

	go func() {
		time.Sleep(time.Second)
		log.Debug(ctx, "start")
		msg := kafka.Message{
			Key:   []byte("sample key"),
			Value: []byte("sample value"),
		}
		die(producer.Produce(ctx, msg))
	}()

	die(app.Run(ctx))
}

func die(args ...any) {
	if len(args) == 0 {
		return
	}
	if err, ok := args[len(args)-1].(error); ok && err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: %s", file, line, err.Error())
		os.Exit(1)
	}
}
