package user

import (
	"context"
	"errors"
	"fmt"
	itmgate "github.com/JMURv/e-commerce/users/internal/gateway/items"
	repo "github.com/JMURv/e-commerce/users/internal/repository"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"time"
)

const cacheKey = "users:%v"

var ErrNotFound = errors.New("user not found")

type BrokerRepository interface{}

type CacheRepository interface {
	Get(ctx context.Context, key string) (*model.User, error)
	Set(ctx context.Context, t time.Duration, key string, r *model.User) error
	Delete(ctx context.Context, key string) error
}

type userRepository interface {
	GetUsersList(ctx context.Context) ([]*model.User, error)
	GetByID(ctx context.Context, userID uint64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, userData *model.User) (*model.User, error)
	Update(ctx context.Context, userID uint64, newData *model.User) (*model.User, error)
	Delete(ctx context.Context, userID uint64) error
}

type Controller struct {
	repo       userRepository
	cache      CacheRepository
	broker     BrokerRepository
	itmGateway *itmgate.Gateway
}

func New(repo userRepository, cache CacheRepository, broker BrokerRepository, itmGateway *itmgate.Gateway) *Controller {
	return &Controller{
		repo:       repo,
		cache:      cache,
		broker:     broker,
		itmGateway: itmGateway,
	}
}

func (c *Controller) GetUsersList(ctx context.Context) ([]*model.User, error) {
	res, err := c.repo.GetUsersList(ctx)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Controller) GetUserByID(ctx context.Context, userID uint64) (*model.User, error) {
	cached, err := c.cache.Get(ctx, fmt.Sprintf(cacheKey, userID))
	if err == nil {
		return cached, nil
	}

	res, err := c.repo.GetByID(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(cacheKey, userID), res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (c *Controller) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	res, err := c.repo.GetByEmail(ctx, email)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) CreateUser(ctx context.Context, userData *model.User) (*model.User, error) {
	res, err := c.repo.Create(ctx, userData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(cacheKey, res.ID), res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (c *Controller) UpdateUser(ctx context.Context, userID uint64, newData *model.User) (*model.User, error) {
	res, err := c.repo.Update(ctx, userID, newData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, repo.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(cacheKey, res.ID), res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (c *Controller) DeleteUser(ctx context.Context, userID uint64) error {
	if err := c.repo.Delete(ctx, userID); err != nil {
		return err
	}

	if err := c.cache.Delete(ctx, fmt.Sprintf(cacheKey, userID)); err != nil {
		return err
	}
	return nil
}
