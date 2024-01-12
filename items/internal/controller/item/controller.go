package item

import (
	"context"
	"errors"
	"github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/JMURv/protos/ecom/common"
)

var ErrNotFound = errors.New("not found")

type itemRepository interface {
	GetByID(ctx context.Context, id uint64) (*model.Item, error)
	Create(ctx context.Context, i *model.Item) (*model.Item, error)
	Update(ctx context.Context, itemID uint64, newData *model.Item) (*model.Item, error)
	Delete(ctx context.Context, itemID uint64) error
}

type Controller struct {
	repo itemRepository
}

func New(repo itemRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) GetItemByID(ctx context.Context, id uint64) (*model.Item, error) {
	res, err := c.repo.GetByID(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (c *Controller) CreateItem(ctx context.Context, i *common.Item) (*common.Item, error) {
	res, err := c.repo.Create(ctx, model.ItemFromProto(i))
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return model.ItemToProto(res), err
}

func (c *Controller) UpdateItem(ctx context.Context, itemID uint64, newData *common.Item) (*common.Item, error) {
	res, err := c.repo.Update(ctx, itemID, model.ItemFromProto(newData))
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return model.ItemToProto(res), err
}

func (c *Controller) DeleteItem(ctx context.Context, itemID uint64) error {
	err := c.repo.Delete(ctx, itemID)
	if err != nil {
		return err
	}
	return nil
}
