package repo_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/paranoideed/uni-products-svc/internal/repo"
	"go.yaml.in/yaml/v3"
)

type config struct {
	DB struct {
		URL string `yaml:"url"`
	} `yaml:"db"`
}

var testRepo *repo.Repo

func TestMain(m *testing.M) {
	cfg, err := loadConfig("test_config.yaml")
	if err != nil {
		panic("load test config: " + err.Error())
	}

	pool, err := pgxpool.New(context.Background(), cfg.DB.URL)
	if err != nil {
		panic("connect to db: " + err.Error())
	}
	defer pool.Close()

	if err = pool.Ping(context.Background()); err != nil {
		panic("ping db: " + err.Error())
	}

	testRepo = repo.NewRepo(pool)

	os.Exit(m.Run())
}

func loadConfig(path string) (*config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
