package item

import (
	"context"
	"errors"
	"fmt"
	usrgate "github.com/JMURv/e-commerce/items/internal/gateway/users"
	"github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"time"
)

const itemCacheKey = "item:%v"

var ErrNotFound = errors.New("not found")

type BrokerRepository interface{}

type CacheRepository interface {
	Get(ctx context.Context, key string) (*model.Item, error)
	Set(ctx context.Context, t time.Duration, key string, r *model.Item) error
	Delete(ctx context.Context, key string) error
}

type itemRepository interface {
	ListItem(_ context.Context) ([]*model.Item, error)
	GetItemByID(_ context.Context, itemID uint64) (*model.Item, error)
	ListUserItemsByID(_ context.Context, userID uint64) ([]*model.Item, error)
	CreateItem(_ context.Context, i *model.Item) (*model.Item, error)
	UpdateItem(_ context.Context, itemID uint64, newData *model.Item) (*model.Item, error)
	DeleteItem(_ context.Context, itemID uint64) error

	GetAllCategories(_ context.Context) ([]*model.Category, error)
	CreateCategory(_ context.Context, c *model.Category) (*model.Category, error)
	UpdateCategory(_ context.Context, categoryID uint64, newData *model.Category) (*model.Category, error)
	DeleteCategory(_ context.Context, categoryID uint64) error

	ListTags(_ context.Context) ([]*model.Tag, error)
	CreateTag(_ context.Context, t *model.Tag) (*model.Tag, error)
	DeleteTag(_ context.Context, tagName string) error
}

type Controller struct {
	repo       itemRepository
	cache      CacheRepository
	broker     BrokerRepository
	usrGateway usrgate.Gateway
}

func New(repo itemRepository, cache CacheRepository, broker BrokerRepository, usrGateway usrgate.Gateway) *Controller {
	return &Controller{
		repo:       repo,
		cache:      cache,
		broker:     broker,
		usrGateway: usrGateway,
	}
}

func (c *Controller) ListItem(ctx context.Context) ([]*model.Item, error) {
	res, err := c.repo.ListItem(ctx)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}

	return res, nil
}

func (c *Controller) GetItemByID(ctx context.Context, itemID uint64) (*model.Item, error) {
	cached, err := c.cache.Get(ctx, fmt.Sprintf(itemCacheKey, itemID))
	if err == nil {
		return cached, nil
	}

	res, err := c.repo.GetItemByID(ctx, itemID)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}

	return res, nil
}

func (c *Controller) ListUserItemsByID(ctx context.Context, userID uint64) ([]*model.Item, error) {
	res, err := c.repo.ListUserItemsByID(ctx, userID)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, nil
}

func (c *Controller) CreateItem(ctx context.Context, i *model.Item) (*model.Item, error) {
	res, err := c.repo.CreateItem(ctx, i)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(itemCacheKey, res.ID), res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Controller) UpdateItem(ctx context.Context, itemID uint64, newData *model.Item) (*model.Item, error) {
	res, err := c.repo.UpdateItem(ctx, itemID, newData)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(itemCacheKey, res.ID), res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Controller) DeleteItem(ctx context.Context, itemID uint64) error {
	if err := c.repo.DeleteItem(ctx, itemID); err != nil {
		return err
	}

	if err := c.cache.Delete(ctx, fmt.Sprintf(itemCacheKey, itemID)); err != nil {
		return err
	}

	return nil
}

// TODO: Прогрев кэша на категории и тэги?

func (c *Controller) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	res, err := c.repo.GetAllCategories(ctx)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, nil
}

func (c *Controller) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	res, err := c.repo.CreateCategory(ctx, category)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, nil
}

func (c *Controller) UpdateCategory(ctx context.Context, categoryID uint64, newData *model.Category) (*model.Category, error) {
	res, err := c.repo.UpdateCategory(ctx, categoryID, newData)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, nil
}

func (c *Controller) DeleteCategory(ctx context.Context, categoryID uint64) error {
	if err := c.repo.DeleteCategory(ctx, categoryID); err != nil {
		return nil
	}
	return nil
}

func (c *Controller) ListTags(ctx context.Context) ([]*model.Tag, error) {
	res, err := c.repo.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Controller) CreateTag(ctx context.Context, t *model.Tag) (*model.Tag, error) {
	res, err := c.repo.CreateTag(ctx, t)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Controller) DeleteTag(ctx context.Context, tagName string) error {
	if err := c.repo.DeleteTag(ctx, tagName); err != nil {
		return err
	}
	return nil
}
