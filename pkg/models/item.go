package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Item struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CategoryID  uint      `json:"categoryID"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	UserID      uint      `json:"userID"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Reviews     []Review  `json:"reviews"`
	Tags        []Tag     `json:"tags" gorm:"many2many:item_tags;"`
	Quantity    int32     `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

func GetItemByID(id uint) *Item {
	var getItem Item
	db.Preload("Category").Preload("Tags").Preload("User").Preload("Reviews").Where("ID=?", id).First(&getItem)
	return &getItem
}

func GetAllItems() []Item {
	var Items []Item
	db.Preload("Category").Preload("Tags").Preload("User").Preload("User").Preload("Reviews").Find(&Items)
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
	if i.Quantity == 0 {
		i.Quantity = 1
	}

	// Link or create tags if specified
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

	i.CreatedAt = time.Now()
	// Perform item's save
	if err := db.Create(&i).Error; err != nil {
		return nil, err
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

	if newData.CategoryID != 0 {
		i.CategoryID = newData.CategoryID
	}
	if i.Quantity != 0 {
		i.Quantity = newData.Quantity
	}

	if err := db.Save(&i).Error; err != nil {
		return nil, err
	}
	return i, nil
}

func DeleteItem(id uint) (*Item, error) {
	var item Item

	if err := db.Delete(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
