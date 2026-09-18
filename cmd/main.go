package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/s-588/tms/internal/config"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/grpc"
	"github.com/s-588/tms/internal/http"
	"github.com/s-588/tms/internal/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("can't start app", "error", err)
		return
	}
	closeLogFile, err := logger.SetupSLog(cfg.Logger)
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

	grpcSrv := grpc.NewServer(dbConn)
	go func() {
		slog.Info("grpc server started")
		if err := grpcSrv.Run(cfg.Server.GRPCPort); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	s := http.New(context.Background(), dbConn, cfg.Server)
	slog.Info("server ready to start")

	slog.Info("starting server")
	if err := s.Start(); err != nil {
		slog.Error("can't start server", "error", err)
		return
	}
	slog.Info("shuting down app")
}
