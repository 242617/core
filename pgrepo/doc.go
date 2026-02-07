// Package pgrepo provides PostgreSQL connection pool management with master/replica support.
//
// It implements the protocol.Lifecycle interface and offers graceful shutdown.
// Use Master() for write operations and Replica() for read operations.
//
// Example:
//
//	db, err := pgrepo.New(pgrepo.WithConfig(cfg.DB))
//	db.Start(ctx)
//
//	// Write
//	pgxscan.Select(ctx, db.Master(), &items, "SELECT * FROM items WHERE id = $1", id)
//
//	// Read
//	pgxscan.Select(ctx, db.Replica(ctx), &items, "SELECT * FROM items")
//
//	// Transaction
//	tx := db.BeginTx(ctx)
//	defer tx.Rollback(ctx)
//	tx.Commit(ctx)
package pgrepo
