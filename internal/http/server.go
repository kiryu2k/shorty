package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kiryu-dev/shorty/internal/config"
	"github.com/kiryu-dev/shorty/internal/http/handlers"
	"github.com/kiryu-dev/shorty/internal/http/validator"
)

type HTTPServer struct {
	server     *http.Server
	mux        *chi.Mux
	validator  *validator.RequestValidator
	urlService handlers.URLShortener
}

func New(cfg *config.HTTPServer, urlService handlers.URLShortener) *HTTPServer {
	httpServer := &HTTPServer{
		mux:        chi.NewMux(),
		validator:  validator.NewValidator(),
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
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)
	s.mux.Route("/url", func(r chi.Router) {
		r.Post("/", handlers.CreateShortURL(s.validator, s.urlService))
		r.Get("/{alias}", handlers.Redirect(s.urlService))
	})
}

func (s *HTTPServer) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
