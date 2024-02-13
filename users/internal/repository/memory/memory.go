package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/users/internal/repository"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"sync"
	"time"
)

type Repository struct {
	sync.RWMutex
	data map[uint64]*model.User
}

func New() *Repository {
	return &Repository{data: map[uint64]*model.User{}}
}

func (r *Repository) GetUsersList(_ context.Context) (*[]model.User, error) {
	r.RLock()
	defer r.RUnlock()

	res := make([]model.User, 0, len(r.data))
	for _, v := range r.data {
		res = append(res, *v)
	}
	return &res, nil
}

func (r *Repository) GetByID(_ context.Context, userID uint64) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()
	for _, v := range r.data {
		if v.ID == userID {
			return v, nil
		}
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) GetByEmail(_ context.Context, email string) (*model.User, error) {
	r.RLock()
	defer r.RUnlock()
	for _, v := range r.data {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) Create(_ context.Context, u *model.User) (*model.User, error) {
	u.ID = uint64(time.Now().Unix())

	if u.Username == "" {
		return nil, repo.ErrUsernameIsRequired
	}

	if u.Email == "" {
		return nil, repo.ErrEmailIsRequired
	}

	r.Lock()
	r.data[u.ID] = u
	r.Unlock()
	return u, nil
}

func (r *Repository) Update(_ context.Context, userID uint64, newData *model.User) (*model.User, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.data[userID]; ok {
		r.data[userID] = newData
		return newData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) Delete(_ context.Context, userID uint64) error {
	r.Lock()
	delete(r.data, userID)
	r.Unlock()
	return nil
}
