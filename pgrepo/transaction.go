package pgrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// txKey is the context key for storing a transaction.
type txKey struct{}

// WithTx executes a function within a transaction.
// The transaction is automatically committed if fn returns nil.
// If fn returns an error, the transaction is rolled back.
// The transaction is passed via context, so it can be retrieved using GetTx.
//
// This allows for cross-repository transactions by passing the same context
// to multiple repository methods.
func WithTx(ctx context.Context, db *DB, fn func(ctx context.Context, tx pgx.Tx) error) error {
	if db == nil {
		return errors.New("db is required")
	}
	if db.master == nil {
		return ErrDatabaseNotStarted
	}

	tx, err := db.master.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	// Rollback is safe to call multiple times
	defer tx.Rollback(ctx)

	// Pass transaction via context
	ctxWithTx := context.WithValue(ctx, txKey{}, tx)

	// Execute user function
	if err := fn(ctxWithTx, tx); err != nil {
		return errors.Wrap(err, "execute transaction")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "commit transaction")
	}

	return nil
}

// GetTx retrieves the transaction from context.
// Returns the transaction and true if found, nil and false otherwise.
func GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

// Exec executes a query without returning any rows.
// Uses transaction if available in context, otherwise uses the provided pool.
func Exec(ctx context.Context, pool *pgxpool.Pool, query string, args ...any) (int64, error) {
	// Check if transaction exists in context
	if tx, ok := GetTx(ctx); ok {
		result, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected(), nil
	}

	// Use pool directly
	result, err := pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

// QueryRow executes a query that is expected to return at most one row.
// Uses transaction if available in context, otherwise uses the provided pool.
func QueryRow(ctx context.Context, pool *pgxpool.Pool, query string, args []any, dest ...any) error {
	// Check if transaction exists in context
	if tx, ok := GetTx(ctx); ok {
		return tx.QueryRow(ctx, query, args...).Scan(dest...)
	}

	// Use pool directly
	return pool.QueryRow(ctx, query, args...).Scan(dest...)
}

// Query executes a query that returns multiple rows.
// Uses transaction if available in context, otherwise uses the provided pool.
func Query(ctx context.Context, pool *pgxpool.Pool, query string, args []any, fn func(rows pgx.Rows) error) error {
	// Check if transaction exists in context
	if tx, ok := GetTx(ctx); ok {
		rows, err := tx.Query(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		return fn(rows)
	}

	// Use pool directly
	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return fn(rows)
}
