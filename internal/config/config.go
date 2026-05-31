package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type RestConfig struct {
	Port     int `mapstructure:"port"`
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (c PostgresConfig) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DB, c.SSLMode)
}

type Config struct {
	Log      LogConfig      `mapstructure:"log"`
	Rest     RestConfig     `mapstructure:"rest"`
	Postgres PostgresConfig `mapstructure:"postgres"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		return nil, fmt.Errorf("KV_VIPER_FILE env var is not set")
	}
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	for _, key := range viper.AllKeys() {
		if val, ok := viper.Get(key).(string); ok {
			viper.Set(key, os.ExpandEnv(val))
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &config, nil
}
