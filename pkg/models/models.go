package models

import "gorm.io/gorm"

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
