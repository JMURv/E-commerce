package favorite

import (
	"context"
	"errors"
	"fmt"
	repo "github.com/JMURv/e-commerce/favorites/internal/repository"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
	"github.com/opentracing/opentracing-go"
	"time"
)

const cacheKey = "favorite:%v"

var ErrNotFound = errors.New("not found")

type BrokerRepository interface{}

type CacheRepository interface {
	Get(ctx context.Context, key string) (*model.Favorite, error)
	Set(ctx context.Context, t time.Duration, key string, r *model.Favorite) error
	Delete(ctx context.Context, key string) error
}

type favoriteRepository interface {
	GetAllUserFavorites(ctx context.Context, userID uint64) ([]*model.Favorite, error)
	GetFavoriteByID(ctx context.Context, favoriteID uint64) (*model.Favorite, error)
	CreateFavorite(ctx context.Context, favData *model.Favorite) (*model.Favorite, error)
	DeleteFavorite(ctx context.Context, favoriteID uint64) error
}

type Controller struct {
	repo   favoriteRepository
	cache  CacheRepository
	broker BrokerRepository
}

func New(repo favoriteRepository, cache CacheRepository, broker BrokerRepository) *Controller {
	return &Controller{
		repo:   repo,
		cache:  cache,
		broker: broker,
	}
}

func (c *Controller) GetAllUserFavorites(ctx context.Context, userID uint64) ([]*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.GetAllUserFavorites.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	r, err := c.repo.GetAllUserFavorites(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, err
}

func (c *Controller) GetFavoriteByID(ctx context.Context, favoriteID uint64) (*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.GetFavoriteByID.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached, err := c.cache.Get(ctx, fmt.Sprintf(cacheKey, favoriteID))
	if err == nil {
		return cached, nil
	}

	res, err := c.repo.GetFavoriteByID(ctx, favoriteID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, err
}

func (c *Controller) CreateFavorite(ctx context.Context, favData *model.Favorite) (*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.CreateFavorite.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.CreateFavorite(ctx, favData)
	if err != nil {
		return nil, ErrNotFound
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(cacheKey, res.ID), res)
	if err != nil {
		return res, err
	}

	return res, err
}

func (c *Controller) DeleteFavorite(ctx context.Context, favoriteID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.DeleteFavorite.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	err := c.repo.DeleteFavorite(ctx, favoriteID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(cacheKey, favoriteID)); err != nil {
		return err
	}

	return nil
}
