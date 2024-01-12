package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Review struct {
	ID             uint64                `gorm:"primaryKey"`
	UserID         uint64                `json:"user_id"`
	ItemID         uint64                `json:"item_id"`
	ReviewedUserID uint64                `json:"reviewed_user_id"`
	Advantages     string                `json:"advantages"`
	Disadvantages  string                `json:"disadvantages"`
	ReviewText     string                `json:"review_text"`
	Rating         uint32                `json:"rating"`
	CreatedAt      timestamppb.Timestamp `json:"created_at"`
	UpdatedAt      timestamppb.Timestamp `json:"updated_at"`
}
