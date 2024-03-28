package grpc

import (
	"context"
	"errors"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/review"
	controller "github.com/JMURv/e-commerce/reviews/internal/controller/review"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.ReviewServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) AggregateUserRatingByID(ctx context.Context, req *pb.ByUserIDRequest) (*pb.AggregateUserRatingByIDResponse, error) {
	userID := req.UserId
	if req == nil || userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	r, err := h.ctrl.AggregateUserRatingByID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.AggregateUserRatingByIDResponse{Rating: r}, nil
}

func (h *Handler) GetReviewsByUserID(ctx context.Context, req *pb.ByUserIDRequest) (*pb.ListReviewResponse, error) {
	userID := req.UserId
	if req == nil || userID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	r, err := h.ctrl.GetReviewsByUserID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.ListReviewResponse{Reviews: model.ReviewsToProto(*r)}, nil
}

func (h *Handler) GetReviewByID(ctx context.Context, req *pb.GetReviewByIDRequest) (*cm.Review, error) {
	reviewID := req.ReviewId
	if req == nil || reviewID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	r, err := h.ctrl.GetReviewByID(ctx, reviewID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.ReviewToProto(r), nil
}

func (h *Handler) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*cm.Review, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	r, err := h.ctrl.CreateReview(ctx, model.ReviewFromProto(&cm.Review{
		UserId:         req.UserId,
		ItemId:         req.ItemId,
		ReviewedUserId: req.ReviewedUserId,
		Advantages:     req.Advantages,
		Disadvantages:  req.Disadvantages,
		ReviewText:     req.ReviewText,
		Rating:         req.Rating,
	}))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return model.ReviewToProto(r), nil
}

func (h *Handler) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*cm.Review, error) {
	reviewID := req.ReviewId
	if req == nil || reviewID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	r, err := h.ctrl.UpdateReview(ctx, req.ReviewId, model.ReviewFromProto(&cm.Review{
		ReviewId:       req.ReviewId,
		UserId:         req.UserId,
		ItemId:         req.ItemId,
		ReviewedUserId: req.ReviewedUserId,
		Advantages:     req.Advantages,
		Disadvantages:  req.Disadvantages,
		ReviewText:     req.ReviewText,
		Rating:         req.Rating,
	}))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return model.ReviewToProto(r), nil
}

func (h *Handler) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.ReviewId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteReview(ctx, req.ReviewId); err != nil {
		return nil, err
	}
	return &pb.EmptyResponse{}, nil
}
