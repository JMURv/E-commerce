package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Item struct {
	ID          uint64                `gorm:"primaryKey"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Price       float64               `json:"price"`
	CategoryID  uint64                `json:"categoryID"`
	Category    Category              `json:"category" gorm:"foreignKey:CategoryID"`
	UserID      uint64                `json:"userID"`
	Tags        []Tag                 `json:"tags" gorm:"many2many:item_tags;"`
	Status      string                `json:"status"`
	Quantity    int32                 `json:"quantity"`
	CreatedAt   timestamppb.Timestamp `json:"createdAt"`
	UpdatedAt   timestamppb.Timestamp `json:"updatedAt"`
}
