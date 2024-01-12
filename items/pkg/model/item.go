package model

import (
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func init() {
	var err error

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	dsn := os.Getenv("DSN")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db = database

	err = db.AutoMigrate(
		&Item{},
		&Category{},
		&Tag{},
	)
	if err != nil {
		log.Fatal(err)
	}

}

type Item struct {
	ID          uint64                `gorm:"primaryKey"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Price       float64               `json:"price"`
	CategoryID  uint64                `json:"categoryID"`
	Category    Category              `json:"category" gorm:"foreignKey:CategoryID"`
	UserID      uint64                `json:"userID"`
	Tags        []Tag                 `json:"tags" gorm:"many2many:item_tags;"`
	Status      string                `json:"status"`
	Quantity    int32                 `json:"quantity"`
	CreatedAt   timestamppb.Timestamp `json:"createdAt"`
	UpdatedAt   timestamppb.Timestamp `json:"updatedAt"`
}

func GetAllItems() ([]Item, error) {
	var Items []Item
	if err := db.Preload("Category").Preload("Tags").Preload("User").Preload("Reviews").Find(&Items).Error; err != nil {
		return nil, err
	}
	return Items, nil
}
