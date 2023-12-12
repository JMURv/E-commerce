package models

import "gorm.io/gorm"

type ItemReview struct {
	gorm.Model
	Author         User   `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID       uint   `json:"authorID"`
	ReviewedItem   Item   `json:"reviewedItem" gorm:"foreignKey:ReviewedItemID"`
	ReviewedItemID uint   `json:"reviewedItemID"`
	Advantages     string `json:"advantages"`
	Disadvantages  string `json:"disadvantages"`
	ReviewText     string `json:"reviewText"`
	Rating         int    `json:"rating"`
}

type SellerReview struct {
	gorm.Model
	Author           User   `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID         uint   `json:"authorID"`
	ReviewedSeller   Item   `json:"reviewedSeller" gorm:"foreignKey:ReviewedItemID"`
	ReviewedSellerID uint   `json:"reviewedSellerID"`
	Advantages       string `json:"advantages"`
	Disadvantages    string `json:"disadvantages"`
	ReviewText       string `json:"reviewText"`
	Rating           int    `json:"rating"`
}
