package model

import (
	"github.com/JMURv/protos/ecom/common"
)

func tagsToProto(tags []Tag) []*common.Tag {
	var pbTags []*common.Tag
	for _, tag := range tags {
		pbTag := &common.Tag{
			Name: tag.Name,
		}
		pbTags = append(pbTags, pbTag)
	}
	return pbTags
}

func TagsFromProto(tags []*common.Tag) []Tag {
	var modelTags []Tag
	for _, tag := range tags {
		modelTag := Tag{
			Name: tag.Name,
		}
		modelTags = append(modelTags, modelTag)
	}
	return modelTags
}

func ItemFromProto(i *common.Item) *Item {
	return &Item{
		ID:          i.Id,
		Name:        i.Name,
		Description: i.Description,
		Price:       i.Price,
		CategoryID:  i.CategoryId,
		Category: Category{
			ID:               i.Category.Id,
			Name:             i.Category.Name,
			Description:      i.Category.Description,
			ParentCategoryID: i.Category.ParentCategoryId,
			CreatedAt:        *i.Category.CreatedAt,
			UpdatedAt:        *i.Category.UpdatedAt,
		},
		UserID:    i.UserId,
		Tags:      TagsFromProto(i.Tags),
		Status:    i.Status,
		Quantity:  i.Quantity,
		CreatedAt: *i.CreatedAt,
		UpdatedAt: *i.UpdatedAt,
	}
}

func ItemToProto(i *Item) *common.Item {
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
			ParentCategoryId: i.Category.ParentCategoryID,
			CreatedAt:        &i.Category.CreatedAt,
			UpdatedAt:        &i.Category.UpdatedAt,
		},
		Tags:      tagsToProto(i.Tags),
		Status:    i.Status,
		Quantity:  i.Quantity,
		CreatedAt: &i.CreatedAt,
		UpdatedAt: &i.UpdatedAt,
	}
}
