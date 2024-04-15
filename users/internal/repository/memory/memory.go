package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/users/internal/repository"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"github.com/opentracing/opentracing-go"
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

func (r *Repository) GetUsersList(ctx context.Context) ([]*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "users.GetUsersList.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()

	res := make([]*model.User, 0, len(r.data))
	for _, v := range r.data {
		res = append(res, v)
	}
	return res, nil
}

func (r *Repository) GetByID(ctx context.Context, userID uint64) (*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "users.GetByID.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()
	for _, v := range r.data {
		if v.ID == userID {
			return v, nil
		}
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "users.GetByEmail.repo")
	defer span.Finish()

	r.RLock()
	defer r.RUnlock()
	for _, v := range r.data {
		if v.Email == email {
			return v, nil
		}
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "users.CreateUser.repo")
	defer span.Finish()

	u.ID = uint64(time.Now().Unix())

	if u.Username == "" {
		return nil, repo.ErrUsernameIsRequired
	}

	if u.Email == "" {
		return nil, repo.ErrEmailIsRequired
	}

	r.Lock()
	r.data[u.ID] = u
	defer r.Unlock()
	return u, nil
}

func (r *Repository) Update(ctx context.Context, userID uint64, u *model.User) (*model.User, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "users.UpdateUser.repo")
	defer span.Finish()

	if u.Username == "" {
		return nil, repo.ErrUsernameIsRequired
	}

	if u.Email == "" {
		return nil, repo.ErrEmailIsRequired
	}

	r.Lock()
	defer r.Unlock()
	if _, ok := r.data[userID]; ok {
		r.data[userID] = u
		return u, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) Delete(ctx context.Context, userID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "users.DeleteUser.repo")
	defer span.Finish()

	r.Lock()
	defer r.Unlock()

	for k, v := range r.data {
		if v.ID == userID {
			delete(r.data, k)
			return nil
		}
	}
	return repo.ErrNotFound
}
