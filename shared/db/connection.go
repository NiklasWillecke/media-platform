package db

import (
	"context"
	"fmt"

	db "streaming-platform/shared/db/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	Q    *db.Queries
}

func NewDB(url string) *DB {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		fmt.Println(err.Error())
	}

	return &DB{
		pool: pool,
		Q:    db.New(pool),
	}
}

func (db *DB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
