package review

import (
	"context"
	"errors"
	"fmt"
	notifygate "github.com/JMURv/e-commerce/reviews/internal/gateway/notifications"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"time"
)

var ErrNotFound = errors.New("not found")

type CacheRepository interface {
	GetReviewFromCache(ctx context.Context, key string) (*model.Review, error)
	SetReviewToCache(ctx context.Context, t time.Duration, r *model.Review) error
}

type reviewRepository interface {
	GetByID(ctx context.Context, reviewID uint64) (*model.Review, error)
	GetReviewsByUserID(ctx context.Context, userID uint64) (*[]*model.Review, error)
	AggregateUserRatingByID(ctx context.Context, userID uint64) (float32, error)
	Create(ctx context.Context, review *model.Review) (*model.Review, error)
	Update(ctx context.Context, reviewID uint64, newData *model.Review) (*model.Review, error)
	Delete(ctx context.Context, reviewID uint64) error
}

type Controller struct {
	repo          reviewRepository
	cache         CacheRepository
	notifyGateway *notifygate.Gateway
}

func New(repo reviewRepository, cache CacheRepository, gate *notifygate.Gateway) *Controller {
	return &Controller{repo, cache, gate}
}

func (c *Controller) GetReviewByID(ctx context.Context, reviewID uint64) (*model.Review, error) {
	cached, err := c.cache.GetReviewFromCache(ctx, fmt.Sprintf("review:%d", reviewID))
	if err == nil {
		return cached, nil
	}

	r, err := c.repo.GetByID(ctx, reviewID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}

	err = c.cache.SetReviewToCache(ctx, time.Hour, r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Controller) AggregateUserRatingByID(ctx context.Context, userID uint64) (float32, error) {
	r, err := c.repo.AggregateUserRatingByID(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return 0.0, ErrNotFound
	}
	return r, err
}

func (c *Controller) GetReviewsByUserID(ctx context.Context, userID uint64) (*[]*model.Review, error) {
	r, err := c.repo.GetReviewsByUserID(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}
	return r, err
}

func (c *Controller) CreateReview(ctx context.Context, review *model.Review) (*model.Review, error) {
	r, err := c.repo.Create(ctx, review)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}

	// Save review to cache
	err = c.cache.SetReviewToCache(ctx, time.Hour, r)
	if err != nil {
		return r, err
	}
	// Create notification for new review
	err = c.notifyGateway.CreateReviewNotification(ctx, r)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (c *Controller) UpdateReview(ctx context.Context, reviewID uint64, newData *model.Review) (*model.Review, error) {
	res, err := c.repo.Update(ctx, reviewID, newData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (c *Controller) DeleteReview(ctx context.Context, reviewID uint64) error {
	err := c.repo.Delete(ctx, reviewID)
	if err != nil {
		return err
	}
	return nil
}
