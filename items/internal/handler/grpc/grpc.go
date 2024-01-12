package grpc

import (
	"context"
	"errors"
	"github.com/JMURv/e-commerce/items/internal/controller/item"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/JMURv/protos/ecom/common"
	pb "github.com/JMURv/protos/ecom/item"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	pb.ItemServiceServer
	ctrl *item.Controller
}

func New(ctrl *item.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) ListItem(ctx context.Context, req *pb.EmptyRequest) (*pb.ListItemResponse, error) {
	return &pb.ListItemResponse{}, nil
}

func (h *Handler) GetItemByID(ctx context.Context, req *pb.GetItemByIDRequest) (*common.Item, error) {
	itemID := req.ItemId
	if req == nil || itemID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	i, err := h.ctrl.GetItemByID(ctx, itemID)
	if err != nil && errors.Is(err, item.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.ItemToProto(i), nil
}

func (h *Handler) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*common.Item, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	requestData, err := proto.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
	}

	newItem := &common.Item{}
	if err = proto.Unmarshal(requestData, newItem); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal request: %v", err)
	}

	i, err := h.ctrl.CreateItem(ctx, newItem)
	return i, err
}

func (h *Handler) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*common.Item, error) {
	itemID := req.ItemId
	if req == nil || itemID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	requestData, err := proto.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
	}

	updateItemData := &common.Item{}
	if err = proto.Unmarshal(requestData, updateItemData); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal request: %v", err)
	}

	i, err := h.ctrl.UpdateItem(ctx, itemID, updateItemData)
	return i, nil
}

func (h *Handler) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.ItemId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteItem(ctx, req.ItemId); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
