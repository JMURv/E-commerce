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
		//&Category{},
		//&Tag{},
		//&User{},
		//&Item{},
		//&Review{},
		&Favorite{},
		&Room{},
		&Message{},
		&Notification{},
	)
	if err != nil {
		log.Fatal(err)
	}

}
