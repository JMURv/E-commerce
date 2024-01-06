package model

import (
	"errors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
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
	ID          uint64    `gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  uint64    `json:"categoryID"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	UserID      uint64    `json:"userID"`
	Tags        []Tag     `json:"tags" gorm:"many2many:item_tags;"`
	Status      string    `json:"status"`
	Quantity    int32     `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func GetItemByID(id uint64) (*Item, error) {
	var getItem Item
	if err := db.Preload("Category").Preload("Tags").Where("ID=?", id).First(&getItem).Error; err != nil {
		return nil, err
	}
	return &getItem, nil
}

func GetAllItems() ([]Item, error) {
	var Items []Item
	if err := db.Preload("Category").Preload("Tags").Preload("User").Preload("User").Preload("Reviews").Find(&Items).Error; err != nil {
		return nil, err
	}
	return Items, nil
}

func (i *Item) CreateItem() (*Item, error) {
	if i.UserID == 0 {
		return i, errors.New("UserID is required")
	}
	if i.CategoryID == 0 {
		return i, errors.New("CategoryID is required")
	}

	if i.Name == "" {
		return i, errors.New("name is required")
	}
	if i.Description == "" {
		return i, errors.New("description is required")
	}
	if i.Price == 0 {
		return i, errors.New("price is required")
	}
	if i.Quantity == 0 {
		i.Quantity = 1
	}

	// Link or create tags if specified
	for idx := range i.Tags {
		existingTag := &Tag{}
		if err := db.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&i.Tags[idx]).Error; err != nil {
					return nil, err
				}
				if err := db.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		i.Tags[idx] = *existingTag
	}

	i.Status = "created"
	i.CreatedAt = time.Now()
	// Perform item's save
	if err := db.Create(&i).Error; err != nil {
		return nil, err
	}

	return i, nil
}

func UpdateItem(itemID uint64, newData *Item) (*Item, error) {
	i, err := GetItemByID(itemID)
	if err != nil {
		return i, err
	}

	if newData.Name != "" {
		i.Name = newData.Name
	}

	if newData.Description != "" {
		i.Description = newData.Description
	}

	if newData.Price != 0 {
		i.Price = newData.Price
	}

	if newData.CategoryID != 0 {
		i.CategoryID = newData.CategoryID
	}
	if i.Quantity != 0 {
		i.Quantity = newData.Quantity
	}

	if err := db.Save(&i).Error; err != nil {
		return nil, err
	}
	return i, nil
}

func DeleteItem(id uint64) error {
	var item Item

	if err := db.Delete(&item, id).Error; err != nil {
		return err
	}
	return nil
}
