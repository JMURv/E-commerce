package main

import (
	"e-commerce/items/model"
	"github.com/JMURv/protos/ecom/common"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func responseItem(i *model.Item) (*common.Item, error) {
	return &common.Item{
		Id:         i.ID,
		UserId:     i.UserID,
		CategoryId: i.CategoryID,

		Name:        i.Name,
		Description: i.Description,
		Price:       i.Price,
		Category: &common.Category{
			Id:               i.Category.ID,
			Name:             i.Category.Name,
			Description:      i.Category.Description,
			ParentCategoryId: *i.Category.ParentCategoryID,
			CreatedAt:        timestamppb.New(i.Category.CreatedAt),
			UpdatedAt:        timestamppb.New(i.Category.UpdatedAt),
		},
		Tags:      tagsToPb(i.Tags),
		Status:    i.Status,
		Quantity:  i.Quantity,
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
	}, nil
}

func tagsToModel(tags []*common.Tag) []model.Tag {
	var modelTags []model.Tag
	for _, tag := range tags {
		modelTag := model.Tag{
			Name: tag.Name,
		}
		modelTags = append(modelTags, modelTag)
	}
	return modelTags
}

func tagsToPb(tags []model.Tag) []*common.Tag {
	var pbTags []*common.Tag
	for _, tag := range tags {
		pbTag := &common.Tag{
			Name: tag.Name,
		}
		pbTags = append(pbTags, pbTag)
	}
	return pbTags
}
