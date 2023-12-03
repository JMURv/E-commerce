package models

import (
	"e-commerce/pkg/config"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

type Item struct {
	gorm.Model
	ID          string `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	Price       float64
	CategoryID  uint
	Category    Category
	Tags        []Tag `gorm:"many2many:item_tags;"`
}

type Category struct {
	gorm.Model
	ID            uint
	Name          string
	SubcategoryID uint
	Subcategory   Subcategory
}

type Subcategory struct {
	gorm.Model
	ID   uint
	Name string
}

type Tag struct {
	gorm.Model
	ID   uint
	Name string
}

type User struct {
	gorm.Model
	ID       string `gorm:"type:uuid;primaryKey"`
	Username string
	Email    string
}

func init() {
	var err error
	config.Connect()
	db = config.GetDB()

	err = db.AutoMigrate(&Item{}, &Category{}, &Tag{})
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Item) CreateItem() *Item {
	db.Create(&i)
	return i
}

func GetAllItems() []Item {
	var Items []Item
	db.Find(&Items)
	return Items
}

func GetItemByID(id int64) (*Item, *gorm.DB) {
	var getItem Item
	db := db.Where("ID=?", id).Find(&getItem)
	return &getItem, db
}

func DeleteItem(id int64) Item {
	var item Item
	db.Delete(&item, id)
	return item
}
