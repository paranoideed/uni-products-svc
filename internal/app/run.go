package app

import (
	"context"
	"sync"

	"github.com/paranoideed/uni-products-svc/internal/domain"
	"github.com/paranoideed/uni-products-svc/internal/repo"
	"github.com/paranoideed/uni-products-svc/internal/rest"
	"github.com/paranoideed/uni-products-svc/internal/rest/controller"
	"github.com/paranoideed/uni-products-svc/internal/rest/middlewares"
)

func (a *App) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	run := func(f func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}

	log := a.Logger()

	pool, err := a.PoolDB(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	repository := repo.NewRepo(pool)
	core := domain.NewService(repository)
	router := rest.NewServer(middlewares.New(), controller.New(core))

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
