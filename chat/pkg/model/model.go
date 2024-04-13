package model

import "time"

type Media struct {
	ID       uint64 `gorm:"primaryKey"`
	FilePath string `json:"file_path"`
}

type Room struct {
	ID        uint64     `gorm:"primaryKey"`
	SellerID  uint64     `json:"seller_id"`
	BuyerID   uint64     `json:"buyer_id"`
	ItemID    uint64     `json:"item_id"`
	Item      string     `json:"item"`
	Messages  []*Message `json:"messages"`
	CreatedAt time.Time  `json:"created_at"`
}

type Message struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `json:"user_id"`
	RoomID    uint64    `json:"room_id"`
	Text      string    `json:"text"`
	Seen      bool      `json:"seen"`
	Edited    bool      `json:"edited"`
	ReplyToID *uint64   `json:"reply_to_id"`
	ReplyTo   *Message  `json:"reply_to"`
	Media     []*Media  `json:"media"`
	CreatedAt time.Time `json:"created_at"`
}
