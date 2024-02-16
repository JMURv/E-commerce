package model

import (
	cm "github.com/JMURv/e-commerce/api/pb/common"
)

func tagsToProto(tags []Tag) []*cm.Tag {
	var pbTags []*cm.Tag
	for _, tag := range tags {
		pbTag := &cm.Tag{
			Name: tag.Name,
		}
		pbTags = append(pbTags, pbTag)
	}
	return pbTags
}

func TagsFromProto(tags []*cm.Tag) []Tag {
	t := make([]Tag, 0, len(tags))
	for i := range tags {
		t = append(t, Tag{
			Name: tags[i].Name,
		})
	}
	return t
}

func ItemFromProto(i *cm.Item) *Item {
	return &Item{
		ID:          i.Id,
		Name:        i.Name,
		Description: i.Description,
		Price:       i.Price,
		CategoryID:  i.CategoryId,
		//Category: Category{
		//	ID:               i.Category.Id,
		//	Name:             i.Category.Name,
		//	Description:      i.Category.Description,
		//	ParentCategoryID: i.Category.ParentCategoryId,
		//},
		UserID:   i.UserId,
		Tags:     TagsFromProto(i.Tags),
		Status:   i.Status,
		Quantity: i.Quantity,
	}
}

func ItemToProto(i *Item) *cm.Item {
	return &cm.Item{
		Id:         i.ID,
		UserId:     i.UserID,
		CategoryId: i.CategoryID,

		Name:        i.Name,
		Description: i.Description,
		Price:       i.Price,
		Category: &cm.Category{
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

func ItemsToProto(items []*Item) []*cm.Item {
	res := make([]*cm.Item, 0, len(items))
	for i := range items {
		res = append(res, ItemToProto(items[i]))
	}

	return res
}

func CategoryFromProto(c *cm.Category) *Category {
	return &Category{
		Name:             c.Name,
		Description:      c.Description,
		ParentCategoryID: c.ParentCategoryId,
	}
}

func CategoryToProto(c *Category) *cm.Category {
	return &cm.Category{
		Id:               c.ID,
		Name:             c.Name,
		Description:      c.Description,
		ParentCategoryId: c.ParentCategoryID,
	}
}

func CategoriesToProto(categories []*Category) []*cm.Category {
	res := make([]*cm.Category, 0, len(categories))
	for i := range categories {
		res = append(res, CategoryToProto(categories[i]))
	}

	return res
}
