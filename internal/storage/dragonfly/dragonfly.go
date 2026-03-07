package dragonfly

import (
	"context"
	"fmt"
	"strconv"
	"urlshort/internal/config"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	rdb *redis.Client
}

func New(cfg config.Dragonfly) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port)),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Storage{client}
}

func (s *Storage) Ping() error {
	return s.rdb.Ping(context.Background()).Err()
}
