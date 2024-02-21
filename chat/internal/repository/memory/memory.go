package memory

import (
	"context"
	repo "github.com/JMURv/e-commerce/chat/internal/repository"
	mdl "github.com/JMURv/e-commerce/chat/pkg/model"
	"sync"
	"time"
)

type Repository struct {
	sync.RWMutex
	messageData map[uint64]*mdl.Message
}

func New() *Repository {
	return &Repository{messageData: map[uint64]*mdl.Message{}}
}

func (r *Repository) GetMessageByID(_ context.Context, msgID uint64) (*mdl.Message, error) {
	r.RLock()
	defer r.RUnlock()

	if f, ok := r.messageData[msgID]; ok {
		return f, nil
	} else {
		return nil, repo.ErrNotFound
	}
}

func (r *Repository) CreateMessage(_ context.Context, msgData *mdl.Message) (*mdl.Message, error) {
	msgData.ID = uint64(time.Now().Unix())

	if msgData.UserID == 0 {
		return nil, repo.ErrUserIDRequired
	}

	if msgData.RoomID == 0 {
		return nil, repo.ErrRoomIDRequired
	}

	if msgData.Text == "" {
		return nil, repo.ErrTextRequired
	}

	r.Lock()
	r.messageData[msgData.ID] = msgData
	r.Unlock()

	return msgData, nil
}

func (r *Repository) UpdateMessage(_ context.Context, msgID uint64, msgData *mdl.Message) (*mdl.Message, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.messageData[msgID]; ok {
		r.messageData[msgID] = msgData
		return msgData, nil
	}
	return nil, repo.ErrNotFound
}

func (r *Repository) DeleteMessage(_ context.Context, msgID uint64) error {
	r.Lock()
	delete(r.messageData, msgID)
	r.Unlock()
	return nil
}
