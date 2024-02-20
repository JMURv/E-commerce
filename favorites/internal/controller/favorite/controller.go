package favorite

import (
	"context"
	"errors"
	repo "github.com/JMURv/e-commerce/favorites/internal/repository"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
)

var ErrNotFound = errors.New("not found")

type favoriteRepository interface {
	GetAllUserFavorites(ctx context.Context, userID uint64) ([]*model.Favorite, error)
	GetFavoriteByID(ctx context.Context, favoriteID uint64) (*model.Favorite, error)
	CreateFavorite(ctx context.Context, favData *model.Favorite) (*model.Favorite, error)
	DeleteFavorite(ctx context.Context, favoriteID uint64) error
}

type Controller struct {
	repo favoriteRepository
}

func New(repo favoriteRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) GetAllUserFavorites(ctx context.Context, userID uint64) ([]*model.Favorite, error) {
	r, err := c.repo.GetAllUserFavorites(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, err
}

func (c *Controller) GetFavoriteByID(ctx context.Context, favoriteID uint64) (*model.Favorite, error) {
	r, err := c.repo.GetFavoriteByID(ctx, favoriteID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, err
}

func (c *Controller) CreateFavorite(ctx context.Context, favData *model.Favorite) (*model.Favorite, error) {
	r, err := c.repo.CreateFavorite(ctx, favData)
	if err != nil {
		return nil, ErrNotFound
	}
	return r, err
}

func (c *Controller) DeleteFavorite(ctx context.Context, favoriteID uint64) error {
	err := c.repo.DeleteFavorite(ctx, favoriteID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}
