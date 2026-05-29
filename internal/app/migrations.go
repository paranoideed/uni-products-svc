package app

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/paranoideed/uni-products-svc/migrations"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

func (a *App) MigrateUp(ctx context.Context) error {
	pool, err := a.config.PoolDB(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			a.log.Error("failed to close database connection", "error", err)
			return
		}
	}(db)

	if err = db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping database")
	}

	applied, err := migrate.ExecContext(ctx, db, "postgres", migrations.Migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations (up)")
	}
	a.log.Info("applied migrations", "applied", applied)

	return nil
}

func (a *App) MigrateDown(ctx context.Context) error {
	pool, err := a.config.PoolDB(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			a.log.Error("failed to close database connection", "error", err)
			return
		}
	}(db)

	if err = db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "failed to ping database")
	}

	applied, err := migrate.ExecContext(ctx, db, "postgres", migrations.Migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations (down)")
	}
	a.log.Info("rolled back migrations", "applied", applied)

	return nil
}
