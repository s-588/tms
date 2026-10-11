package http

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"net"
	"net/http"
	"time"

	"github.com/s-588/tms/internal/config"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/http/handler"
	"github.com/s-588/tms/internal/ui"
)

type Server struct {
	Port    string
	Cfg     config.ServerConfig
	Handler handler.Handler
	mux     *http.ServeMux
}

func New(ctx context.Context, db db.DB, cfg config.ServerConfig) *Server {
	if err := mime.AddExtensionType(".css", "text/css"); err != nil {
		slog.Warn("set .css mime type", "error", err)
	}
	if err := mime.AddExtensionType(".js", "application/javascript"); err != nil {
		slog.Warn("set .js mime type: %w", "error", err)
	}
	return &Server{
		Port:    cfg.HTTPPort,
		Cfg:     cfg,
		mux:     http.NewServeMux(),
		Handler: handler.NewHandler(db),
	}
}

func (s Server) Start() error {
	s.setHandlers()
	s.mux.HandleFunc("/", IndexHandler)
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	addr := "0.0.0.0:" + s.Cfg.HTTPPort

	ln, err := net.Listen("tcp4", addr) // force IPv4 socket
	if err != nil {
		return fmt.Errorf("can't listen on %s: %w", addr, err)
	}

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      LogMiddleware(s.mux),
		ReadTimeout:  time.Duration(s.Cfg.HTTPTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Cfg.HTTPTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Cfg.HTTPTimeout) * time.Second,
	}

	if s.Cfg.HTTPS {
		err = httpServer.ServeTLS(ln, "server.crt", "server.key")
	} else {
		err = httpServer.Serve(ln)
	}
	return err
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("serving home page")
	err := ui.Index().Render(r.Context(), w)
	if err != nil {
		slog.Error("can't render response", "error", err)
	}
}

func (s Server) Stop() {
	s.Handler.DB.Close()
}
