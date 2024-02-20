package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
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

func (r *Repository) GetByID(_ context.Context, id uint64) (*model.Review, error) {
	r.RLock()
	defer r.RUnlock()

	rev, ok := r.data[id]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return rev, nil
}

func (r *Repository) GetReviewsByUserID(_ context.Context, userID uint64) ([]*model.Review, error) {
	return nil, nil
}

func (r *Repository) AggregateUserRatingByID(_ context.Context, userID uint64) (float32, error) {
	return 0.0, nil
}

func (r *Repository) Create(_ context.Context, rev *model.Review) (*model.Review, error) {
	rev.ID = uint64(time.Now().Unix())
	r.Lock()
	defer r.Unlock()

	r.data[rev.ID] = rev
	return rev, nil
}

func (r *Repository) Update(_ context.Context, reviewID uint64, newData *model.Review) (*model.Review, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.data[reviewID]; ok {
		r.data[reviewID] = newData
		return newData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) Delete(_ context.Context, reviewID uint64) error {
	r.Lock()
	delete(r.data, reviewID)
	r.Unlock()
	return nil
}
