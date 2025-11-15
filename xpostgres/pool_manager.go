package xpostgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.olapie.com/x/xlog"
	"go.olapie.com/x/xsync"
)

type DBTX interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
}

type PoolManager struct {
	mu         sync.Mutex
	m          xsync.Map[string, DBTX]
	connString string
	config     *Config
}

func NewPoolManager(ctx context.Context, connString string, config *Config) *PoolManager {
	m := &PoolManager{
		m:          xsync.Map[string, DBTX]{},
		connString: connString,
		config:     config,
	}

	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to postgres database: %v", err))
	}
	defer func() {
		_ = conn.Close(ctx)
	}()
	return m
}

func (m *PoolManager) Get(ctx context.Context, schema string) (DBTX, error) {
	// TODO: search_path param not working after migration from self-hosted Postgres to supabase
	// so we set search_path explicitly
	connString := SetParameterInConnString(m.connString, "search_path", schema)
	m.mu.Lock()
	defer m.mu.Unlock()
	pool, err := NewPool(ctx, connString, m.config)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}
	if _, err = pool.Exec(ctx, string("set search_path to "+schema)); err != nil {
		return nil, fmt.Errorf("set search_path to %s: %w", schema, err)
	}
	xlog.FromContext(ctx).InfoContext(ctx, string("new pool successfully for schema "+schema))
	m.m.Store(schema, pool)
	return pool, nil
}
