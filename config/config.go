package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	DefaultDBConfig = DBConfig{
		Port: "5432",
		User: "tms_user",
		DB:   "tms",
		Addr: "postgres",
	}
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Port     string
	Addr     string
	User     string
	Password string
	DB       string
}

func New() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("can't load config: %w", err)
	}

	dbcfg, err := parseDBCfg()
	if err != nil {
		return Config{}, fmt.Errorf("can't load postgres config: %w", err)
	}

	return Config{
		DB: dbcfg,
	}, nil
}

func parseDBCfg() (DBConfig, error) {
	cfg := DBConfig{
		Port:     os.Getenv("POSTGRES_PORT"),
		Addr:     os.Getenv("POSTGRES_ADDR"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DB:       os.Getenv("POSTGRES_DB"),
	}

	if cfg.Port == "" {
		cfg.Port = DefaultDBConfig.Port
	}
	if cfg.Addr == "" {
		cfg.Addr = DefaultDBConfig.Addr
	}
	if cfg.Password == "" {
		return cfg, errors.New("POSTGRES_PASSWORD is not set, set it in the .env file")
	}
	if cfg.User == "" {
		return cfg, errors.New("POSTGRES_USER is not set, set it in the .env file")
	}
	if cfg.DB == "" {
		return cfg, errors.New("POSTGRES_DB is not set, set it in the .env file")
	}
	return cfg, nil
}
