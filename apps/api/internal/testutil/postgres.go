// Author: Caden Lund
// Created: 4/11/2026
// Last updated: 4/11/2026
// Notes:
// - helpers for testing
// - one returns pooler from container
// - other destroys the in memory container

package testutil

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var Container *postgres.PostgresContainer

func Setup() (*pgxpool.Pool, error) {
	//1. Setup container
	var err error
	Container, err = postgres.Run(context.Background(), "postgres:16-alpine",
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	//2. Get the connection string
	connStr, err := Container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		return nil, err
	}

	//3. Get the pool
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	//4. IMPORTANT: must convert from pool to *sql.DB which is what goose expects
	conn := stdlib.OpenDBFromPool(pool)

	//5. Set sql flavor & run migrations
	goose.SetDialect("postgres")
	err = goose.Up(conn, "../../migrations")
	if err != nil {
		return nil, err
	}

	//4. Return
	return pool, nil
}

func Cleanup() error {
	err := testcontainers.TerminateContainer(Container)
	if err != nil {
		return err
	}

	return nil
}

func WithTx(t *testing.T, pool *pgxpool.Pool) pgx.Tx {
	tx, err := pool.Begin(context.Background())
	require.NoError(t, err)
	t.Cleanup(func() { tx.Rollback(context.Background()) })
	return tx
}
