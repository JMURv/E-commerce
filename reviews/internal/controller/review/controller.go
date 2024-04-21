package review

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/notification"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"github.com/opentracing/opentracing-go"
	"log"
	"time"
)

const cacheKey = "review:%v"

var ErrNotFound = errors.New("not found")

type BrokerRepository interface {
	NewReviewNotification(reviewID uint64, notification []byte) error
}

type CacheRepository interface {
	Get(ctx context.Context, key string) (*model.Review, error)
	Set(ctx context.Context, t time.Duration, key string, r *model.Review) error
	Delete(ctx context.Context, key string) error
}

type ReviewRepository interface {
	GetByID(ctx context.Context, reviewID uint64) (*model.Review, error)
	GetReviewsByUserID(ctx context.Context, userID uint64) (*[]*model.Review, error)
	AggregateUserRatingByID(ctx context.Context, userID uint64) (float32, error)
	Create(ctx context.Context, review *model.Review) (*model.Review, error)
	Update(ctx context.Context, reviewID uint64, newData *model.Review) (*model.Review, error)
	Delete(ctx context.Context, reviewID uint64) error
}

type Controller struct {
	repo   ReviewRepository
	cache  CacheRepository
	broker BrokerRepository
}

func New(repo ReviewRepository, cache CacheRepository, broker BrokerRepository) *Controller {
	return &Controller{
		repo:   repo,
		cache:  cache,
		broker: broker,
	}
}

func (c *Controller) GetReviewByID(ctx context.Context, reviewID uint64) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.GetReviewByID.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	cached, err := c.cache.Get(ctx, fmt.Sprintf(cacheKey, reviewID))
	if err == nil {
		return cached, nil
	}

	r, err := c.repo.GetByID(ctx, reviewID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(cacheKey, reviewID), r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (c *Controller) AggregateUserRatingByID(ctx context.Context, userID uint64) (float32, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.AggregateUserRatingByID.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	r, err := c.repo.AggregateUserRatingByID(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return 0.0, ErrNotFound
	}
	return r, nil
}

func (c *Controller) GetReviewsByUserID(ctx context.Context, userID uint64) (*[]*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.GetReviewsByUserID.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	r, err := c.repo.GetReviewsByUserID(ctx, userID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Controller) CreateReview(ctx context.Context, review *model.Review) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.CreateReview.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.Create(ctx, review)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	// Save review to cache
	err = c.cache.Set(ctx, time.Hour, fmt.Sprintf(cacheKey, res.ID), res)
	if err != nil {
		return res, err
	}

	// Create notification for new review
	bytes, err := json.Marshal(pb.Notification{
		Type:       "new_review",
		UserId:     res.UserID,
		ReceiverId: res.ReviewedUserID,
		Message:    "New review received!",
	})
	if err != nil {
		log.Printf("Error notification encoding: %v", err)
		return res, err
	}

	if err = c.broker.NewReviewNotification(res.ID, bytes); err != nil {
		log.Println(err)
	}
	return res, nil
}

func (c *Controller) UpdateReview(ctx context.Context, reviewID uint64, newData *model.Review) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.UpdateReview.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	res, err := c.repo.Update(ctx, reviewID, newData)
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

func (c *Controller) DeleteReview(ctx context.Context, reviewID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.DeleteReview.ctrl")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer span.Finish()

	if err := c.repo.Delete(ctx, reviewID); err != nil {
		return err
	}

	if err := c.cache.Delete(ctx, fmt.Sprintf(cacheKey, reviewID)); err != nil {
		return err
	}
	return nil
}
