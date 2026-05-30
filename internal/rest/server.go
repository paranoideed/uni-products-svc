package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Controller interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	Logger(log *slog.Logger) func(next http.Handler) http.Handler
	CorsDocs() func(next http.Handler) http.Handler
}

type Server struct {
	middlewares Middlewares
	controller  Controller
}

func NewServer(
	middlewares Middlewares,
	controller Controller,
) *Server {
	return &Server{
		middlewares: middlewares,
		controller:  controller,
	}
}

type Config struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (s *Server) Run(ctx context.Context, log *slog.Logger, cfg Config) {
	r := chi.NewRouter()

	r.Use(
		s.middlewares.Logger(log),
		s.middlewares.CorsDocs(),
	)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", s.controller.GetProducts)
		r.Post("/", s.controller.CreateProduct)
		r.Delete("/{product_id}", s.controller.DeleteProduct)
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           otelhttp.NewHandler(r, "uni-products-svc"),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	log.Info("starting http service", "port", cfg.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			log.Error("http server error", "error", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Error("failed to shutdown http server gracefully", "error", err)
	} else {
		log.Info("http server shutdown gracefully")
	}
}
