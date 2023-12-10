package models

import (
	"errors"
	"gorm.io/gorm"
)

type Seller struct {
	gorm.Model
	Username  string `json:"username"`
	Email     string `json:"email"`
	SoldItems []Item `json:"soldItems" gorm:"many2many:seller_items;"`
}

func GetSellerByID(sellerId string) *Seller {
	var seller Seller
	db.Where("ID=?", sellerId).First(&seller)
	return &seller
}

func GetAllSellers() []Seller {
	var sellers []Seller
	db.Find(&sellers)
	return sellers
}

func (s *Seller) CreateSeller() (*Seller, error) {
	if s.Username == "" {
		return s, errors.New("username is required")
	}
	if s.Email == "" {
		return s, errors.New("email is required")
	}
	return s, nil
}

func UpdateSeller(sellerId string, newData *Seller) (*Seller, error) {
	seller := GetSellerByID(sellerId)
	if newData.Username != "" {
		seller.Username = newData.Username
	}

	if newData.Email != "" {
		seller.Email = newData.Email
	}
	db.Save(&seller)
	return seller, nil
}

func DeleteSeller(id string) Seller {
	var seller Seller
	db.Delete(&seller, id)
	return seller
}
