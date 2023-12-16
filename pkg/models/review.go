package models

import (
	"errors"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	User          User   `json:"user" gorm:"foreignKey:UserID"`
	UserID        uint   `json:"userID"`
	Item          Item   `json:"item" gorm:"foreignKey:ItemID"`
	ItemID        uint   `json:"itemID"`
	Seller        Seller `json:"seller" gorm:"foreignKey:SellerID"`
	SellerID      uint   `json:"sellerID"`
	Advantages    string `json:"advantages"`
	Disadvantages string `json:"disadvantages"`
	ReviewText    string `json:"reviewText"`
	Rating        uint   `json:"rating"`
}

func GetReviewByID(reviewID uint) (*Review, error) {
	var review Review
	if err := db.Preload("User").Preload("Item").Preload("Seller").Where("ID=?", reviewID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *Review) CreateReview() (*Review, error) {
	if r.UserID == 0 {
		return nil, errors.New("userID is required")
	}

	if r.ItemID == 0 && r.SellerID == 0 {
		return nil, errors.New("itemID or sellerID is required")
	}

	if r.Rating == 0 {
		return nil, errors.New("rating is required")
	}

	if err := db.Create(&r).Error; err != nil {
		return nil, err
	}

	return r, nil
}

func UpdateReview(reviewID uint, newData *Review) (*Review, error) {
	review, err := GetReviewByID(reviewID)
	if err != nil {
		return nil, err
	}

	if newData.UserID != review.UserID {
		review.UserID = newData.UserID
	}

	if newData.ItemID != review.ItemID {
		review.ItemID = newData.ItemID
	}

	if newData.SellerID != review.SellerID {
		review.SellerID = newData.SellerID
	}

	if newData.Advantages != "" {
		review.Advantages = newData.Advantages
	}

	if newData.Disadvantages != "" {
		review.Disadvantages = newData.Disadvantages
	}

	if newData.ReviewText != "" {
		review.ReviewText = newData.ReviewText
	}

	if newData.Rating != 0 {
		review.Rating = newData.Rating
	}

	if err := db.Save(&review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

func DeleteReview(reviewID uint) (Review, error) {
	var review Review
	if err := db.Delete(&review, reviewID).Error; err != nil {
		return review, err
	}
	return review, nil
}
