package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/s-588/tms/config"
	"github.com/s-588/tms/db/generated"
)

type DB struct {
	queries *generated.Queries
	cfg     config.DBConfig
}

func New(ctx context.Context, cfg config.DBConfig) (*DB, error) {
	conn, err := pgx.Connect(ctx, getConnStr(cfg))
	if err != nil {
		return nil, fmt.Errorf("can't create database connection: %w", err)
	}

	quieries := generated.New(conn)
	return &DB{
		queries: quieries,
		cfg:     cfg,
	}, nil
}

func getConnStr(cfg config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)
}
