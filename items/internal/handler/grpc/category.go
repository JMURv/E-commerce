package grpc

import (
	"context"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) GetAllCategories(ctx context.Context, req *pb.EmptyRequest) (*pb.ListCategoriesResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	if c, err := h.ctrl.GetAllCategories(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		return &pb.ListCategoriesResponse{Categories: model.CategoriesToProto(c)}, nil
	}
}

func (h *Handler) GetCategoryByID(ctx context.Context, req *pb.GetCategoryByIDRequest) (*cm.Category, error) {
	return &cm.Category{}, nil
}

func (h *Handler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*cm.Category, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}

	c, err := h.ctrl.CreateCategory(ctx, model.CategoryFromProto(&cm.Category{
		Name:             req.Name,
		Description:      req.Description,
		ParentCategoryId: req.ParentCategoryId,
	}))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return model.CategoryToProto(c), nil
}

func (h *Handler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*cm.Category, error) {
	categoryID := req.CategoryId
	if req == nil || categoryID == 0 {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id")
	}

	c, err := h.ctrl.UpdateCategory(ctx, categoryID, model.CategoryFromProto(&cm.Category{
		Name:             req.Name,
		Description:      req.Description,
		ParentCategoryId: req.ParentCategoryId,
	}))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return model.CategoryToProto(c), nil
}

func (h *Handler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.EmptyResponse, error) {
	categoryID := req.CategoryId
	if req == nil || categoryID == 0 {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteItem(ctx, categoryID); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
