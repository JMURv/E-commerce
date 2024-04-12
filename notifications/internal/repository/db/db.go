package db

import (
	"context"
	"fmt"
	repo "github.com/JMURv/e-commerce/notifications/internal/repository"
	conf "github.com/JMURv/e-commerce/notifications/pkg/config"
	"github.com/JMURv/e-commerce/notifications/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository struct {
	conn *gorm.DB
}

func New(conf *conf.Config) *Repository {
	DSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Database,
	)

	conn, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = conn.AutoMigrate(
		&model.Notification{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: conn}
}

func (r *Repository) ListUserNotifications(_ context.Context, userID uint64) (*[]*model.Notification, error) {
	var n []*model.Notification
	if err := r.conn.Where("ReceiverID=?", userID).Find(&n).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) CreateNotification(_ context.Context, notify *model.Notification) (*model.Notification, error) {
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

	if err := r.conn.Create(notify).Error; err != nil {
		return nil, err
	}

	return notify, nil
}

func (r *Repository) DeleteNotification(_ context.Context, notificationID uint64) error {
	if err := r.conn.Delete(&model.Notification{}, notificationID).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteAllNotifications(_ context.Context, userID uint64) error {
	if err := r.conn.Where("ReceiverID=?", userID).Delete(&model.Notification{}).Error; err != nil {
		return err
	}
	return nil
}
