package cacheutil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ErrNotFoundInCache = errors.New("not found")

type Cache struct {
	cli *redis.Client
}

func New(cli *redis.Client) *Cache {
	return &Cache{cli}
}

func (c *Cache) GetReviewFromCache(ctx context.Context, key string) (*model.Review, error) {
	cachedData, err := c.cli.Get(ctx, key).Result()
	if err == nil {
		log.Printf("found obj in redis cache: %s", key)

		var r model.Review
		if err = json.Unmarshal([]byte(cachedData), &r); err != nil {
			return nil, err
		}
		return &r, nil
	}
	return nil, ErrNotFoundInCache
}

func (c *Cache) SetReviewToCache(ctx context.Context, t time.Duration, r *model.Review) error {
	newDataJSON, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if err = c.cli.Set(ctx, fmt.Sprintf("review:%d", r.ID), newDataJSON, t).Err(); err != nil {
		return err
	}
	return nil
}
