package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"

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
	poolCfg, err := pgxpool.ParseConfig(a.config.Postgres.DSN())
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

	switch strings.ToLower(strings.TrimSpace(a.config.Log.Format)) {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	}
}
