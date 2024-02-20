package models

import (
	"errors"
	"time"
)

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"userID"`
	User      string    `json:"user" gorm:"foreignKey:UserID"`
	RoomID    uint      `json:"roomID"`
	Room      Room      `json:"room" gorm:"foreignKey:RoomID"`
	Text      string    `json:"text"`
	Seen      bool      `json:"seen"`
	Edited    bool      `json:"edited"`
	ReplyToID *uint     `json:"replyToID"`
	ReplyTo   *Message  `json:"replyTo" gorm:"foreignKey:ReplyToID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func GetMessageByID(messageID uint) (*Message, error) {
	var message Message
	if err := db.Preload("User").
		Preload("Room").
		Preload("ReplyTo").
		Where("ID=?", messageID).
		First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
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

func UpdateMessage(messageID uint, newData *Message) (*Message, error) {
	message, err := GetMessageByID(messageID)
	if err != nil {
		return nil, err
	}

	if newData.Text != "" {
		message.Text = newData.Text
	}

	message.Edited = true
	message.UpdatedAt = time.Now()

	if err := db.Save(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func DeleteMessage(messageID uint) (Message, error) {
	var message Message
	if err := db.Delete(&message, messageID).Error; err != nil {
		return message, err
	}
	return message, nil
}
