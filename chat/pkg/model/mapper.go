package model

import (
	pb "github.com/JMURv/e-commerce/api/pb/chat"
	"time"
)

func MessageToProto(msg *Message) *pb.Message {
	return &pb.Message{
		Id:        msg.ID,
		UserId:    msg.UserID,
		RoomId:    msg.RoomID,
		ReplyToId: *msg.ReplyToID,
		Text:      msg.Text,
		Seen:      msg.Seen,
		Edited:    msg.Edited,
		CreatedAt: uint64(msg.CreatedAt.Unix()),
	}
}

func MessageFromProto(msg *pb.Message) *Message {
	return &Message{
		ID:        msg.Id,
		UserID:    msg.UserId,
		RoomID:    msg.RoomId,
		Text:      msg.Text,
		Seen:      msg.Seen,
		Edited:    msg.Edited,
		ReplyToID: &msg.ReplyToId,
		CreatedAt: time.Unix(int64(msg.CreatedAt), 0),
	}
}
