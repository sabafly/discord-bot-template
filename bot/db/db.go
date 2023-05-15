package db

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sabafly/sabafly-lib/db"
)

func SetupDatabase(cfg db.DBConfig) (*DB, error) {
	db := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DB:      cfg.DB,
	})
	return &DB{
		db: db,
	}, nil
}

var _ db.DB = (*DB)(nil)

type DB struct {
	db *redis.Client
}

func (d *DB) Close() error {
	return d.db.Close()
}
