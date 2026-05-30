package app

import (
	"context"
	"sync"

	"github.com/paranoideed/uni-products-svc/internal/domain"
	"github.com/paranoideed/uni-products-svc/internal/metrics"
	"github.com/paranoideed/uni-products-svc/internal/repo"
	"github.com/paranoideed/uni-products-svc/internal/rest"
	"github.com/paranoideed/uni-products-svc/internal/rest/controller"
	"github.com/paranoideed/uni-products-svc/internal/rest/middlewares"
	"github.com/paranoideed/uni-products-svc/internal/telemetry"
)

func (a *App) Run(ctx context.Context) error {
	shutdown, err := telemetry.Setup(ctx, "uni-products-svc")
	if err != nil {
		return err
	}

	a.initLogger()
	log := a.Logger()

	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Error("failed to shutdown telemetry gracefully", "error", err)
		}
	}()

	pool, err := a.PoolDB(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	repository := repo.NewRepo(pool)
	core := domain.NewService(repository)

	m, err := metrics.New()
	if err != nil {
		return err
	}

	router := rest.NewServer(middlewares.New(), controller.New(core, m))

	var wg sync.WaitGroup
	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	run(func() {
		router.Run(ctx, log, rest.Config{
			Port:              8000,
			ReadTimeout:       a.config.Rest.Timeouts.Read,
			ReadHeaderTimeout: a.config.Rest.Timeouts.ReadHeader,
			WriteTimeout:      a.config.Rest.Timeouts.Write,
			IdleTimeout:       a.config.Rest.Timeouts.Idle,
		})
	})

	log.Info("starting application")

	wg.Wait()

	return nil
}
