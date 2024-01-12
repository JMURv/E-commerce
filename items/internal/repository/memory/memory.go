package memory

import (
	"context"
	"github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	data map[uint64]*model.Item
}

func New() *Repository {
	return &Repository{data: map[uint64]*model.Item{}}
}

func (r *Repository) GetByID(_ context.Context, id uint64) (*model.Item, error) {
	r.RLock()
	defer r.RUnlock()

	i, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return i, nil
}

func (r *Repository) Create(_ context.Context, i *model.Item) (*model.Item, error) {
	r.Lock()
	defer r.Unlock()

	r.data[i.ID] = i
	return i, nil
}

func (r *Repository) Update(_ context.Context, itemID uint64, newData *model.Item) (*model.Item, error) {
	r.Lock()
	defer r.Unlock()

	i, ok := r.data[itemID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return i, nil
}

func (r *Repository) Delete(_ context.Context, itemID uint64) error {

	return nil
}
