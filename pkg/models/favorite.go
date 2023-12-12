package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	Author         User `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID       uint `json:"authorID"`
	FavoriteItem   Item `json:"favoriteItem" gorm:"foreignKey:FavoriteItemID"`
	FavoriteItemID uint `json:"favoriteItemID"`
}
