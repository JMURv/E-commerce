package model

import (
	"gorm.io/gorm"
)

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
