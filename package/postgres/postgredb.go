package postgres

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	defaultMaxPoolSize        = 1
	defaultConnectionAttempts = 10
	defaultConnectionTimeout  = time.Second
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type PostgreDB struct {
	maxPoolSize        int
	connectionAttempts int
	connectionTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Pool    PgxPool
}

func New(host, port, database, username, password string, opts ...Option) (*PostgreDB, error) {
	pg := &PostgreDB{
		maxPoolSize:        defaultMaxPoolSize,
		connectionAttempts: defaultConnectionAttempts,
		connectionTimeout:  defaultConnectionTimeout,
	}

	for _, option := range opts {
		option(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres properties from the config due to error: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connectionAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Trying to connect to the postgresql database, attempts left: %d", pg.connectionAttempts)
		time.Sleep(pg.connectionTimeout)
		pg.connectionAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql database due to error: %w", err)
	}

	return pg, nil
}

func (pg *PostgreDB) Close() {
	if pg.Pool != nil {
		pg.Pool.Close()
	}
}
