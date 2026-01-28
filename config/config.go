package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	DefaultDBConfig = DBConfig{
		Port: "5432",
		User: "tms_user",
		DB:   "tms",
		Addr: "postgres",
	}
	DefaultServerConfig = ServerConfig{
		Port:  "8080",
		HTTPS: false,
	}
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	Logger LoggerConfig
}

type LoggerConfig struct {
	File  string
	Level string
}

type ServerConfig struct {
	Port  string
	HTTPS bool
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

	cfg, err := parseServerCfg()
	if err != nil {
		return Config{}, fmt.Errorf("can't load server config: %w", err)
	}

	logcfg, err := parseLoggerCfg()
	if err != nil {
		return Config{}, fmt.Errorf("can't load logger config: %w", err)
	}

	return Config{
		DB:     dbcfg,
		Server: cfg,
		Logger: logcfg,
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

func parseServerCfg() (ServerConfig, error) {
	cfg := ServerConfig{
		Port: os.Getenv("SERVER_PORT"),
	}
	https := os.Getenv("SERVER_HTTPS")
	if https == "1" || https == "true" || https == "yes" || https == "y" {
		cfg.HTTPS = true
	}

	if cfg.Port == "" {
		cfg.Port = DefaultServerConfig.Port
	}

	return cfg, nil
}

func parseLoggerCfg() (LoggerConfig, error) {
	cfg := LoggerConfig{
		Level: os.Getenv("LOG_LEVEL"),
		File:  os.Getenv("LOG_FILE"),
	}

	// Don't check cfg.Level because logger already handles uknown levels.

	if cfg.File == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg, fmt.Errorf("can't assemble default log file path: %w", err)
		}
		cfg.File = filepath.Join(home, ".local", "share", "tms", "tms.log")
	}

	return cfg, nil
}
