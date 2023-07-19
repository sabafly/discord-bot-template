package db

import (
	"fmt"

	"github.com/go-redis/redis"
)

type DBConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	DB   int    `json:"db"`
}

func SetupDatabase(cfg DBConfig) (*DB, error) {
	db := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		DB:      cfg.DB,
	})
	return &DB{
		db: db,
	}, nil
}

type DB struct {
	db *redis.Client
}

func (d *DB) Close() error {
	return d.db.Close()
}
