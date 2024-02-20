package item

import (
	"context"
	"errors"
	usrgate "github.com/JMURv/e-commerce/items/internal/gateway/users"
	"github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
)

var ErrNotFound = errors.New("not found")

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
	usrGateway usrgate.Gateway
}

func New(repo itemRepository, usrGateway usrgate.Gateway) *Controller {
	return &Controller{
		repo:       repo,
		usrGateway: usrGateway,
	}
}

func (ctrl *Controller) ListItem(ctx context.Context) ([]*model.Item, error) {
	res, err := ctrl.repo.ListItem(ctx)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) GetItemByID(ctx context.Context, id uint64) (*model.Item, error) {
	res, err := ctrl.repo.GetItemByID(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) ListUserItemsByID(ctx context.Context, userID uint64) ([]*model.Item, error) {
	res, err := ctrl.repo.ListUserItemsByID(ctx, userID)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) CreateItem(ctx context.Context, i *model.Item) (*model.Item, error) {
	res, err := ctrl.repo.CreateItem(ctx, i)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) UpdateItem(ctx context.Context, itemID uint64, newData *model.Item) (*model.Item, error) {
	res, err := ctrl.repo.UpdateItem(ctx, itemID, newData)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) DeleteItem(ctx context.Context, itemID uint64) error {
	if err := ctrl.repo.DeleteItem(ctx, itemID); err != nil {
		return err
	}
	return nil
}

func (ctrl *Controller) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	res, err := ctrl.repo.GetAllCategories(ctx)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) CreateCategory(ctx context.Context, c *model.Category) (*model.Category, error) {
	res, err := ctrl.repo.CreateCategory(ctx, c)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) UpdateCategory(ctx context.Context, categoryID uint64, newData *model.Category) (*model.Category, error) {
	res, err := ctrl.repo.UpdateCategory(ctx, categoryID, newData)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (ctrl *Controller) DeleteCategory(ctx context.Context, categoryID uint64) error {
	if err := ctrl.repo.DeleteCategory(ctx, categoryID); err != nil {
		return err
	}
	return nil
}

func (ctrl *Controller) ListTags(ctx context.Context) ([]*model.Tag, error) {
	res, err := ctrl.repo.ListTags(ctx)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (ctrl *Controller) CreateTag(ctx context.Context, t *model.Tag) (*model.Tag, error) {
	res, err := ctrl.repo.CreateTag(ctx, t)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (ctrl *Controller) DeleteTag(ctx context.Context, tagName string) error {
	if err := ctrl.repo.DeleteTag(ctx, tagName); err != nil {
		return err
	}
	return nil
}
