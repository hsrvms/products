package db

import (
	"context"
	"log"
	"products/pkg/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
	cfg  *config.Config
}

func New(cfg *config.Config) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDBConnString())
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return &Database{Pool: pool, cfg: cfg}, nil
}

func (d *Database) Close() {
	if d.Pool != nil {
		d.Pool.Close()
		log.Println("Database connection closed")
	}
}
