package models

import (
	"e-commerce/pkg/config"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

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
