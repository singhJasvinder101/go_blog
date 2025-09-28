package storage

import (
	"github.com/singhJasvinder101/go_blog/internal/config"
	"github.com/singhJasvinder101/go_blog/storage/postgres"
)

type Storage interface {
	NewPostgres(cfg *config.Config) (*postgres.Postgres, error)
	InitSchema() error
	Close() error
}
