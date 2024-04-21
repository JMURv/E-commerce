package grpc

import (
	"context"
	"errors"
	"github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/user"
	controller "github.com/JMURv/e-commerce/users/internal/controller/user"
	metrics "github.com/JMURv/e-commerce/users/internal/metrics/prometheus"
	"github.com/JMURv/e-commerce/users/pkg/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"time"
)

type Handler struct {
	pb.UserServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) ListUser(ctx context.Context, req *pb.EmptyRequest) (*pb.ListUserResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("users.ListUser.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "ListUser")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req")
	}

	if u, err := h.ctrl.GetUsersList(ctx); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	} else {
		statusCode = codes.OK
		return &pb.ListUserResponse{Users: model.UsersToProto(u)}, nil
	}
}

func (h *Handler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*common.User, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("users.GetUserByID.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetUserByID")
	}()

	userID := req.UserId
	if req == nil || userID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	u, err := h.ctrl.GetUserByID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		statusCode = codes.NotFound
		return nil, status.Errorf(statusCode, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.UserToProto(u), nil
}

func (h *Handler) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*common.User, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("users.GetUserByEmail.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetUserByEmail")
	}()

	userEmail := req.Email
	if req == nil || userEmail == "" {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	u, err := h.ctrl.GetUserByEmail(ctx, userEmail)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		statusCode = codes.NotFound
		return nil, status.Errorf(statusCode, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.UserToProto(u), nil
}

func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*common.User, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("users.CreateUser.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "CreateUser")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req")
	}

	u, err := h.ctrl.CreateUser(ctx, model.UserFromProto(&common.User{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
	}))
	if err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.UserToProto(u), nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*common.User, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("users.UpdateUser.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "UpdateUser")
	}()

	userID := req.UserId
	if req == nil || userID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	updateUserData := &common.User{}
	if err = proto.Unmarshal(reqData, updateUserData); err != nil {
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	u, err := h.ctrl.UpdateUser(ctx, userID, model.UserFromProto(updateUserData))
	if err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.UserToProto(u), nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("users.DeleteUser.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteUser")
	}()

	if req == nil || req.UserId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	if err := h.ctrl.DeleteUser(ctx, req.UserId); err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.EmptyResponse{}, nil
}
