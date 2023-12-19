package models

import (
	"errors"
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	User   User `json:"user" gorm:"foreignKey:UserID"`
	UserID uint `json:"userID"`
	Item   Item `json:"item" gorm:"foreignKey:ItemID"`
	ItemID uint `json:"itemID"`
}

func GetAllUserFavorites(userID uint) ([]Favorite, error) {
	var favorites []Favorite
	if err := db.Preload("User").Preload("Item").Where("favorites.user_id = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func GetFavoriteByID(favoriteID uint) (*Favorite, error) {
	var favorite Favorite
	if err := db.Preload("User").Preload("Item").Where("ID = ?", favoriteID).First(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (f *Favorite) CreateFavorite() (*Favorite, error) {
	if f.UserID == 0 {
		return nil, errors.New("userID is required")
	}
	if f.ItemID == 0 {
		return nil, errors.New("itemID is required")
	}

	if err := db.Create(&f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func DeleteFavorite(favoriteID uint) error {
	var favorite Favorite
	if err := db.Delete(&favorite, favoriteID).Error; err != nil {
		return err
	}
	return nil
}
