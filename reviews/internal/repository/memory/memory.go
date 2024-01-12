package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"sync"
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

	i, ok := r.data[id]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return i, nil
}

func (r *Repository) Create(_ context.Context, i *model.Review) (*model.Review, error) {
	r.Lock()
	defer r.Unlock()

	r.data[i.ID] = i
	return i, nil
}

func (r *Repository) Update(_ context.Context, itemID uint64, newData *model.Review) (*model.Review, error) {
	r.Lock()
	defer r.Unlock()

	i, ok := r.data[itemID]
	if !ok {
		return nil, repo.ErrNotFound
	}
	return i, nil
}

func (r *Repository) Delete(_ context.Context, itemID uint64) error {

	return nil
}
