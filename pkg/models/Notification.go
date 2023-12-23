package models

import (
	"errors"
	"time"
)

type Notification struct {
	Type       string    `json:"type"`
	UserID     uint      `json:"userID"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	ReceiverID uint      `json:"receiverID"`
	Receiver   User      `json:"receiver" gorm:"foreignKey:ReceiverID"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (n *Notification) CreateNotification() (*Notification, error) {
	if n.Type == "" {
		return nil, errors.New("type is required")
	}

	if n.UserID == 0 {
		return nil, errors.New("userID is required")
	}

	if n.ReceiverID == 0 {
		return nil, errors.New("receiverID is required")
	}

	if n.Message == "" {
		return nil, errors.New("message is required")
	}

	n.CreatedAt = time.Now()
	if err := db.Preload("User").Preload("Receiver").Create(&n).Error; err != nil {
		return nil, err
	}

	return n, nil
}
