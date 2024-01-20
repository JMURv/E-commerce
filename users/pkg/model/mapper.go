package model

import "github.com/JMURv/protos/ecom/common"

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
