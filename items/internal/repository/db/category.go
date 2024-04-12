package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
)

func (r *Repository) GetAllCategories(_ context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.conn.Preload("ParentCategory").Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func (r *Repository) GetCategoryByID(_ context.Context, categoryID uint64) (*model.Category, error) {
	var category model.Category
	if err := r.conn.Preload("ParentCategory").Where("ID = ?", categoryID).First(&category).Error; err != nil {
		return &category, err
	}
	return &category, nil
}

func (r *Repository) CreateCategory(_ context.Context, c *model.Category) (*model.Category, error) {
	if c.Name == "" {
		return nil, repo.ErrCategoryNameIsRequired
	}
	if c.Description == "" {
		return nil, repo.ErrCategoryDescriptionIsRequired
	}
	result := r.conn.Create(&c)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}

func (r *Repository) UpdateCategory(ctx context.Context, categoryID uint64, newData *model.Category) (*model.Category, error) {
	category, err := r.GetCategoryByID(ctx, categoryID)
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

	if err := r.conn.Save(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func (r *Repository) DeleteCategory(_ context.Context, categoryID uint64) error {
	var category model.Category
	if err := r.conn.Delete(&category, categoryID).Error; err != nil {
		return err
	}
	return nil
}
