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
		DB          pgrepo.Config   `yaml:"db"`
		RawMessages consumer.Config `yaml:"raw_messages"`
	}
	die(config.New().With(file.YAML("config.yaml")).Scan(&cfg))

	db, err := pgrepo.New(pgrepo.WithLogger(log.New("pgrepo")), pgrepo.WithConfig(cfg.DB))
	die(err)

	rawMessages, err := consumer.New(
		consumer.WithLogger(log.New("raw messages")),
		consumer.WithConfig(cfg.RawMessages),
		// consumer.WithGroupID(fmt.Sprintf("temporary_group_%v", rand.Int())),
		consumer.WithStartOffset(kafka.StartOffsetEarliest),
		consumer.WithHandler(func(ctx context.Context, msg kafka.Message) error {
			log.Info(ctx, "message",
				"key", string(msg.Key),
				"value", string(msg.Value),
			)
			return nil
		}),
	)
	die(err)

	app, err := application.New(
		application.WithLogger(log.New("application")),
		application.WithName("main"),
		application.WithComponents(
			application.NewLifecycleComponent("db", db),
			application.NewLifecycleComponent("raw messages", rawMessages),
		),
	)
	die(err)

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
