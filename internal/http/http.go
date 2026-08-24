package http

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"mime"
	"net"
	"net/http"

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
		Port:    cfg.Port,
		Cfg:     cfg,
		mux:     http.NewServeMux(),
		Handler: handler.NewHandler(db),
	}
}

func (s Server) Start() error {
	s.setHandlers()
	s.mux.HandleFunc("/", IndexHandler)
	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	addr := "0.0.0.0:" + s.Cfg.Port
	log.Printf("listening on %s (HTTPS=%v)", addr, s.Cfg.HTTPS) // ← add this so you can see it in the logs

	ln, err := net.Listen("tcp4", addr) // force IPv4 socket
	if err != nil {
		return fmt.Errorf("can't listen on %s: %w", addr, err)
	}

	if s.Cfg.HTTPS {
		err = http.ServeTLS(ln, s.mux, "server.crt", "server.key")
	} else {
		err = http.Serve(ln, LogMiddleware(s.mux))
	}
	return fmt.Errorf("can't serve requests: %w", err)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug("serving home page")
	ui.Index().Render(r.Context(), w)
}

func (s Server) Stop() {
	s.Handler.DB.Close()
}
