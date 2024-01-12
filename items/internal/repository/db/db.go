package db

import (
	"context"
	"errors"
	repo "github.com/JMURv/e-commerce/items/internal/repository"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DSN string

type Repository struct {
	conn *gorm.DB
}

func New() *Repository {
	var err error
	var db *gorm.DB

	db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&model.Item{},
		&model.Category{},
		&model.Tag{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: db}
}

func (r *Repository) GetByID(_ context.Context, id uint64) (*model.Item, error) {
	var getItem model.Item
	if err := r.conn.Preload("Category").Preload("Tags").Where("ID=?", id).First(&getItem).Error; err != nil {
		return nil, repo.ErrNotFound
	}
	return &getItem, nil

}

func (r *Repository) Create(_ context.Context, i *model.Item) (*model.Item, error) {
	if i.UserID == 0 {
		return i, repo.ErrUserIDRequired
	}
	if i.CategoryID == 0 {
		return i, repo.ErrCategoryIDRequired
	}

	if i.Name == "" {
		return i, repo.ErrNameRequired
	}
	if i.Description == "" {
		return i, repo.ErrDescriptionRequired
	}
	if i.Price == 0 {
		return i, repo.ErrPriceRequired
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
		i.Tags[idx] = *existingTag
	}

	i.Status = "created"
	i.CreatedAt = *timestamppb.Now()
	// Perform item's save
	if err := r.conn.Create(&i).Error; err != nil {
		return nil, err
	}

	return i, nil
}

func (r *Repository) Update(ctx context.Context, itemID uint64, newData *model.Item) (*model.Item, error) {
	i, err := r.GetByID(ctx, itemID)
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

func (r *Repository) Delete(_ context.Context, itemID uint64) error {
	var item model.Item

	if err := r.conn.Delete(&item, itemID).Error; err != nil {
		return err
	}
	return nil
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	DSN = os.Getenv("DSN")
}
