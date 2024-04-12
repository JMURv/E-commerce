package db

import (
	"context"
	"errors"
	"fmt"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	conf "github.com/JMURv/e-commerce/items/pkg/config"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Repository struct {
	conn *gorm.DB
}

func New(conf *conf.Config) *Repository {
	DSN := fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Database,
	)

	conn, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = conn.AutoMigrate(
		&model.Item{},
		&model.Category{},
		&model.Tag{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: conn}
}

func (r *Repository) ListItem(_ context.Context) ([]*model.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetItemByID(_ context.Context, id uint64) (*model.Item, error) {
	var getItem model.Item
	if err := r.conn.Preload("Category").Preload("Tags").Where("ID=?", id).First(&getItem).Error; err != nil {
		return nil, repo.ErrNotFound
	}
	return &getItem, nil

}

func (r *Repository) ListUserItemsByID(_ context.Context, userID uint64) ([]*model.Item, error) {
	var items []*model.Item
	if err := r.conn.Preload("Category").Preload("Tags").Where("UserID=?", userID).Find(&items).Error; err != nil {
		return nil, repo.ErrNotFound
	}
	return items, nil
}

func (r *Repository) CreateItem(_ context.Context, i *model.Item) (*model.Item, error) {
	if i.UserID == 0 {
		return nil, repo.ErrUserIDRequired
	}
	if i.CategoryID == 0 {
		return nil, repo.ErrCategoryIDRequired
	}

	if i.Name == "" {
		return nil, repo.ErrNameRequired
	}
	if i.Description == "" {
		return nil, repo.ErrDescriptionRequired
	}
	if i.Price == 0 {
		return nil, repo.ErrPriceRequired
	}
	if i.Quantity == 0 {
		i.Quantity = 1
	}

	// Link or create tags if specified
	for idx := range i.Tags {
		existingTag := &model.Tag{}
		if err := r.conn.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := r.conn.Create(&i.Tags[idx]).Error; err != nil {
					return nil, err
				}
				if err := r.conn.Where("name = ?", i.Tags[idx].Name).First(existingTag).Error; err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		i.Tags[idx] = existingTag
	}

	i.Status = "created"
	i.CreatedAt = *timestamppb.Now()
	// Perform item's save
	if err := r.conn.Create(&i).Error; err != nil {
		return nil, err
	}

	return i, nil
}

func (r *Repository) UpdateItem(ctx context.Context, itemID uint64, newData *model.Item) (*model.Item, error) {
	i, err := r.GetItemByID(ctx, itemID)
	if err != nil {
		return i, err
	}

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

	if err = r.conn.Save(&i).Error; err != nil {
		return nil, err
	}
	return i, nil
}

func (r *Repository) DeleteItem(_ context.Context, itemID uint64) error {
	var item model.Item

	if err := r.conn.Delete(&item, itemID).Error; err != nil {
		return err
	}
	return nil
}
