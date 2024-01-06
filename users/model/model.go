package model

import (
	"e-commerce/pkg/auth"
	"errors"
	"fmt"
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

	dsn := os.Getenv("DSN")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db = database

	err = db.AutoMigrate(
		&User{},
	)
	if err != nil {
		log.Fatal(err)
	}

}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type User struct {
	gorm.Model
	ID       uint64 `gorm:"primaryKey"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"isAdmin"`
}

func GetUserByID(id uint64) (*User, error) {
	var user User
	if err := db.Where("ID=?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := db.Where("Email=?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() []User {
	var Users []User
	db.Find(&Users)
	return Users
}

func (u *User) CreateUser() (*User, string, error) {

	if u.Username == "" {
		return u, "", errors.New("username is required")
	}

	if u.Email == "" {
		return u, "", errors.New("email is required")
	}

	if err := db.Create(&u).Error; err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return u, token, nil
}

func UpdateUser(userId uint64, newData *User) (*User, error) {
	user, err := GetUserByID(userId)
	if err != nil {
		return user, err
	}

	if newData.Username != "" {
		user.Username = newData.Username
	}

	if newData.Email != "" {
		user.Email = newData.Email
	}
	db.Save(&user)
	return user, nil
}

func DeleteUser(id uint64) error {
	var user User
	if err := db.Delete(&user, id).Error; err != nil {
		return err
	}
	return nil
}
