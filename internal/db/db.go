package db

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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

//go:embed migrations/*.sql
var embeddedMigrations embed.FS

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
	db := DB{
		queries: quieries,
		cfg:     cfg,
		pool:    pool,
	}
	db.initDB()
	return db, nil
}

func getConnStr(cfg config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) initDB() error {
	goose.SetBaseFS(embeddedMigrations)
	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}

	dbFromPool := stdlib.OpenDBFromPool(db.pool)

	if err := goose.Up(dbFromPool, "migrations"); err != nil {
		return fmt.Errorf("execute migrations: %w", err)
	}
	return nil
}
