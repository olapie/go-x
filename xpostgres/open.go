package xpostgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.olapie.com/x/xlog"
)

type Config struct {
	// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	MaxConnLifetime time.Duration

	// MaxConnLifetimeJitter is the duration after MaxConnLifetime to randomly decide to close a connection.
	// This helps prevent all connections from being closed at the exact same time, starving the pool.
	MaxConnLifetimeJitter time.Duration

	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	MaxConnIdleTime time.Duration

	// MaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU().
	MaxConns int32

	// MinConns is the minimum size of the pool. After connection closes, the pool might dip below MinConns. A low
	// number of MinConns might mean the pool is empty after MaxConnLifetime until the health check has a chance
	// to create new connections.
	MinConns int32

	// HealthCheckPeriod is the duration between checks of the health of idle connections.
	HealthCheckPeriod time.Duration
}

func NewPool(ctx context.Context, connString string, config *Config) (*pgxpool.Pool, error) {
	logger := xlog.FromContext(ctx)
	logger.Debug("opening postgres: " + connString)
	//db, err := sql.Open("postgres", connString)
	//if err != nil {
	//	return nil, fmt.Errorf("open: %s, %w", connString, err)
	//}

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	if config != nil {
		poolConfig.MaxConns = config.MaxConns
		poolConfig.MinConns = config.MinConns
		poolConfig.MaxConnIdleTime = config.MaxConnIdleTime
		poolConfig.MaxConnLifetimeJitter = config.MaxConnLifetimeJitter
		poolConfig.MaxConnLifetime = config.MaxConnLifetime
		poolConfig.HealthCheckPeriod = config.HealthCheckPeriod
	}

	return pgxpool.NewWithConfig(ctx, poolConfig)
}

func Open(ctx context.Context, connString string, config *Config) (*sql.DB, error) {
	pool, err := NewPool(ctx, connString, config)
	if err != nil {
		return nil, fmt.Errorf("open: %s, %w", connString, err)
	}
	db := stdlib.OpenDBFromPool(pool)
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping: %s, %w", connString, err)
	}
	return db, nil
}

func MustOpen(ctx context.Context, connString string, config *Config) *sql.DB {
	return MustGet(Open(ctx, connString, config))
}

func OpenLocal(ctx context.Context) (*sql.DB, error) {
	if db, err := Open(ctx, NewConnStringBuilder().UseUnixDomainSocket(true).Build(), nil); err == nil {
		return db, nil
	}
	return Open(ctx, NewConnStringBuilder().Build(), nil)
}

func MustOpenLocal(ctx context.Context) *sql.DB {
	return MustGet(OpenLocal(ctx))
}

var connStringToDBCache sync.Map

type RepoFactory[T any] interface {
	Get(ctx context.Context, schema string) T
}

type NewRepoFunc[T any] func(ctx context.Context, db *sql.DB) T

type repoFactoryImpl[T any] struct {
	mu         sync.RWMutex
	cache      map[string]T
	connString string
	config     *Config
	fn         NewRepoFunc[T]
}

func NewRepoFactory[T any](connString string, config *Config, fn NewRepoFunc[T]) RepoFactory[T] {
	f := &repoFactoryImpl[T]{
		connString: connString,
		config:     config,
		cache:      make(map[string]T),
		fn:         fn,
	}
	return f
}

func (f *repoFactoryImpl[T]) Get(ctx context.Context, schema string) T {
	f.mu.RLock()
	repo, ok := f.cache[schema]
	f.mu.RUnlock()
	if ok {
		return repo
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	repo, ok = f.cache[schema]
	if ok {
		return repo
	}

	connString := SetParameterInConnString(f.connString, "search_path", schema)
	if dbVal, ok := connStringToDBCache.Load(connString); ok {
		repo = f.fn(ctx, dbVal.(*sql.DB))
	} else {
		db := MustOpen(ctx, connString, f.config)
		connStringToDBCache.Store(connString, db)
		repo = f.fn(ctx, db)
	}
	f.cache[schema] = repo
	return repo
}

// MustGet eliminates nil err and panics if err isn't nil
func MustGet[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
