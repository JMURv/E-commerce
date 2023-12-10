package models

import (
	"errors"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  *uint     `json:"categoryID"`
	Category    *Category `json:"category" gorm:"foreignKey:CategoryID"`
	Sellers     []Seller  `json:"sellers" gorm:"many2many:seller_items;"`
	Tags        []Tag     `json:"tags" gorm:"many2many:item_tags;"`
}

func GetItemByID(id string) *Item {
	var getItem Item
	db.Preload("Tags").Where("ID=?", id).First(&getItem)
	return &getItem
}

func GetAllItems() []Item {
	var Items []Item
	db.Preload("Tags").Find(&Items)
	return Items
}

func (i *Item) CreateItem() (*Item, error) {

	if i.Name == "" {
		return i, errors.New("name is required")
	}
	if i.Description == "" {
		return i, errors.New("description is required")
	}
	if i.Price == 0 {
		return i, errors.New("price is required")
	}

	for idx := range i.Tags {
		existingTag := &Tag{}
		if err := db.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&i.Tags[idx]).Error; err != nil {
					return nil, err
				}
				if err := db.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		i.Tags[idx] = *existingTag
	}

	result := db.Create(&i)
	if result.Error != nil {
		return nil, result.Error
	}
	return i, nil
}

func (i *Item) UpdateItem(newData *Item) (*Item, error) {
	if newData.Name != "" {
		i.Name = newData.Name
	}

	if newData.Description != "" {
		i.Description = newData.Description
	}

	if newData.Price != 0 {
		i.Price = newData.Price
	}

	if newData.CategoryID != i.CategoryID {
		i.CategoryID = newData.CategoryID
	}
	db.Save(&i)
	return i, nil
}

func DeleteItem(id string) Item {
	var item Item
	db.Delete(&item, id)
	return item
}
