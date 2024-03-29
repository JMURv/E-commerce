package grpc

import (
	"context"
	"errors"
	"github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/user"
	controller "github.com/JMURv/e-commerce/users/internal/controller/user"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	pb.UserServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) ListUser(ctx context.Context, req *pb.EmptyRequest) (*pb.ListUserResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	if u, err := h.ctrl.GetUsersList(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		return &pb.ListUserResponse{Users: model.UsersToProto(*u)}, nil
	}
}

func (h *Handler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*common.User, error) {
	userID := req.UserId
	if req == nil || userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	u, err := h.ctrl.GetUserByID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.UserToProto(u), nil
}

func (h *Handler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*common.User, error) {
	userEmail := req.Email
	if req == nil || userEmail == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	u, err := h.ctrl.GetUserByEmail(ctx, userEmail)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.UserToProto(u), nil
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*common.User, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	u, err := h.ctrl.CreateUser(ctx, model.UserFromProto(&common.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.UserToProto(u), nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*common.User, error) {
	userID := req.UserId
	if req == nil || userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	updateUserData := &common.User{}
	if err = proto.Unmarshal(reqData, updateUserData); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	u, err := h.ctrl.UpdateUser(ctx, userID, model.UserFromProto(updateUserData))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.UserToProto(u), nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteUser(ctx, req.UserId); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.EmptyResponse{}, nil
}
