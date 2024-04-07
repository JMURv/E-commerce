package controller

import (
	"context"
	"time"
)

type entityType string

var (
	EntityTypeUser         entityType = "user"
	EntityTypeItem         entityType = "item"
	EntityTypeReview       entityType = "review"
	EntityTypeNotification entityType = "notification"
	EntityTypeFavorite     entityType = "favorite"
)

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, t time.Duration, key string, value []byte) error
	Delete(ctx context.Context, key string) error
}

type QueryRepository interface {
	FindAll(ctx context.Context, entityType entityType) ([]any, error)
	FindByID(ctx context.Context, entityType entityType, id string) (any, error)
	UserPage(ctx context.Context, userID uint64) ([]byte, error)
}

type Controller struct {
	repo  QueryRepository
	cache CacheRepository
}

func New(repo QueryRepository, cache CacheRepository) *Controller {
	return &Controller{
		repo:  repo,
		cache: cache,
	}
}
