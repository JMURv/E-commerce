package models

import (
	"time"
)

type Item struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  uint      `json:"categoryID"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	UserID      uint      `json:"userID"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Reviews     []Review  `json:"reviews"`
	Tags        []Tag     `json:"tags" gorm:"many2many:item_tags;"`
	Status      string    `json:"status"`
	Quantity    int32     `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
