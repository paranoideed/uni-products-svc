package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/paranoideed/uni-products-svc/internal/app"
	"github.com/paranoideed/uni-products-svc/internal/config"
)

func Run(args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	application := app.New(cfg)
	log := application.Logger()

	var (
		service    = kingpin.New("service", "")
		runCmd     = service.Command("run", "run command flags: service")
		serviceCmd = runCmd.Command("service", "starting all service processes")

		migrateCmd     = service.Command("migrate", "migrate command")
		migrateUpCmd   = migrateCmd.Command("up", "migrate db up")
		migrateDownCmd = migrateCmd.Command("down", "migrate db down")
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	command, err := service.Parse(args[1:])
	if err != nil {
		return fmt.Errorf("error parsing command: %w", err)
	}

	switch command {
	case serviceCmd.FullCommand():
		err = application.Run(ctx)
	case migrateUpCmd.FullCommand():
		err = application.MigrateUp(ctx)
	case migrateDownCmd.FullCommand():
		err = application.MigrateDown(ctx)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	if err != nil {
		log.Error("error executing command", "command", command, "err", err)
		return err
	}

	log.Info("all processes finished successfully")
	return nil
}
