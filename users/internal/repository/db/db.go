package db

import (
	"context"
	repo "github.com/JMURv/e-commerce/users/internal/repository"
	"github.com/JMURv/e-commerce/users/pkg/model"
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
		&model.User{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: db}
}

func (r *Repository) GetByID(_ context.Context, userID uint64) (*model.User, error) {
	var user model.User
	if err := r.conn.Where("ID=?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByEmail(_ context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.conn.Where("Email=?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Create(_ context.Context, u *model.User) (*model.User, error) {
	if u.Username == "" {
		return u, repo.ErrUsernameIsRequired
	}

	if u.Email == "" {
		return u, repo.ErrEmailIsRequired
	}

	if err := r.conn.Create(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) Update(ctx context.Context, userID uint64, newData *model.User) (*model.User, error) {
	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return user, err
	}

	if newData.Username != "" {
		user.Username = newData.Username
	}

	if newData.Email != "" {
		user.Email = newData.Email
	}
	r.conn.Save(&user)
	return user, nil
}

func (r *Repository) Delete(_ context.Context, userID uint64) error {
	var user model.User
	if err := r.conn.Delete(&user, userID).Error; err != nil {
		return err
	}
	return nil
}
