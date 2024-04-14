package model

import (
	pb "github.com/JMURv/e-commerce/api/pb/chat"
	"time"
)

func MessageToProto(msg *Message) *pb.Message {
	pbMedia := make([]*pb.Media, 0, len(msg.Media))
	for _, v := range msg.Media {
		pbMedia = append(pbMedia, &pb.Media{Id: v.ID, Url: v.FilePath})
	}

	return &pb.Message{
		Id:        msg.ID,
		UserId:    msg.UserID,
		RoomId:    msg.RoomID,
		ReplyToId: *msg.ReplyToID,
		Text:      msg.Text,
		Seen:      msg.Seen,
		Edited:    msg.Edited,
		CreatedAt: uint64(msg.CreatedAt.Unix()),
		Media:     pbMedia,
	}
}

func MessageFromProto(msg *pb.Message) *Message {
	mdMedia := make([]*Media, 0, len(msg.Media))
	for _, v := range msg.Media {
		mdMedia = append(mdMedia, &Media{ID: v.Id, FilePath: v.Url})
	}

	return &Message{
		ID:        msg.Id,
		UserID:    msg.UserId,
		RoomID:    msg.RoomId,
		Text:      msg.Text,
		Seen:      msg.Seen,
		Edited:    msg.Edited,
		ReplyToID: &msg.ReplyToId,
		CreatedAt: time.Unix(int64(msg.CreatedAt), 0),
		Media:     mdMedia,
	}
}

func MessagesToProto(msgs []*Message) []*pb.Message {
	r := make([]*pb.Message, 0, len(msgs))
	for _, msg := range msgs {
		r = append(r, MessageToProto(msg))
	}
	return r
}

func MessagesFromProto(msgs []*pb.Message) []*Message {
	r := make([]*Message, 0, len(msgs))
	for _, msg := range msgs {
		r = append(r, MessageFromProto(msg))
	}
	return r
}

func RoomToProto(room *Room) *pb.Room {
	return &pb.Room{
		Id:        room.ID,
		SellerId:  room.SellerID,
		BuyerId:   room.BuyerID,
		ItemId:    room.ItemID,
		Messages:  MessagesToProto(room.Messages),
		CreatedAt: uint64(room.CreatedAt.Unix()),
	}
}

func RoomFromProto(room *pb.Room) *Room {
	return &Room{
		ID:        room.Id,
		SellerID:  room.SellerId,
		BuyerID:   room.BuyerId,
		ItemID:    room.ItemId,
		Messages:  MessagesFromProto(room.Messages),
		CreatedAt: time.Unix(int64(room.CreatedAt), 0),
	}
}

func RoomsToProto(rooms []*Room) []*pb.Room {
	r := make([]*pb.Room, 0, len(rooms))
	for _, room := range rooms {
		r = append(r, RoomToProto(room))
	}
	return r
}
