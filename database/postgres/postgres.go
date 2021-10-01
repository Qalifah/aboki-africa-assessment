package postgres

import (
	"context"
	"fmt"
	"time"
	"log"
	"strings"

	core "github.com/Qalifah/aboki-africa-assessment"
	"github.com/Qalifah/aboki-africa-assessment/config"

	pool "github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgconn"
	"github.com/pkg/errors"
)

var defaultOptions = pgx.TxOptions{
	IsoLevel:       pgx.ReadCommitted,
	DeferrableMode: pgx.NotDeferrable,
	AccessMode:     pgx.ReadWrite,
}

// Tx represents a database transaction
type Tx interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)

	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row

	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)

	Commit(ctx context.Context) error

	Rollback(ctx context.Context) error
}

type Client struct {
	pool *pool.Pool
}

// New Returns a new database initialized with credentials from config
func New(ctx context.Context, config *config.PostgresConfig) (*Client, error) {
	const format = "postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_max_conns=%d"
	uri := fmt.Sprintf(format, config.Username, config.Password, config.Host, config.Port, config.Database, config.MaxConn)

	cfg, err := pool.ParseConfig(uri)
	if err != nil {
		log.Fatalf("failed to parse pgx config: %v", err)
		return nil, err
	}

	cfg.ConnConfig.ConnectTimeout = time.Minute

	pool, err := pool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("pgx pool failed to connect: %v", err)
		return nil, err 
	}

	return &Client{
		pool: pool,
	}, nil
}

func (c *Client) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rs, err := c.pool.Query(ctx, query, args...)
	return rs, err
}

func (c *Client) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	row := c.pool.QueryRow(ctx, query, args...)
	return row
}

func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	_, err := c.pool.Exec(ctx, query, args...)
	return nil, err
}

func (c *Client) Commit(ctx context.Context) error {
	return nil
}

func (c *Client) Rollback(ctx context.Context) error {
	return nil
}


func (c *Client) BeginTx() (pgx.Tx, error) {
	tx, err := c.pool.BeginTx(context.Background(), defaultOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin new transaction")
	}
	return tx, nil
}


func (c *Client) GetTx(ctx context.Context) (Tx, error) {
	tx := ctx.Value(core.TxContextKey)
	if tx != nil {
		return tx.(Tx), nil
	}
	return c, nil
}

func IsDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}