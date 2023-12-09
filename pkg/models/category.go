package models

import (
	"errors"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ParentCategoryID *uint     `json:"parentCategoryID"`
	ParentCategory   *Category `json:"parentCategory" gorm:"foreignKey:ParentCategoryID"`
}

func (c *Category) CreateCategory() (*Category, error) {
	if c.Name == "" {
		return c, errors.New("name is required")
	}
	if c.Description == "" {
		return c, errors.New("description is required")
	}
	result := db.Create(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}

func GetAllCategories() []Category {
	var Categories []Category
	db.Preload("ParentCategory").Find(&Categories)
	return Categories
}
