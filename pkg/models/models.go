package models

import (
	"e-commerce/pkg/config"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

type Item struct {
	gorm.Model
	ID          uint     `json:"id" gorm:"type:serial;primaryKey"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	CategoryID  *uint    `json:"categoryID"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
	Tags        []Tag    `gorm:"many2many:item_tags;"`
}

type Category struct {
	gorm.Model
	Name          string      `json:"name"`
	SubcategoryID uint        `json:"subcategoryID"`
	Subcategory   Subcategory `json:"subcategory"`
}

type Subcategory struct {
	gorm.Model
	Name string `json:"name"`
}

type Tag struct {
	gorm.Model
	Name string `json:"name"`
}

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
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
