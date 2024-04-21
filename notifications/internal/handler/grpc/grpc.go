package handler

import (
	"context"
	pb "github.com/JMURv/e-commerce/api/pb/notification"
	ctrl "github.com/JMURv/e-commerce/notifications/internal/controller/notifications"
	metrics "github.com/JMURv/e-commerce/notifications/internal/metrics/prometheus"
	mdl "github.com/JMURv/e-commerce/notifications/pkg/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"sync"
	"time"
)

type BrokerRepository interface{}

type Connection struct {
	stream pb.Broadcast_CreateStreamServer
	userID uint64
	active bool
	error  chan error
}

type Pool struct {
	Connection []*Connection
}

type Handler struct {
	pb.BroadcastServer
	pb.NotificationsServer
	ctrl *ctrl.Controller
	pool *Pool
}

func New(ctrl *ctrl.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
		pool: &Pool{
			Connection: []*Connection{},
		},
	}
}

// Boradcasting
func (h *Handler) CreateStream(pbConn *pb.Connect, stream pb.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		userID: pbConn.User.Id,
		active: true,
		error:  make(chan error),
	}

	h.pool.Connection = append(h.pool.Connection, conn)
	log.Printf("UserID: %v has been connected\n", pbConn.User.Id)
	return <-conn.error
}

func (h *Handler) Broadcast(ctx context.Context, msg *mdl.Notification) error {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("notifications.Broadcast.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "Broadcast")
	}()

	var wg sync.WaitGroup
	for _, conn := range h.pool.Connection {
		wg.Add(1)
		go func(msg *mdl.Notification, conn *Connection) {
			defer wg.Done()
			if conn.active && conn.userID == msg.ReceiverID {
				log.Printf("Sending message to: %v from %v\n", conn.userID, msg.UserID)
				if err := conn.stream.Send(mdl.NotificationToProto(msg)); err != nil {
					log.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	wg.Wait()
	return nil
}

func (h *Handler) ListUserNotifications(ctx context.Context, req *pb.ByUserIDRequest) (*pb.ListNotificationResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("notifications.ListUserNotifications.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "ListUserNotifications")
	}()

	userID := req.UserId
	if req == nil || userID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	n, err := h.ctrl.ListUserNotifications(ctx, userID)
	if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.ListNotificationResponse{Notifications: mdl.NotificationsToProto(*n)}, nil
}

func (h *Handler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("notifications.DeleteNotification.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteNotification")
	}()

	if req == nil || req.Id == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	if err := h.ctrl.DeleteNotification(ctx, req.Id); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.EmptyResponse{}, nil
}

func (h *Handler) DeleteAllNotifications(ctx context.Context, req *pb.ByUserIDRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("notifications.DeleteAllNotifications.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteAllNotifications")
	}()

	if req == nil || req.UserId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	if err := h.ctrl.DeleteAllNotifications(ctx, req.UserId); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.EmptyResponse{}, nil
}
