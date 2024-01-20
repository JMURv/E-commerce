package memory

import (
	"context"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"sync"
)

type Repository struct {
	sync.RWMutex
	data map[uint64]*model.User
}

func New() *Repository {
	return &Repository{data: map[uint64]*model.User{}}
}

func (r *Repository) GetByID(ctx context.Context, userID uint64) (*model.User, error) {
	panic("implement me")
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	panic("implement me")
}

func (r *Repository) Create(ctx context.Context, userData *model.User) (*model.User, error) {
	panic("implement me")
}

func (r *Repository) Update(ctx context.Context, userID uint64, newData *model.User) (*model.User, error) {
	panic("implement me")
}

func (r *Repository) Delete(ctx context.Context, userID uint64) error {
	panic("implement me")
}
