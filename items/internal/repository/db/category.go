package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.GetAllCategories.repo")
	defer span.Finish()

	var categories []*model.Category
	if err := r.conn.Preload("ParentCategory").Find(&categories).Error; err != nil {
		return categories, err
	}
	return categories, nil
}

func (r *Repository) GetCategoryByID(ctx context.Context, categoryID uint64) (*model.Category, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.GetCategoryByID.repo")
	defer span.Finish()

	var category model.Category
	if err := r.conn.Preload("ParentCategory").Where("ID = ?", categoryID).First(&category).Error; err != nil {
		return &category, err
	}
	return &category, nil
}

func (r *Repository) CreateCategory(ctx context.Context, c *model.Category) (*model.Category, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.CreateCategory.repo")
	defer span.Finish()

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
	span, _ := opentracing.StartSpanFromContext(ctx, "items.UpdateCategory.repo")
	defer span.Finish()

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

func (r *Repository) DeleteCategory(ctx context.Context, categoryID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.DeleteCategory.repo")
	defer span.Finish()

	var category model.Category
	if err := r.conn.Delete(&category, categoryID).Error; err != nil {
		return err
	}
	return nil
}
