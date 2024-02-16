package grpc

import (
	"context"
	"errors"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	controller "github.com/JMURv/e-commerce/items/internal/controller/item"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.ItemServiceServer
	pb.CategoryServiceServer
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

	i, err := h.ctrl.ListUserItemsByID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListItemResponse{Items: model.ItemsToProto(i)}, nil
}

func (h *Handler) ListItem(ctx context.Context, req *pb.EmptyRequest) (*pb.ListItemResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	if i, err := h.ctrl.ListItem(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		return &pb.ListItemResponse{Items: model.ItemsToProto(i)}, nil
	}
}

func (h *Handler) GetItemByID(ctx context.Context, req *pb.GetItemByIDRequest) (*cm.Item, error) {
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

func (h *Handler) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*cm.Item, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	i, err := h.ctrl.CreateItem(ctx, model.ItemFromProto(&cm.Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		UserId:      req.UserId,
		Tags:        req.Tags,
		Quantity:    req.Quantity,
	}))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return model.ItemToProto(i), nil
}

func (h *Handler) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*cm.Item, error) {
	itemID := req.ItemId
	if req == nil || itemID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	i, err := h.ctrl.UpdateItem(ctx, itemID, model.ItemFromProto(&cm.Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryId,
		UserId:      req.UserId,
		Tags:        req.Tags,
		Quantity:    req.Quantity,
	}))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update item: %v", err)
	}
	return model.ItemToProto(i), nil
}

func (h *Handler) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.ItemId == 0 {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteItem(ctx, req.ItemId); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
