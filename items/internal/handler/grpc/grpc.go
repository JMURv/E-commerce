package grpc

import (
	"context"
	"errors"
	"github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	controller "github.com/JMURv/e-commerce/items/internal/controller/item"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	pb.ItemServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) ListUserItemsByID(ctx context.Context, req *pb.ListUserItemsByIDRequest) (*pb.ListItemResponse, error) {
	userID := req.UserId
	if req == nil || userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	items, err := h.ctrl.ListUserItemsByID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListItemResponse{Items: model.ItemsToProto(*items)}, nil
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
	if err != nil && errors.Is(err, controller.ErrNotFound) {
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

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
	}

	newItem := &common.Item{}
	if err = proto.Unmarshal(reqData, newItem); err != nil {
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

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
	}

	updateItemData := &common.Item{}
	if err = proto.Unmarshal(reqData, updateItemData); err != nil {
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
