package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/notifications/internal/repository"
	"github.com/JMURv/e-commerce/notifications/pkg/model"
	"sync"
	"time"
)

type Repository struct {
	sync.RWMutex
	data map[uint64]*model.Notification
}

func New() *Repository {
	return &Repository{data: map[uint64]*model.Notification{}}
}

func (r *Repository) ListUserNotifications(_ context.Context, userID uint64) ([]*model.Notification, error) {
	r.RLock()
	defer r.RUnlock()
	n := make([]*model.Notification, 0, len(r.data))
	for _, v := range r.data {
		if v.ReceiverID == userID {
			n = append(n, v)
		}
	}
	return n, nil
}

func (r *Repository) CreateNotification(_ context.Context, notify *model.Notification) (*model.Notification, error) {
	notify.ID = uint64(time.Now().Unix())

	if notify.Type == "" {
		return nil, repo.ErrTypeIsRequired
	}

	if notify.UserID == 0 {
		return nil, repo.ErrUserIDIsRequired
	}

	if notify.ReceiverID == 0 {
		return nil, repo.ErrIRecieverIDIsRequired
	}

	if notify.Message == "" {
		return nil, repo.ErrMessageIsRequired
	}

	notify.CreatedAt = time.Now()

	r.Lock()
	defer r.Unlock()
	r.data[notify.ID] = notify

	return notify, nil
}

func (r *Repository) DeleteNotification(_ context.Context, notificationID uint64) error {
	r.Lock()
	defer r.Unlock()
	delete(r.data, notificationID)
	return nil
}

func (r *Repository) DeleteAllNotifications(_ context.Context, userID uint64) error {
	r.Lock()
	defer r.Unlock()
	for i, n := range r.data {
		if userID == n.ReceiverID {
			delete(r.data, i)
			return nil
		}
	}
	return repo.ErrNotFound
}
