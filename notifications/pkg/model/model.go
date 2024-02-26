package model

import "time"

type Notification struct {
	Type       string    `json:"type"`
	UserID     uint64    `json:"user_id"`
	User       string    `json:"user"`
	ReceiverID uint64    `json:"receiver_id"`
	Receiver   string    `json:"receiver"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"createdAt"`
}
