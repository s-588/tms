package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/s-588/tms/config"
	"github.com/s-588/tms/db"
	"github.com/s-588/tms/static/template"
)

type Server struct {
	Port string
	Cfg  config.ServerConfig
	DB   db.DB
	mux  *http.ServeMux
}

func New(ctx context.Context, db db.DB, cfg config.ServerConfig) *Server {
	return &Server{
		Port: cfg.Port,
		DB:   db,
		Cfg:  cfg,
		mux:  http.NewServeMux(),
	}
}

func (s Server) Start() error {
	s.setHandlers()
	s.mux.HandleFunc("/", IndexHandler)

	var err error
	if s.Cfg.HTTPS {
		err = http.ListenAndServeTLS(":"+s.Cfg.Port, "server.crt", "server.key", s.mux)
	} else {
		err = http.ListenAndServe(":"+s.Cfg.Port, LogMiddleware(s.mux))
	}
	return fmt.Errorf("can't serve requests: %w", err)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	template.Index().Render(r.Context(), w)
}

func (s Server) Stop() {
	s.DB.Close()
}
