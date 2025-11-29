package xpostgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.olapie.com/x/xlog"
	"go.olapie.com/x/xsync"
)

type PoolManager struct {
	mu         sync.Mutex
	m          xsync.Map[string, *pgxpool.Pool]
	connString string
	config     *Config
}

func NewPoolManager(connString string, config *Config) *PoolManager {
	m := &PoolManager{
		m:          xsync.Map[string, *pgxpool.Pool]{},
		connString: connString,
		config:     config,
	}
	return m
}

func (m *PoolManager) Get(ctx context.Context, schema string) (*pgxpool.Pool, error) {
	// TODO: search_path param not working after migration from self-hosted Postgres to supabase
	// so we set search_path explicitly
	connString := SetParameterInConnString(m.connString, "search_path", schema)
	m.mu.Lock()
	defer m.mu.Unlock()

	if pool, ok := m.m.Load(schema); ok {
		return pool, nil
	}

	pool, err := NewPool(ctx, connString, m.config)
	if err != nil {
		return nil, fmt.Errorf("new pool: %w", err)
	}
	if _, err = pool.Exec(ctx, "set search_path to "+schema); err != nil {
		return nil, fmt.Errorf("set search_path to %s: %w", schema, err)
	}
	xlog.FromContext(ctx).InfoContext(ctx, "new pool successfully for schema "+schema)
	m.m.Store(schema, pool)
	return pool, nil
}
