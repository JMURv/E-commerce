package model

import "github.com/JMURv/e-commerce/api/pb/common"

func UserFromProto(r *common.User) *User {
	return &User{
		ID:       r.Id,
		Username: r.Username,
		Email:    r.Email,
		IsAdmin:  r.IsAdmin,
	}
}

func UserToProto(r *User) *common.User {
	return &common.User{
		Id:       r.ID,
		Username: r.Username,
		Email:    r.Email,
		IsAdmin:  r.IsAdmin,
	}
}

func UsersToProto(u []User) []*common.User {
	var users []*common.User

	for i := range u {
		users = append(users, UserToProto(&u[i]))
	}

	return users
}
