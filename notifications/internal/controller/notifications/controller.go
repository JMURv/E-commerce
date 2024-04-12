package controller

import (
	"context"
	"errors"
	"fmt"
	repo "github.com/JMURv/e-commerce/notifications/internal/repository"
	"github.com/JMURv/e-commerce/notifications/pkg/model"
	"time"
)

const cacheKey = "notification:%v"

var ErrNotFound = errors.New("not found")

type BrokerRepository interface{}

type CacheRepository interface {
	Get(ctx context.Context, key string) (*model.Notification, error)
	Set(ctx context.Context, t time.Duration, key string, r *model.Notification) error
	Delete(ctx context.Context, key string) error
}

type notificationsRepository interface {
	ListUserNotifications(ctx context.Context, userID uint64) (*[]*model.Notification, error)
	CreateNotification(ctx context.Context, data *model.Notification) (*model.Notification, error)
	DeleteNotification(ctx context.Context, notificationID uint64) error
	DeleteAllNotifications(ctx context.Context, userID uint64) error
}

type Controller struct {
	repo   notificationsRepository
	cache  CacheRepository
	broker BrokerRepository
}

func New(repo notificationsRepository, cache CacheRepository, broker BrokerRepository) *Controller {
	return &Controller{
		repo:   repo,
		cache:  cache,
		broker: broker,
	}
}

func (c *Controller) ListUserNotifications(ctx context.Context, userID uint64) (*[]*model.Notification, error) {
	res, err := c.repo.ListUserNotifications(ctx, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Controller) CreateNotification(ctx context.Context, data *model.Notification) (*model.Notification, error) {
	res, err := c.repo.CreateNotification(ctx, data)
	// TODO: Check for errs
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(cacheKey, res.ID), res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *Controller) DeleteNotification(ctx context.Context, notificationID uint64) error {
	err := c.repo.DeleteNotification(ctx, notificationID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	if err = c.cache.Delete(ctx, fmt.Sprintf(cacheKey, notificationID)); err != nil {
		return err
	}

	return nil
}

func (c *Controller) DeleteAllNotifications(ctx context.Context, userID uint64) error {
	err := c.repo.DeleteAllNotifications(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}
