package models

import (
	"e-commerce/pkg/auth"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

type User struct {
	gorm.Model
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	IsAdmin   bool       `json:"isAdmin"`
	Items     []Item     `json:"items"`
	Favorites []Favorite `json:"favorites"`
	Reviews   []Review   `json:"reviews"`
}

func GetUserByID(id uint) (*User, error) {
	var user User
	if err := db.Preload("Items").Preload("Reviews").Preload("Favorites").Where("ID=?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) *User {
	var user User
	db.Where("Email=?", email).First(&user)
	return &user
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

func UpdateUser(userId uint, newData *User) (*User, error) {
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

func DeleteUser(id uint) error {
	var user User
	if err := db.Delete(&user, id).Error; err != nil {
		return err
	}
	return nil
}
