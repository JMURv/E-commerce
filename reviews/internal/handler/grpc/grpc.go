package grpc

import (
	"context"
	"errors"
	"github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/review"
	controller "github.com/JMURv/e-commerce/reviews/internal/controller/review"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	pb.ReviewServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetReviewsByUserID(ctx context.Context, req *pb.GetReviewByIDRequest) (*common.Review, error) {
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

func (h *Handler) GetReviewByID(ctx context.Context, req *pb.GetReviewByIDRequest) (*common.Review, error) {
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

func (h *Handler) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*common.Review, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req")
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
	}

	newReview := &common.Review{}
	if err = proto.Unmarshal(reqData, newReview); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal request: %v", err)
	}

	return h.ctrl.CreateReview(ctx, newReview)
}

func (h *Handler) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*common.Review, error) {
	reviewID := req.ReviewId
	if req == nil || reviewID == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	reqData, err := proto.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
	}

	updateReviewData := &common.Review{}
	if err = proto.Unmarshal(reqData, updateReviewData); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal request: %v", err)
	}

	return h.ctrl.UpdateReview(ctx, reviewID, updateReviewData)
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
