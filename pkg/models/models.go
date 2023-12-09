package models

import (
	"e-commerce/pkg/config"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

type Favorite struct {
	gorm.Model
	Author         User `json:"author" gorm:"foreignKey:AuthorID"`
	AuthorID       uint `json:"authorID"`
	FavoriteItem   Item `json:"favoriteItem" gorm:"foreignKey:FavoriteItemID"`
	FavoriteItemID uint `json:"favoriteItemID"`
}

type Review struct {
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

type Tag struct {
	gorm.Model
	Name string `json:"name"`
}

func init() {
	var err error
	config.Connect()
	db = config.GetDB()

	err = db.AutoMigrate(
		&Item{},
		&Category{},
		&Tag{},
		&User{},
		&Seller{},
		&Cart{},
		&CartItem{},
		&Order{},
		&OrderItem{},
		&Favorite{},
	)
	if err != nil {
		log.Fatal(err)
	}
}
