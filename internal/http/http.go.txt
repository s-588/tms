package http

import (
	"context"
	"fmt"
	"log/slog"
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

	var err error
	if s.Cfg.HTTPS {
		err = http.ListenAndServeTLS(":"+s.Cfg.Port, "server.crt", "server.key", s.mux)
	} else {
		err = http.ListenAndServe(":"+s.Cfg.Port, LogMiddleware(s.mux))
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
