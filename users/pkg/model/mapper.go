package model

import "github.com/JMURv/e-commerce/api/pb/common"

func UserFromProto(u *common.User) *User {
	return &User{
		ID:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		IsAdmin:  u.IsAdmin,
	}
}

func UserToProto(u *User) *common.User {
	return &common.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		IsAdmin:  u.IsAdmin,
	}
}

func UsersToProto(u []*User) []*common.User {
	var users []*common.User
	for i := range u {
		users = append(users, UserToProto(u[i]))
	}
	return users
}
