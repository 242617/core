package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/242617/core/protocol"
)

// Application manages component lifecycle with graceful shutdown and timeout handling.
type Application struct {
	log          protocol.Logger
	startTimeout time.Duration
	stopTimeout  time.Duration
	components   Components
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	started      bool
	startedMu    sync.RWMutex
	stopCh       chan struct{} // For programmatic shutdown
	hostname     string
	name         string
}

// New creates Application with defaults and custom options.
func New(options ...Option) (*Application, error) {
	ctx, cancel := context.WithCancel(context.Background())

	app := &Application{
		ctx:    ctx,
		cancel: cancel,
		stopCh: make(chan struct{}, 1),
	}

	for _, option := range append(defaults(), options...) {
		if err := option(app); err != nil {
			return nil, fmt.Errorf("apply option: %w", err)
		}
	}

	for i, c := range app.components {
		if c == nil {
			return nil, fmt.Errorf("component at index %d is nil", i)
		}
	}

	if app.log == nil {
		return nil, errors.New("empty log")
	}
	if app.startTimeout == 0 {
		return nil, errors.New("empty start timeout")
	}
	if app.stopTimeout == 0 {
		return nil, errors.New("empty stop timeout")
	}
	if app.hostname == "" {
		return nil, errors.New("empty hostname")
	}
	if app.name == "" {
		return nil, errors.New("empty name")
	}

	return app, nil
}

// Exit triggers graceful shutdown. Safe to call multiple times.
func (a *Application) Exit() {
	select {
	case a.stopCh <- struct{}{}:
	default:
	}
}
