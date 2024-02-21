package grcp

import (
	"context"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/chat"
	ctrl "github.com/JMURv/e-commerce/chat/internal/controller/chat"
	mdl "github.com/JMURv/e-commerce/chat/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type Connection struct {
	stream pb.Broadcast_CreateStreamServer
	id     uint64
	active bool
	error  chan error
}

type Pool struct {
	Connection []*Connection
}

type Handler struct {
	pb.MessagesServer
	pb.RoomsServer
	pb.BroadcastServer
	ctrl *ctrl.Controller
	pool *Pool
}

func New(ctrl *ctrl.Controller) *Handler {
	// Create a new connection pool
	return &Handler{
		ctrl: ctrl,
		pool: &Pool{
			Connection: []*Connection{},
		},
	}
}

// Boradcasting
func (h *Handler) CreateStream(pconn *pb.Connect, stream pb.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	h.pool.Connection = append(h.pool.Connection, conn)
	return <-conn.error
}

func (h *Handler) BroadcastMessage(ctx context.Context, msg *pb.Message) (*pb.Close, error) {
	wg := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range h.pool.Connection {
		wg.Add(1)

		go func(msg *pb.Message, conn *Connection) {
			defer wg.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				fmt.Printf("Sending message to: %v from %v\n", conn.id, msg.Id)
				if err != nil {
					fmt.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done
	return &pb.Close{}, nil
}

// Messages
func (h *Handler) GetMessageByID(ctx context.Context, req *pb.GetMessageByIDRequest) (*pb.Message, error) {
	if req == nil || req.MessageId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	msg, err := h.ctrl.GetMessageByID(ctx, req.MessageId)
	if err != nil {
		return nil, err
	}

	return mdl.MessageToProto(msg), nil
}

func (h *Handler) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest) (*pb.Message, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	msg, err := h.ctrl.CreateMessage(ctx, &mdl.Message{
		UserID:    req.UserId,
		RoomID:    req.RoomId,
		ReplyToID: &req.ReplyToId,
		Text:      req.Text,
	})
	if err != nil {
		return nil, err
	}

	return mdl.MessageToProto(msg), nil
}

func (h *Handler) UpdateMessage(ctx context.Context, req *pb.UpdateMessageRequest) (*pb.Message, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	msg, err := h.ctrl.UpdateMessage(ctx, req.Id, &mdl.Message{
		UserID:    req.UserId,
		RoomID:    req.RoomId,
		ReplyToID: &req.ReplyToId,
		Text:      req.Text,
	})
	if err != nil {
		return nil, err
	}

	return mdl.MessageToProto(msg), nil
}

func (h *Handler) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.MessageId == 0 || req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteMessage(ctx, req.MessageId); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}

//func (h *Handler) ListMessages(ctx context.Context, req *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
//	return nil, nil
//}
