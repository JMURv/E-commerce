package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/opentracing/opentracing-go"
)

func (r *Repository) ListTags(ctx context.Context) ([]*model.Tag, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.ListTags.repo")
	defer span.Finish()

	var tags []*model.Tag
	if err := r.conn.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repository) CreateTag(ctx context.Context, t *model.Tag) (*model.Tag, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.CreateTag.repo")
	defer span.Finish()

	if t.Name == "" {
		return nil, repo.ErrTagNameIsRequired
	}

	if err := r.conn.Create(&t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *Repository) DeleteTag(ctx context.Context, tagName string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "items.DeleteTag.repo")
	defer span.Finish()

	if err := r.conn.Where("name=?", tagName).Delete(&model.Tag{}).Error; err != nil {
		return err
	}
	return nil
}
