package handler

import (
	"context"
	pb "github.com/JMURv/e-commerce/api/pb/notification"
	ctrl "github.com/JMURv/e-commerce/notifications/internal/controller/notifications"
	mdl "github.com/JMURv/e-commerce/notifications/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.NotificationsServer
	ctrl *ctrl.Controller
}

func New(ctrl *ctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) ListUserNotifications(ctx context.Context, req *pb.ByUserIDRequest) (*pb.ListNotificationResponse, error) {
	userID := req.UserId
	if req == nil || userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	n, err := h.ctrl.ListUserNotifications(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListNotificationResponse{Notifications: mdl.NotificationsToProto(*n)}, nil
}

func (h *Handler) CreateNotification(ctx context.Context, req *pb.Notification) (*pb.Notification, error) {
	if req == nil || req.Type == "" || req.UserId == 0 || req.ReceiverId == 0 || req.Message == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	n, err := h.ctrl.CreateNotification(ctx, &mdl.Notification{
		Type:       req.Type,
		UserID:     req.UserId,
		ReceiverID: req.ReceiverId,
		Message:    req.Message,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return mdl.NotificationToProto(n), nil
}

func (h *Handler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteNotification(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

func (h *Handler) DeleteAllNotifications(ctx context.Context, req *pb.ByUserIDRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteAllNotifications(ctx, req.UserId); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
