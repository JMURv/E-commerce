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

func GetAllFavorites() ([]Favorite, error) {
	var favorites []Favorite
	if err := db.Preload("User").Preload("Item").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func GetAllUserFavorites(userID uint) ([]Favorite, error) {
	var favorites []Favorite
	if err := db.Preload("User").Preload("Item").Where("UserID = ?", userID).Find(&favorites).Error; err != nil {
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

func UpdateFavorite(favoriteID uint, newData *Favorite) (*Favorite, error) {
	favorite, err := GetFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	if newData.UserID != 0 {
		favorite.UserID = newData.UserID
	}

	if newData.ItemID != 0 {
		favorite.ItemID = newData.ItemID
	}

	if err := db.Save(&favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

func DeleteFavorite(favoriteID uint) (Favorite, error) {
	var favorite Favorite
	if err := db.Delete(&favorite, favoriteID).Error; err != nil {
		return favorite, err
	}
	return favorite, nil
}
