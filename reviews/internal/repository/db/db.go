package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DSN string

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	DSN = os.Getenv("DSN")
}

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
		&model.Review{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: db}
}

func (r *Repository) GetByID(_ context.Context, reviewID uint64) (*model.Review, error) {
	var review model.Review
	if err := r.conn.Where("ID=?", reviewID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// TODO: reviewed user id

func (r *Repository) GetReviewsByUserID(_ context.Context, userID uint64) (*[]model.Review, error) {
	var reviews []model.Review
	if err := r.conn.Where("UserID=?", userID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return &reviews, nil
}

func (r *Repository) AggregateUserRatingByID(_ context.Context, userID uint64) (float32, error) {
	var rating float32
	var reviews []model.Review

	if err := r.conn.Where("UserID=?", userID).Find(&reviews).Error; err != nil {
		return rating, err
	}

	var ratingSum int
	revCount := len(reviews)
	for i := range reviews {
		ratingSum += int(reviews[i].Rating)
	}

	if revCount > 0 {
		rating = float32(ratingSum / revCount)
	}
	return rating, nil
}

func (r *Repository) Create(_ context.Context, review *model.Review) (*model.Review, error) {
	if review.UserID == 0 {
		return nil, repo.ErrUserIDRequired
	}

	if review.ItemID == 0 && review.ReviewedUserID == 0 {
		return nil, repo.ErrItemIDRequired
	}

	if review.Rating == 0 {
		return nil, repo.ErrRatingRequired
	}

	if err := r.conn.Create(&review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

func (r *Repository) Update(ctx context.Context, reviewID uint64, newData *model.Review) (*model.Review, error) {
	review, err := r.GetByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	if newData.UserID != review.UserID {
		review.UserID = newData.UserID
	}

	if newData.ItemID != review.ItemID {
		review.ItemID = newData.ItemID
	}

	if newData.ReviewedUserID != review.ReviewedUserID {
		review.ReviewedUserID = newData.ReviewedUserID
	}

	if newData.Advantages != "" {
		review.Advantages = newData.Advantages
	}

	if newData.Disadvantages != "" {
		review.Disadvantages = newData.Disadvantages
	}

	if newData.ReviewText != "" {
		review.ReviewText = newData.ReviewText
	}

	if newData.Rating != 0 {
		review.Rating = newData.Rating
	}

	if err = r.conn.Save(&review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

func (r *Repository) Delete(_ context.Context, reviewID uint64) error {
	if err := r.conn.Delete(&model.Review{}, reviewID).Error; err != nil {
		return err
	}
	return nil
}
