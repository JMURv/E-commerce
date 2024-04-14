package db

import (
	"context"
	"fmt"
	repo "github.com/JMURv/e-commerce/users/internal/repository"
	conf "github.com/JMURv/e-commerce/users/pkg/config"
	"github.com/JMURv/e-commerce/users/pkg/model"
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
		&model.User{},
	)
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{conn: conn}
}

func (r *Repository) GetUsersList(_ context.Context) ([]*model.User, error) {
	var u []*model.User
	if err := r.conn.Find(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) GetByID(_ context.Context, userID uint64) (*model.User, error) {
	var u model.User
	if err := r.conn.Where("ID=?", userID).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) GetByEmail(_ context.Context, email string) (*model.User, error) {
	var u model.User
	if err := r.conn.Where("Email=?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) Create(_ context.Context, u *model.User) (*model.User, error) {
	if u.Username == "" {
		return nil, repo.ErrUsernameIsRequired
	}

	if u.Email == "" {
		return nil, repo.ErrEmailIsRequired
	}

	if err := r.conn.Create(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *Repository) Update(ctx context.Context, userID uint64, newData *model.User) (*model.User, error) {
	u, err := r.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if newData.Username != "" {
		u.Username = newData.Username
	}

	if newData.Email != "" {
		u.Email = newData.Email
	}
	r.conn.Save(&u)
	return u, nil
}

func (r *Repository) Delete(_ context.Context, userID uint64) error {
	var u model.User
	if err := r.conn.Delete(&u, userID).Error; err != nil {
		return err
	}
	return nil
}
