package model

import (
	"e-commerce/pkg/models"
	"errors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func init() {
	var err error

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	log.Println("Env file has been loaded")

	dsn := os.Getenv("DSN")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")

	db = database

	err = db.AutoMigrate(
		&Review{},
	)
	if err != nil {
		log.Fatal(err)
	}

}

type Review struct {
	gorm.Model
	ID             uint64      `gorm:"primaryKey"`
	UserID         uint64      `json:"user_id"`
	User           models.User `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ItemID         uint64      `json:"item_id"`
	Item           models.Item `json:"item" gorm:"foreignKey:ItemID"`
	ReviewedUserID uint64      `json:"reviewed_user_id"`
	ReviewedUser   models.User `json:"reviewedUser" gorm:"foreignKey:ReviewedUserID;references:ID"`
	Advantages     string      `json:"advantages"`
	Disadvantages  string      `json:"disadvantages"`
	ReviewText     string      `json:"review_text"`
	Rating         uint32      `json:"rating"`
}

func GetReviewByID(reviewID uint64) (*Review, error) {
	var review Review
	if err := db.Preload("User").Preload("Item").Preload("ReviewedUser").Where("ID=?", reviewID).First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *Review) CreateReview() (*Review, error) {
	if r.UserID == 0 {
		return nil, errors.New("userID is required")
	}

	if r.ItemID == 0 && r.ReviewedUserID == 0 {
		return nil, errors.New("itemID or ReviewedUserID is required")
	}

	if r.Rating == 0 {
		return nil, errors.New("rating is required")
	}

	if err := db.Create(&r).Error; err != nil {
		return nil, err
	}

	return r, nil
}

func UpdateReview(reviewID uint64, newData *Review) (*Review, error) {
	review, err := GetReviewByID(reviewID)
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

	if err := db.Save(&review).Error; err != nil {
		return nil, err
	}
	return review, nil
}

func DeleteReview(reviewID uint64) error {
	var review Review
	if err := db.Delete(&review, reviewID).Error; err != nil {
		return err
	}
	return nil
}
