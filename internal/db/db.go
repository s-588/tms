package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-588/tms/internal/config"
	"github.com/s-588/tms/internal/db/generated"
)

var (
	ErrIncorrectPhone = errors.New("phone have incorrect format")
    ErrDuplicateEmail = errors.New("email already exists")
    ErrDuplicatePhone = errors.New("phone already exists")

    ErrDuplicateLicense = errors.New("license plate already exists")

    ErrDuplicatePrice = errors.New("price configuration (cargo type, weight, distance) already exists")
	
	ErrDuplicateNodeAddress = errors.New("node with this address already exists")
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
