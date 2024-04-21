package db

import (
	"context"
	"fmt"
	repo "github.com/JMURv/e-commerce/reviews/internal/repository"
	conf "github.com/JMURv/e-commerce/reviews/pkg/config"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
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

	err = conn.AutoMigrate(
		&model.Review{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: conn}
}

func (r *Repository) GetByID(ctx context.Context, reviewID uint64) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.GetByID.repo")
	defer span.Finish()

	var review model.Review
	if err := r.conn.Where("ID=?", reviewID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// TODO: reviewed user id

func (r *Repository) GetReviewsByUserID(ctx context.Context, userID uint64) (*[]*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.GetReviewsByUserID.repo")
	defer span.Finish()

	var reviews []*model.Review
	if err := r.conn.Where("UserID=?", userID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return &reviews, nil
}

func (r *Repository) AggregateUserRatingByID(ctx context.Context, userID uint64) (float32, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.AggregateUserRatingByID.repo")
	defer span.Finish()

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

func (r *Repository) Create(ctx context.Context, review *model.Review) (*model.Review, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.Create.repo")
	defer span.Finish()

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
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.Update.repo")
	defer span.Finish()

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

func (r *Repository) Delete(ctx context.Context, reviewID uint64) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "reviews.Delete.repo")
	defer span.Finish()

	if err := r.conn.Delete(&model.Review{}, reviewID).Error; err != nil {
		return err
	}
	return nil
}
