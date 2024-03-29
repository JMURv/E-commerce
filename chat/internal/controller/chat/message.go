package controller

import (
	"context"
	"errors"
	repo "github.com/JMURv/e-commerce/chat/internal/repository"
	mdl "github.com/JMURv/e-commerce/chat/pkg/model"
)

var ErrNotFound = errors.New("not found")
var ErrUserIDRequired = errors.New("userID is required")
var ErrRoomIDRequired = errors.New("roomID is required")
var ErrItemIDRequired = errors.New("itemID is required")
var ErrTextRequired = errors.New("text is required")

type ChatRepository interface {
	GetMessageByID(ctx context.Context, msgID uint64) (*mdl.Message, error)
	CreateMessage(ctx context.Context, msgData *mdl.Message) (*mdl.Message, error)
	UpdateMessage(ctx context.Context, msgID uint64, msgData *mdl.Message) (*mdl.Message, error)
	DeleteMessage(ctx context.Context, msgID uint64) error

	CreateRoom(ctx context.Context, room *mdl.Room) (*mdl.Room, error)
	GetUserRooms(ctx context.Context, userID uint64) ([]*mdl.Room, error)
	DeleteRoom(ctx context.Context, roomID uint64) error
}

type Controller struct {
	repo ChatRepository
}

func New(repo ChatRepository) *Controller {
	return &Controller{repo}
}

// Rooms
func (c *Controller) CreateRoom(ctx context.Context, room *mdl.Room) (*mdl.Room, error) {
	r, err := c.repo.CreateRoom(ctx, room)
	if err != nil && errors.Is(err, repo.ErrUserIDRequired) {
		return nil, ErrUserIDRequired
	} else if err != nil && errors.Is(err, repo.ErrItemIDRequired) {
		return nil, ErrItemIDRequired
	} else if err != nil {
		return nil, err
	}
	return r, err
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
	return err
}

//Messages

func (c *Controller) GetMessageByID(ctx context.Context, msgID uint64) (*mdl.Message, error) {
	r, err := c.repo.GetMessageByID(ctx, msgID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, err
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
	return r, err
}

func (c *Controller) UpdateMessage(ctx context.Context, msgID uint64, msgData *mdl.Message) (*mdl.Message, error) {
	r, err := c.repo.UpdateMessage(ctx, msgID, msgData)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return r, err
}

func (c *Controller) DeleteMessage(ctx context.Context, msgID uint64) error {
	err := c.repo.DeleteMessage(ctx, msgID)
	if err != nil && errors.Is(err, repo.ErrNotFound) {
		return ErrNotFound
	} else if err != nil {
		return err
	}
	return err
}
