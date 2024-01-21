package review

import (
	"context"
	"errors"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
)

var ErrNotFound = errors.New("not found")

type reviewRepository interface {
	GetByID(ctx context.Context, reviewID uint64) (*model.Review, error)
	Create(ctx context.Context, review *model.Review) (*model.Review, error)
	Update(ctx context.Context, reviewID uint64, newData *model.Review) (*model.Review, error)
	Delete(ctx context.Context, reviewID uint64) error
}

type Controller struct {
	repo reviewRepository
}

func New(repo reviewRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) GetReviewByID(ctx context.Context, reviewID uint64) (*model.Review, error) {
	res, err := c.repo.GetByID(ctx, reviewID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}

func (c *Controller) CreateReview(ctx context.Context, review *model.Review) (*model.Review, error) {
	res, err := c.repo.Create(ctx, review)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
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
