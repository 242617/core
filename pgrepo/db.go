package pgrepo

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/242617/core/protocol"
)

// DB wraps a master pool and optionally multiple replica pools.
// It provides graceful shutdown.
type DB struct {
	master   *pgxpool.Pool
	replicas []*pgxpool.Pool

	// Logger
	log protocol.Logger

	// Configuration
	cfg Config

	// Lifecycle
	ctx     context.Context
	cancel  context.CancelFunc
	started bool
	mu      sync.RWMutex
}

// New creates a new DB instance with the given configuration and options.
// The DB must be started with Start() before use.
//
// Example:
//
//	db, err := dbrepo.New(cfg, dbrepo.WithLogger(log))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if err := db.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
func New(options ...Option) (*DB, error) {
	// Create DB instance
	db := &DB{}

	// Apply options
	for _, option := range append(defaults(), options...) {
		if err := option(db); err != nil {
			return nil, errors.Wrap(err, "apply option")
		}
	}

	// Validate configuration
	if err := db.cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid config")
	}

	// Create context for lifecycle management
	db.ctx, db.cancel = context.WithCancel(context.Background())

	return db, nil
}

// Start initializes the database connection pools and starts background services.
// This method must be called before using the DB instance.
func (db *DB) Start(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.started {
		return nil
	}

	// Create master pool config
	masterCfg, err := pgxpool.ParseConfig(db.cfg.String())
	if err != nil {
		return errors.Wrap(err, "parse master config")
	}

	// Configure master pool
	masterCfg.MaxConnLifetime = db.cfg.ConnMaxLifeTime
	masterCfg.MaxConnIdleTime = db.cfg.ConnMaxIdleTime
	masterCfg.MinConns = int32(db.cfg.MinConns)
	masterCfg.MaxConns = int32(db.cfg.MaxConns)

	// Create master pool
	db.master, err = pgxpool.NewWithConfig(ctx, masterCfg)
	if err != nil {
		return errors.Wrap(err, "create master pool")
	}

	// Ping master to verify connectivity
	if err := db.master.Ping(ctx); err != nil {
		db.master.Close()
		return errors.Wrap(err, "ping master")
	}

	db.log.Info(ctx, "master pool started", "dsn", db.cfg.RedactedDSN())

	// Create replica pools if configured
	if len(db.cfg.Replicas) > 0 {
		db.replicas = make([]*pgxpool.Pool, 0, len(db.cfg.Replicas))

		for i, replicaCfg := range db.cfg.Replicas {
			cfg, err := pgxpool.ParseConfig(replicaCfg.String())
			if err != nil {
				db.log.Warn(ctx, "parse replica config", "index", i, "error", err)
				continue
			}

			cfg.MaxConnLifetime = replicaCfg.ConnMaxLifeTime
			cfg.MaxConnIdleTime = replicaCfg.ConnMaxIdleTime
			cfg.MinConns = int32(replicaCfg.MinConns)
			cfg.MaxConns = int32(replicaCfg.MaxConns)

			replica, err := pgxpool.NewWithConfig(ctx, cfg)
			if err != nil {
				db.log.Warn(ctx, "create replica pool", "index", i, "error", err)
				continue
			}

			if err := replica.Ping(ctx); err != nil {
				replica.Close()
				db.log.Warn(ctx, "ping replica", "index", i, "error", err)
				continue
			}

			db.replicas = append(db.replicas, replica)
			db.log.Info(ctx, "replica pool started", "index", i, "dsn", replicaCfg.RedactedDSN())
		}
	}

	db.started = true
	db.log.Info(ctx, "db started", "replicas", len(db.replicas))
	return nil
}

// Stop gracefully shuts down the database connection pools.
// It waits for in-flight queries to complete or until shutdown timeout.
func (db *DB) Stop(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if !db.started {
		return nil
	}

	// Cancel context to stop background goroutines
	db.cancel()

	// Shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, db.cfg.ShutdownTimeout)
	defer cancel()

	// Close master pool
	if db.master != nil {
		db.log.Info(shutdownCtx, "closing master pool")
		db.master.Close()
	}

	// Close replica pools
	for i, replica := range db.replicas {
		if replica != nil {
			db.log.Info(shutdownCtx, "closing replica pool", "index", i)
			replica.Close()
		}
	}

	db.started = false
	db.log.Info(ctx, "db stopped")
	return nil
}

// Ping verifies connectivity to master and all replicas.
func (db *DB) Ping(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if !db.started {
		return ErrDatabaseNotStarted
	}

	// Ping master
	if err := db.master.Ping(ctx); err != nil {
		return errors.Wrap(err, "ping master")
	}

	// Ping replicas
	for i, replica := range db.replicas {
		if err := replica.Ping(ctx); err != nil {
			return errors.Wrapf(err, "ping replica[%d]", i)
		}
	}

	return nil
}

// Master returns the master pool for write operations.
func (db *DB) Master() *pgxpool.Pool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.master
}

// Replica returns a replica pool for read operations.
// If no replicas are available, returns the master pool as fallback.
func (db *DB) Replica(ctx context.Context) *pgxpool.Pool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if !db.started {
		return nil
	}

	// If no replicas configured, return master
	if len(db.replicas) == 0 {
		return db.master
	}

	// Return first replica
	return db.replicas[0]
}

// Config returns the DB configuration.
func (db *DB) Config() Config {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.cfg
}

// IsStarted returns true if the DB has been started.
func (db *DB) IsStarted() bool {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.started
}
