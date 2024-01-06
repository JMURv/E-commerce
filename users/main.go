package main

import (
	"e-commerce/users/model"
	common "github.com/JMURv/protos/ecom/common"
	pb "github.com/JMURv/protos/ecom/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type userServer struct {
	pb.UserServiceServer
}

//func (s *userServer) GetUsersList(ctx context.Context, req *pb.EmptyRequest) (*pb.ListUserResponse, error) {
//	users := model.GetAllUsers()
//	return users, nil
//}

func (s *userServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*common.User, error) {
	userID := req.GetUserId()

	u, err := model.GetUserByID(userID)
	if err != nil {
		log.Printf("Error getting user by ID: %v", err)
		return &common.User{}, err
	}

	response := &common.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		IsAdmin:  u.IsAdmin,
	}

	return response, nil
}

func (s *userServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*common.User, error) {
	userEmail := req.GetEmail()

	u, err := model.GetUserByEmail(userEmail)
	if err != nil {
		log.Printf("Error getting user by email: %v", err)
		return &common.User{}, err
	}

	response := &common.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		IsAdmin:  u.IsAdmin,
	}

	return response, nil
}

func (s *userServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	newUser := &model.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	u, token, err := newUser.CreateUser()
	if err != nil {
		return nil, err
	}

	response := &pb.CreateUserResponse{
		User: &common.User{
			Id:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			IsAdmin:  u.IsAdmin,
		},
		Token:        token,
		ErrorMessage: "",
	}

	return response, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*common.User, error) {
	userID := req.GetUserId()

	newData := &model.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}

	u, err := model.UpdateUser(userID, newData)
	if err != nil {
		return nil, err
	}

	response := &common.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}

	return response, nil
}

func (s *userServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.EmptyResponse, error) {
	userID := req.GetUserId()

	err := model.DeleteUser(userID)
	if err != nil {
		return nil, err
	}

	return &pb.EmptyResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50075")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &userServer{})

	reflection.Register(server)

	log.Println("User service is listening")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
