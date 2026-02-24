package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-588/tms/internal/config"
	"github.com/s-588/tms/internal/db/generated"
)

type DB struct {
	queries *generated.Queries
	pool    *pgxpool.Pool
	cfg     config.DBConfig
}

func New(ctx context.Context, cfg config.DBConfig) (DB, error) {
	pool, err := pgxpool.New(ctx, getConnStr(cfg))
	if err != nil {
		return DB{}, fmt.Errorf("can't create database connection: %w", err)
	}

	quieries := generated.New(pool)
	return DB{
		queries: quieries,
		cfg:     cfg,
		pool:    pool,
	}, nil
}

func getConnStr(cfg config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)
}

func (db DB) Close() {
	db.pool.Close()
}
