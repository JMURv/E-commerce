package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/favorites/internal/repository"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
	"sync"
	"time"
)

type Repository struct {
	sync.RWMutex
	data map[uint64]*model.Favorite
}

func New() *Repository {
	return &Repository{data: map[uint64]*model.Favorite{}}
}

func (r *Repository) GetAllUserFavorites(_ context.Context, userID uint64) ([]*model.Favorite, error) {
	r.RLock()
	defer r.RUnlock()

	res := make([]*model.Favorite, 0, len(r.data))
	for _, v := range r.data {
		if v.UserID == userID {
			res = append(res, v)
		}
	}
	return res, nil
}

func (r *Repository) GetFavoriteByID(_ context.Context, favoriteID uint64) (*model.Favorite, error) {
	r.RLock()
	defer r.RUnlock()

	if f, ok := r.data[favoriteID]; ok {
		return f, nil
	} else {
		return nil, repo.ErrNotFound
	}
}

func (r *Repository) CreateFavorite(_ context.Context, favData *model.Favorite) (*model.Favorite, error) {
	favData.ID = uint64(time.Now().Unix())

	r.Lock()
	r.data[favData.ID] = favData
	r.Unlock()

	return favData, nil
}

func (r *Repository) DeleteFavorite(_ context.Context, favoriteID uint64) error {
	r.Lock()
	delete(r.data, favoriteID)
	r.Unlock()
	return nil
}
