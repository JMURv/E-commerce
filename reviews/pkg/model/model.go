package model

import (
	"time"
)

type Review struct {
	ID             uint64    `gorm:"primaryKey"`
	UserID         uint64    `json:"user_id"`
	ItemID         uint64    `json:"item_id"`
	ReviewedUserID uint64    `json:"reviewed_user_id"`
	Advantages     string    `json:"advantages"`
	Disadvantages  string    `json:"disadvantages"`
	ReviewText     string    `json:"review_text"`
	Rating         uint32    `json:"rating"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
