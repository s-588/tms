package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/s-588/tms/config"
	"github.com/s-588/tms/db"
	"github.com/s-588/tms/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("can't start app", "error", err)
		return
	}
	slog.Info("config successfully parsed", "config", cfg)
	closeLogFile, err := SetupSLog(cfg.Logger)
	if err != nil {
		slog.Error("can't start app", "error", err)
		return
	}
	defer closeLogFile()
	slog.Info("slog configured")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	dbConn, err := db.New(ctx, cfg.DB)
	if err != nil {
		slog.Error("can't start app", "error", err)
		return
	}
	slog.Info("database connected")

	s := http.New(context.Background(), dbConn, cfg.Server)
	slog.Info("server ready to start")

	slog.Info("starting server")
	if err := s.Start(); err != nil {
		slog.Error("can't start server", "error", err)
		return
	}
	slog.Info("shuting down app")
}
