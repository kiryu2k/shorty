package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kiryu-dev/shorty/internal/config"
)

type HTTPServer struct {
	server     *http.Server
	mux        *chi.Mux
	urlService urlShortener
}

type urlShortener interface {
	MakeShort(context.Context, string) string
}

func New(cfg *config.HTTPServer, urlService urlShortener) *HTTPServer {
	httpServer := &HTTPServer{
		mux:        chi.NewMux(),
		urlService: urlService,
	}
	httpServer.setupMux()
	httpServer.server = &http.Server{
		Addr:         cfg.Address,
		Handler:      httpServer.mux,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return httpServer
}

func (s *HTTPServer) setupMux() {
	s.mux.Use(middleware.RequestID)
	s.mux.Route("/url", func(r chi.Router) {
		r.Post("/", nil)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("HELLO BACKEND WORLD!\n"))
		})
	})
}

func (s *HTTPServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}
