package grpc

import (
	"context"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) ListTags(ctx context.Context, req *pb.EmptyRequest) (*pb.ListTagsResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	t, err := h.ctrl.ListTags(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListTagsResponse{Tags: model.TagsToProto(t)}, nil
}

func (h *Handler) CreateTag(ctx context.Context, req *pb.TagRequest) (*cm.Tag, error) {
	if req == nil || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty name")
	}

	t, err := h.ctrl.CreateTag(ctx, model.TagFromProto(&cm.Tag{
		Name: req.Name,
	}))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return model.TagToProto(t), nil
}

func (h *Handler) DeleteTag(ctx context.Context, req *pb.TagRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "nil req or empty name")
	}

	if err := h.ctrl.DeleteTag(ctx, req.Name); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
