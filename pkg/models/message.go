package models

import (
	"errors"
	"time"
)

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"userID"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	RoomID    uint      `json:"roomID"`
	Room      Room      `json:"room" gorm:"foreignKey:RoomID"`
	Text      string    `json:"text"`
	Seen      bool      `json:"seen"`
	Edited    bool      `json:"edited"`
	ReplyToID *uint     `json:"replyToID"`
	ReplyTo   *Message  `json:"replyTo" gorm:"foreignKey:ReplyToID"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m *Message) CreateMessage() (*Message, error) {
	if m.UserID == 0 {
		return nil, errors.New("userID is required")
	}

	if m.RoomID == 0 {
		return nil, errors.New("roomID is required")
	}

	if m.Text == "" {
		return nil, errors.New("text is required")
	}

	if err := db.Preload("User").Create(&m).Error; err != nil {
		return nil, err
	}

	return m, nil
}
