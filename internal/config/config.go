package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type RestConfig struct {
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
}

type Config struct {
	Log  LogConfig  `mapstructure:"log"`
	Rest RestConfig `mapstructure:"rest"`
}

func (cfg *Config) PoolDB(ctx context.Context) (*pgxpool.Pool, error) {
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

func (cfg *Config) Logger() *slog.Logger {
	lvl := slog.LevelInfo
	switch strings.ToLower(strings.TrimSpace(cfg.Log.Level)) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	}

	var handler slog.Handler
	switch strings.ToLower(strings.TrimSpace(cfg.Log.Format)) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: lvl,
		})
	}

	return slog.New(handler)
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		panic(fmt.Errorf("KV_VIPER_FILE env var is not set"))
	}
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}

	for _, key := range viper.AllKeys() {
		if val, ok := viper.Get(key).(string); ok {
			viper.Set(key, os.ExpandEnv(val))
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("error unmarshalling config: %s", err))
	}

	return &config, nil
}
