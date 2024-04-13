package controller

import (
	"context"
	"errors"
	repo "github.com/JMURv/e-commerce/chat/internal/repository"
	mdl "github.com/JMURv/e-commerce/chat/pkg/model"
	"time"
)

const roomsCacheKey = "room:%v"

var ErrNotFound = errors.New("not found")
var ErrUserIDRequired = errors.New("userID is required")
var ErrRoomIDRequired = errors.New("roomID is required")
var ErrItemIDRequired = errors.New("itemID is required")
var ErrTextRequired = errors.New("text is required")

type BrokerRepository interface {
	NewMessageNotification(msgID uint64, msg []byte) error
}

type CacheRepository interface {
	Get(ctx context.Context, key string) (*mdl.Room, error)
	Set(ctx context.Context, t time.Duration, key string, r *mdl.Room) error
	Delete(ctx context.Context, key string) error
}

type ChatRepository interface {
	GetMessageByID(ctx context.Context, msgID uint64) (*mdl.Message, error)
	CreateMessage(ctx context.Context, msgData *mdl.Message) (*mdl.Message, error)
	UpdateMessage(ctx context.Context, msgID uint64, msgData *mdl.Message) (*mdl.Message, error)
	DeleteMessage(ctx context.Context, msgID uint64) error

	UploadMedia(ctx context.Context, file []byte) (*mdl.Media, error)

	GetRoomByID(ctx context.Context, roomID uint64) (*mdl.Room, error)
	CreateRoom(ctx context.Context, room *mdl.Room) (*mdl.Room, error)
	GetUserRooms(ctx context.Context, userID uint64) ([]*mdl.Room, error)
	DeleteRoom(ctx context.Context, roomID uint64) error
}

type Controller struct {
	repo   ChatRepository
	cache  CacheRepository
	broker BrokerRepository
}

func New(repo ChatRepository, cache CacheRepository, broker BrokerRepository) *Controller {
	return &Controller{
		repo:   repo,
		cache:  cache,
		broker: broker,
	}
}

// Rooms
func (c *Controller) GetRoomByID(ctx context.Context, roomID uint64) (*mdl.Room, error) {
	r, err := c.repo.GetRoomByID(ctx, roomID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Controller) CreateRoom(ctx context.Context, room *mdl.Room) (*mdl.Room, error) {
	r, err := c.repo.CreateRoom(ctx, room)
	if err != nil && errors.Is(err, repo.ErrUserIDRequired) {
		return nil, ErrUserIDRequired
	} else if err != nil && errors.Is(err, repo.ErrItemIDRequired) {
		return nil, ErrItemIDRequired
	} else if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Controller) GetUserRooms(ctx context.Context, userID uint64) ([]*mdl.Room, error) {
	r, err := c.repo.GetUserRooms(ctx, userID)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (c *Controller) DeleteRoom(ctx context.Context, roomID uint64) error {
	err := c.repo.DeleteRoom(ctx, roomID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

//Messages

func (c *Controller) GetMessageByID(ctx context.Context, msgID uint64) (*mdl.Message, error) {
	r, err := c.repo.GetMessageByID(ctx, msgID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Controller) CreateMessage(ctx context.Context, msgData *mdl.Message) (*mdl.Message, error) {
	r, err := c.repo.CreateMessage(ctx, msgData)
	if err != nil && errors.Is(err, repo.ErrUserIDRequired) {
		return nil, ErrUserIDRequired
	} else if err != nil && errors.Is(err, repo.ErrRoomIDRequired) {
		return nil, ErrRoomIDRequired
	} else if err != nil && errors.Is(err, repo.ErrTextRequired) {
		return nil, ErrTextRequired
	} else if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Controller) UpdateMessage(ctx context.Context, msgID uint64, msgData *mdl.Message) (*mdl.Message, error) {
	r, err := c.repo.UpdateMessage(ctx, msgID, msgData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Controller) DeleteMessage(ctx context.Context, msgID uint64) error {
	err := c.repo.DeleteMessage(ctx, msgID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	return nil
}

// Media
func (c *Controller) UploadMedia(ctx context.Context, file []byte) (*mdl.Media, error) {
	media, err := c.repo.UploadMedia(ctx, file)
	if err != nil {
		return nil, err
	}
	return media, nil
}
