package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetUserByID(id string) *User {
	var user User
	db.Where("ID=?", id).First(&user)
	return &user
}

func GetAllUsers() []User {
	var Users []User
	db.Find(&Users)
	return Users
}

func (u *User) CreateUser() (*User, error) {
	if u.Username == "" {
		return u, errors.New("username is required")
	}
	if u.Email == "" {
		return u, errors.New("email is required")
	}
	result := db.Create(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func UpdateUser(userId string, newData *User) (*User, error) {
	user := GetUserByID(userId)
	if newData.Username != "" {
		user.Username = newData.Username
	}

	if newData.Email != "" {
		user.Email = newData.Email
	}
	db.Save(&user)
	return user, nil
}

func DeleteUser(id string) User {
	var user User
	db.Delete(&user, id)
	return user
}
