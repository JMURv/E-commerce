package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID   uint        `json:"userID"`
	Items    []OrderItem `json:"items" gorm:"many2many:order_items;"`
	Status   string      `json:"status"`
	Payment  string      `json:"payment"`
	Shipping string      `json:"shipping"`
}

type OrderItem struct {
	OrderID  uint `gorm:"primaryKey"`
	ItemID   uint `gorm:"primaryKey"`
	Quantity int  `json:"quantity"`
}
