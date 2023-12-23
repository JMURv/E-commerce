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

func GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := db.Preload("ParentCategory").Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func GetCategoryByID(categoryID uint) (*Category, error) {
	var category Category
	if err := db.Preload("ParentCategory").Where("ID = ?", categoryID).First(&category).Error; err != nil {
		return &category, err
	}
	return &category, nil
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

func UpdateCategory(categoryID uint, newData *Category) (*Category, error) {
	category, err := GetCategoryByID(categoryID)
	if err != nil {
		return category, err
	}

	if newData.Name != "" {
		category.Name = newData.Name
	}

	if newData.Description != "" {
		category.Description = newData.Description
	}

	if newData.ParentCategoryID != category.ParentCategoryID {
		category.ParentCategoryID = newData.ParentCategoryID
	}

	if err := db.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func DeleteCategory(categoryID uint) error {
	var category Category
	if err := db.Delete(&category, categoryID).Error; err != nil {
		return err
	}
	return nil
}
