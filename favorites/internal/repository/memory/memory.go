package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/favorites/internal/repository"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
	"github.com/opentracing/opentracing-go"
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

func (r *Repository) GetAllUserFavorites(ctx context.Context, userID uint64) ([]*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.GetAllUserFavorites.repo")
	defer span.Finish()

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

func (r *Repository) GetFavoriteByID(ctx context.Context, favoriteID uint64) (*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.GetFavoriteByID.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	if f, ok := r.data[favoriteID]; ok {
		return f, nil
	} else {
		return nil, repo.ErrNotFound
	}
}

func (r *Repository) CreateFavorite(ctx context.Context, favData *model.Favorite) (*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.CreateFavorite.repo")
	defer span.Finish()

	favData.ID = uint64(time.Now().Unix())

	r.Lock()
	r.data[favData.ID] = favData
	r.Unlock()

	return favData, nil
}

func (r *Repository) DeleteFavorite(ctx context.Context, favoriteID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.DeleteFavorite.repo")
	defer span.Finish()

	r.Lock()
	delete(r.data, favoriteID)
	r.Unlock()
	return nil
}
