package cockroach

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/audryus/2dpoint.site/internal/config"
	"github.com/audryus/2dpoint.site/pkg/logger"
	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Cockroach interface {
	Acquire() *pgxpool.Pool
	ExecuteTx(ctx context.Context, fn func(pgx.Tx) error) error
	Close()
}

type cockroach struct {
	pool *pgxpool.Pool
}

func New(conf config.Config, l logger.Log) (Cockroach, error) {
	pool, err := pgxpool.NewWithConfig(context.Background(), databaseConfig(conf.Cockroach))

	if err != nil {
		return nil, err
	}

	l.Info("Cockroach connected")

	return cockroach{
		pool: pool,
	}, nil
}

func (c cockroach) Close() {
	c.pool.Close()
}

func (c cockroach) Acquire() *pgxpool.Pool {
	return c.pool
}

func (c cockroach) ExecuteTx(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	adp := pgxTxAdapter{tx}

	err = crdb.ExecuteInTx(ctx, adp, func() error { return fn(tx) })
	return err
}

type pgxTxAdapter struct {
	tx pgx.Tx
}

var _ crdb.Tx = pgxTxAdapter{}

func (tx pgxTxAdapter) Commit(ctx context.Context) error {
	return tx.tx.Commit(ctx)
}

func (tx pgxTxAdapter) Rollback(ctx context.Context) error {
	return tx.tx.Rollback(ctx)
}

// Exec is part of the crdb.Tx interface.
func (tx pgxTxAdapter) Exec(ctx context.Context, q string, args ...interface{}) error {
	_, err := tx.tx.Exec(ctx, q, args...)
	return err
}

func databaseConfig(conf config.Cockroach) *pgxpool.Config {
	const defaultMaxConns = int32(2)
	const defaultMinConns = int32(1)
	const defaultMaxConnLifetime = time.Minute * 15
	const defaultMaxConnIdleTime = time.Minute * 7
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	url := fmt.Sprintf(conf.Url, conf.Port, conf.Database)

	dbConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}
