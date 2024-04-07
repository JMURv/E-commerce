package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ErrNotFoundInCache = errors.New("not found")

type Cache struct {
	cli *redis.Client
}

func New(addr, pass string) *Cache {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})
	pong, err := redisCli.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis:", pong)
	return &Cache{cli: redisCli}
}

func (c *Cache) Close() {
	err := c.cli.Close()
	if err != nil {
		log.Println("Failed to close connection to Redis: ", err)
	}
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.cli.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrNotFoundInCache
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func (c *Cache) Set(ctx context.Context, t time.Duration, key string, value []byte) error {
	if err := c.cli.Set(ctx, key, value, t).Err(); err != nil {
		return err
	}
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := c.cli.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
