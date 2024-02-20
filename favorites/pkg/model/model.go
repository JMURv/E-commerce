package model

type Favorite struct {
	ID     uint64 `gorm:"primaryKey"`
	UserID uint64 `json:"user_id"`
	ItemID uint64 `json:"item_id"`
}
