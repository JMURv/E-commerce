package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"github.com/opentracing/opentracing-go"
	"sync"
	"time"
)

type Repository struct {
	sync.RWMutex
	data map[uint64]*model.Review
}

func New() *Repository {
	return &Repository{data: map[uint64]*model.Review{}}
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.GetByID.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	rev, ok := r.data[id]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return rev, nil
}

func (r *Repository) GetReviewsByUserID(ctx context.Context, userID uint64) (*[]*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.GetReviewsByUserID.repo")
	defer span.Finish()

	return nil, nil
}

func (r *Repository) AggregateUserRatingByID(ctx context.Context, userID uint64) (float32, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.AggregateUserRatingByID.repo")
	defer span.Finish()

	return 0.0, nil
}

func (r *Repository) Create(ctx context.Context, rev *model.Review) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.Create.repo")
	defer span.Finish()

	rev.ID = uint64(time.Now().Unix())
	r.Lock()
	defer r.Unlock()

	r.data[rev.ID] = rev
	return rev, nil
}

func (r *Repository) Update(ctx context.Context, reviewID uint64, newData *model.Review) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.Update.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()
	if _, ok := r.data[reviewID]; ok {
		r.data[reviewID] = newData
		return newData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) Delete(ctx context.Context, reviewID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.Delete.repo")
	defer span.Finish()

	r.Lock()
	delete(r.data, reviewID)
	r.Unlock()
	return nil
}
