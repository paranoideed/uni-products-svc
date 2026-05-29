package app

import (
	"log/slog"

	"github.com/paranoideed/uni-products-svc/internal/config"
)

type App struct {
	log    *slog.Logger
	config *config.Config
}

func New(log *slog.Logger, cfg *config.Config) *App {
	return &App{
		log:    log,
		config: cfg,
	}
}
