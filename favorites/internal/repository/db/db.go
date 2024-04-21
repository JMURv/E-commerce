package db

import (
	"context"
	"fmt"
	repo "github.com/JMURv/e-commerce/favorites/internal/repository"
	conf "github.com/JMURv/e-commerce/favorites/pkg/config"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
	"github.com/opentracing/opentracing-go"
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

	if err = conn.AutoMigrate(&model.Favorite{}); err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: conn}
}

func (r *Repository) GetAllUserFavorites(ctx context.Context, userID uint64) ([]*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.GetAllUserFavorites.repo")
	defer span.Finish()

	var favorites []*model.Favorite
	if err := r.conn.Where("UserID = ?", userID).Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}

func (r *Repository) GetFavoriteByID(ctx context.Context, favoriteID uint64) (*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.GetFavoriteByID.repo")
	defer span.Finish()

	var favorite model.Favorite
	if err := r.conn.Where("ID = ?", favoriteID).First(&favorite).Error; err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *Repository) CreateFavorite(ctx context.Context, favData *model.Favorite) (*model.Favorite, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.CreateFavorite.repo")
	defer span.Finish()

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

func (r *Repository) DeleteFavorite(ctx context.Context, favoriteID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "favorites.DeleteFavorite.repo")
	defer span.Finish()

	var favorite model.Favorite
	if err := r.conn.Delete(&favorite, favoriteID).Error; err != nil {
		return err
	}
	return nil
}
