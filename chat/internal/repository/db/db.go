package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/chat/internal/repository"
	mdl "github.com/JMURv/e-commerce/chat/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DSN string

type Repository struct {
	conn *gorm.DB
}

func New() *Repository {
	var err error
	var db *gorm.DB

	db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&mdl.Message{},
		&mdl.Room{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: db}
}

func (r *Repository) GetMessageByID(ctx context.Context, msgID uint64) (*mdl.Message, error) {
	var msg *mdl.Message
	if err := r.conn.Where("ID=?", msgID).First(&msg).Error; err != nil {
		return nil, repo.ErrNotFound
	}
	return msg, nil
}

func (r *Repository) CreateMessage(ctx context.Context, msgData *mdl.Message) (*mdl.Message, error) {
	var m *mdl.Message

	if msgData.UserID == 0 {
		return nil, repo.ErrUserIDRequired
	}

	if msgData.RoomID == 0 {
		return nil, repo.ErrRoomIDRequired
	}

	if msgData.Text == "" {
		return nil, repo.ErrTextRequired
	}

	if err := r.conn.Create(&m).Error; err != nil {
		return nil, err
	}

	return m, nil
}

func (r *Repository) UpdateMessage(ctx context.Context, msgID uint64, msgData *mdl.Message) (*mdl.Message, error) {
	msg, err := r.GetMessageByID(ctx, msgID)
	if err != nil {
		return nil, repo.ErrNotFound
	}

	if msgData.Text != "" {
		msg.Text = msgData.Text
	}

	msg.Edited = true
	if err = r.conn.Save(&msg).Error; err != nil {
		return nil, err
	}
	return msg, nil
}

func (r *Repository) DeleteMessage(ctx context.Context, msgID uint64) error {
	if err := r.conn.Delete(&mdl.Message{}, msgID).Error; err != nil {
		return err
	}
	return nil
}
