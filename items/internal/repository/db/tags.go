package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
)

func (r *Repository) ListTags(_ context.Context) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.conn.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repository) CreateTag(_ context.Context, t *model.Tag) (*model.Tag, error) {
	if t.Name == "" {
		return nil, repo.ErrTagNameIsRequired
	}

	if err := r.conn.Create(&t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *Repository) DeleteTag(_ context.Context, tagName string) error {
	if err := r.conn.Where("name=?", tagName).Delete(&model.Tag{}).Error; err != nil {
		return err
	}
	return nil
}
