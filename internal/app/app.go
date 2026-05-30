package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/bridges/otelslog"

	"github.com/paranoideed/uni-products-svc/internal/config"
)

type App struct {
	config *config.Config
	log    *slog.Logger
}

func New(cfg *config.Config) *App {
	return &App{config: cfg}
}

func (a *App) Logger() *slog.Logger {
	if a.log == nil {
		a.initLogger()
	}
	return a.log
}

func (a *App) initLogger() {
	a.log = a.buildLogger()
}

func (a *App) PoolDB(ctx context.Context) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_SQL_URL"))
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	poolCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return pool, nil
}

func (a *App) buildLogger() *slog.Logger {
	lvl := slog.LevelInfo
	switch strings.ToLower(strings.TrimSpace(a.config.Log.Level)) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}

	var stdoutHandler slog.Handler
	switch strings.ToLower(strings.TrimSpace(a.config.Log.Format)) {
	case "json":
		stdoutHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	default:
		stdoutHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	}

	otelHandler := otelslog.NewHandler("uni-products-svc")

	return slog.New(&multiHandler{handlers: []slog.Handler{stdoutHandler, otelHandler}})
}

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}
