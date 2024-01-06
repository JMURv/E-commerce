package model

import (
	"errors"
	"time"
)

type Category struct {
	ID               uint64    `gorm:"primaryKey"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ParentCategoryID *uint64   `json:"parent_category_id"`
	ParentCategory   *Category `json:"parent_category" gorm:"foreignKey:ParentCategoryID"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := db.Preload("ParentCategory").Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func GetCategoryByID(categoryID uint64) (*Category, error) {
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

func UpdateCategory(categoryID uint64, newData *Category) (*Category, error) {
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

func DeleteCategory(categoryID uint64) error {
	var category Category
	if err := db.Delete(&category, categoryID).Error; err != nil {
		return err
	}
	return nil
}
