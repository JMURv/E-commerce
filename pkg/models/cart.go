package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint       `json:"userID"`
	Items  []CartItem `json:"items" gorm:"many2many:cart_items;"`
}

type CartItem struct {
	CartID   uint `gorm:"primaryKey"`
	ItemID   uint `gorm:"primaryKey"`
	Quantity int  `json:"quantity"`
}
