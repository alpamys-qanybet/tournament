package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

var connectionPool *pgxpool.Pool

func ConnectionPool() (*pgxpool.Pool, error) { // just get connection
	if connectionPool == nil {
		return nil, errors.New("postgres db connection pool not connected")
	}
	return connectionPool, nil
}

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) { // try to get connection or connect and get
	if connectionPool != nil {
		return connectionPool, nil
	}

	dbpool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	connectionPool = dbpool
	connectionPool.Ping(ctx)
	return dbpool, nil
}
