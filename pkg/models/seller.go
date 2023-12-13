package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SellerItem struct {
	gorm.Model
	SellerID  uint64    `json:"sellerID"`
	Seller    Seller    `json:"seller" gorm:"foreignKey:SellerID;references:ID"`
	ItemID    uint64    `json:"itemID"`
	Item      Item      `json:"item" gorm:"foreignKey:ItemID;references:ID"`
	Quantity  int32     `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

type Seller struct {
	gorm.Model
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	SellerItems []SellerItem `json:"sellerItems"`
}

func GetSellerByID(sellerId string) *Seller {
	var seller Seller
	db.Preload("SellerItems").Preload("SellerItems.Item").Preload("SellerItems.Seller").Where("ID=?", sellerId).First(&seller)
	return &seller
}

func GetSellerWithItems(sellerId string) *Seller {
	var seller Seller
	if err := db.Preload("SellerItems").Preload("SellerItems.Item").Preload("SellerItems.Seller").Where("id = ?", sellerId).First(&seller).Error; err != nil {
		fmt.Println("Error retrieving seller:", err)
		return nil
	}
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
	result := db.Create(&s)
	if result.Error != nil {
		return nil, result.Error
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

func (s *SellerItem) LinkItemToSeller() (*SellerItem, error) {

	if s.Quantity == 0 {
		return s, errors.New("quantity is required")
	}

	s.CreatedAt = time.Now()

	result := db.Create(&s)
	if result.Error != nil {
		return nil, result.Error
	}

	return s, nil
}

func (s *SellerItem) UpdateSellerItem(newData *SellerItem) (*SellerItem, error) {

	if newData.Quantity != s.Quantity {
		s.Quantity = newData.Quantity
	}

	if err := db.Save(&s).Error; err != nil {
		return nil, err
	}
	return s, nil
}
