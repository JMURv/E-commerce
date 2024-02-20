package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/favorites/internal/repository"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type Repository struct {
	conn *gorm.DB
}

func New() *Repository {
	var err error
	var db *gorm.DB

	if err = godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err = gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err = db.AutoMigrate(&model.Favorite{}); err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: db}
}

func (r *Repository) GetAllUserFavorites(_ context.Context, userID uint64) ([]*model.Favorite, error) {
	var favorites []*model.Favorite
	if err := r.conn.Where("UserID = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (r *Repository) GetFavoriteByID(_ context.Context, favoriteID uint64) (*model.Favorite, error) {
	var favorite model.Favorite
	if err := r.conn.Where("ID = ?", favoriteID).First(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *Repository) CreateFavorite(_ context.Context, favData *model.Favorite) (*model.Favorite, error) {
	if favData.UserID == 0 {
		return nil, repo.ErrUserIDRequired
	}
	if favData.ItemID == 0 {
		return nil, repo.ErrItemIDRequired
	}

	if err := r.conn.Create(favData).Error; err != nil {
		return nil, err
	}
	return favData, nil
}

func (r *Repository) DeleteFavorite(_ context.Context, favoriteID uint64) error {
	var favorite model.Favorite
	if err := r.conn.Delete(&favorite, favoriteID).Error; err != nil {
		return err
	}
	return nil
}
