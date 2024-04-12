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

	currRoom := h.ctrl.GetRoomByID()
	newMsg, err := h.ctrl.CreateMessage(ctx, &mdl.Message{
		UserID:    msg.UserId,
		RoomID:    msg.RoomId,
		ReplyToID: &msg.ReplyToId,
		Text:      msg.Text,
	})

	for _, conn := range h.pool.Connection {
		wg.Add(1)
		go func(msg *pb.Message, conn *Connection) {
			defer wg.Done()
			if conn.id == msg.UserId {

			}
			if conn.active {
				fmt.Printf("Sending message to: %v from %v\n", conn.id, msg.UserId)
				err = conn.stream.Send(mdl.MessageToProto(newMsg))
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

// Rooms
func (h *Handler) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.Room, error) {
	if req == nil || req.SellerId == 0 || req.BuyerId == 0 || req.ItemId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	room, err := h.ctrl.CreateRoom(ctx, &mdl.Room{
		SellerID: req.SellerId,
		BuyerID:  req.BuyerId,
		ItemID:   req.ItemId,
	})
	if err != nil {
		return nil, err
	}

	return mdl.RoomToProto(room), nil
}

func (h *Handler) GetUserRooms(ctx context.Context, req *pb.ListRoomRequest) (*pb.ListRoomResponse, error) {
	if req == nil || req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	rooms, err := h.ctrl.GetUserRooms(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.ListRoomResponse{Rooms: mdl.RoomsToProto(rooms)}, nil
}

func (h *Handler) DeleteRoom(ctx context.Context, req *pb.DeleteRoomRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.RoomId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteRoom(ctx, req.RoomId); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
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
