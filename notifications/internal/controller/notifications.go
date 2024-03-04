package controller

import (
	"context"
	"errors"
	repo "github.com/JMURv/e-commerce/notifications/internal/repository"
	"github.com/JMURv/e-commerce/notifications/pkg/model"
)

var ErrNotFound = errors.New("user not found")

type notificationsRepository interface {
	ListUserNotifications(ctx context.Context, userID uint64) ([]*model.Notification, error)
	CreateNotification(ctx context.Context, data *model.Notification) (*model.Notification, error)
	DeleteNotification(ctx context.Context, notificationID uint64) error
	DeleteAllNotifications(ctx context.Context, userID uint64) error
}

type Controller struct {
	repo notificationsRepository
}

func New(repo notificationsRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) ListUserNotifications(ctx context.Context, userID uint64) ([]*model.Notification, error) {
	res, err := c.repo.ListUserNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *Controller) CreateNotification(ctx context.Context, data *model.Notification) (*model.Notification, error) {
	res, err := c.repo.CreateNotification(ctx, data)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return res, err
}

func (c *Controller) DeleteNotification(ctx context.Context, notificationID uint64) error {
	err := c.repo.DeleteNotification(ctx, notificationID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	return err
}

func (c *Controller) DeleteAllNotifications(ctx context.Context, userID uint64) error {
	err := c.repo.DeleteAllNotifications(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	return err
}
