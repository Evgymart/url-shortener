package dragonfly

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"urlshort/internal/config"
	"urlshort/internal/storage"

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

func (s *Storage) Ping(ctx context.Context) error {
	return s.rdb.Ping(ctx).Err()
}

func (s *Storage) SaveURL(ctx context.Context, urlAlias string, longUrl string) error {
	exists, err := s.rdb.Exists(ctx, urlAlias).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		return storage.ErrAliasAlreadyTaken
	}
	return s.rdb.Set(ctx, urlAlias, longUrl, 0).Err()
}

func (s *Storage) GetUrl(ctx context.Context, urlAlias string) (string, error) {
	result, err := s.rdb.Get(ctx, urlAlias).Result()
	if errors.Is(err, redis.Nil) {
		return "", storage.ErrUrlNotFound
	}
	return result, err
}

func (s *Storage) DeleteURL(ctx context.Context, urlAlias string) error {
	return s.rdb.Del(ctx, urlAlias).Err()
}
