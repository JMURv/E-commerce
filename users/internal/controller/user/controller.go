package user

import (
	"context"
	"errors"
	itmgate "github.com/JMURv/e-commerce/users/internal/gateway/items"
	repo "github.com/JMURv/e-commerce/users/internal/repository"
	"github.com/JMURv/e-commerce/users/pkg/model"
)

var ErrNotFound = errors.New("user not found")

type userRepository interface {
	GetUsersList(ctx context.Context) (*[]model.User, error)
	GetByID(ctx context.Context, userID uint64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, userData *model.User) (*model.User, error)
	Update(ctx context.Context, userID uint64, newData *model.User) (*model.User, error)
	Delete(ctx context.Context, userID uint64) error
}

type Controller struct {
	repo       userRepository
	itmGateway *itmgate.Gateway
}

func New(repo userRepository, itmGateway *itmgate.Gateway) *Controller {
	return &Controller{
		repo:       repo,
		itmGateway: itmGateway,
	}
}

func (c *Controller) GetUsersList(ctx context.Context) (*[]model.User, error) {
	res, err := c.repo.GetUsersList(ctx)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}

	return res, err
}

func (c *Controller) GetUserByID(ctx context.Context, userID uint64) (*model.User, error) {
	res, err := c.repo.GetByID(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}

	return res, err
}

func (c *Controller) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	res, err := c.repo.GetByEmail(ctx, email)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, repo.ErrNotFound
	}

	return res, err
}

func (c *Controller) CreateUser(ctx context.Context, userData *model.User) (*model.User, error) {
	res, err := c.repo.Create(ctx, userData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, repo.ErrNotFound
	}

	return res, err
}

func (c *Controller) UpdateUser(ctx context.Context, userID uint64, newData *model.User) (*model.User, error) {
	res, err := c.repo.Update(ctx, userID, newData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, repo.ErrNotFound
	}

	return res, err
}

func (c *Controller) DeleteUser(ctx context.Context, userID uint64) error {
	err := c.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
