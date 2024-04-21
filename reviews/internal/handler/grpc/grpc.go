package grpc

import (
	"context"
	"errors"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/review"
	controller "github.com/JMURv/e-commerce/reviews/internal/controller/review"
	metrics "github.com/JMURv/e-commerce/reviews/internal/metrics/prometheus"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Handler struct {
	pb.ReviewServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) AggregateUserRatingByID(ctx context.Context, req *pb.ByUserIDRequest) (*pb.AggregateUserRatingByIDResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("reviews.AggregateUserRatingByID.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "AggregateUserRatingByID")
	}()

	userID := req.UserId
	if req == nil || userID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	r, err := h.ctrl.AggregateUserRatingByID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		span.SetTag("error", true)
		statusCode = codes.NotFound
		return nil, status.Errorf(statusCode, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.AggregateUserRatingByIDResponse{Rating: r}, nil
}

func (h *Handler) GetReviewsByUserID(ctx context.Context, req *pb.ByUserIDRequest) (*pb.ListReviewResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("reviews.GetReviewsByUserID.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetReviewsByUserID")
	}()

	userID := req.UserId
	if req == nil || userID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	r, err := h.ctrl.GetReviewsByUserID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		span.SetTag("error", true)
		statusCode = codes.NotFound
		return nil, status.Errorf(statusCode, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.ListReviewResponse{Reviews: model.ReviewsToProto(*r)}, nil
}

func (h *Handler) GetReviewByID(ctx context.Context, req *pb.GetReviewByIDRequest) (*cm.Review, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("reviews.GetReviewByID.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetReviewByID")
	}()

	reviewID := req.ReviewId
	if req == nil || reviewID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	r, err := h.ctrl.GetReviewByID(ctx, reviewID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		span.SetTag("error", true)
		statusCode = codes.NotFound
		return nil, status.Errorf(statusCode, err.Error())
	} else if err != nil {
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.ReviewToProto(r), nil
}

func (h *Handler) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*cm.Review, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("reviews.CreateReview.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "CreateReview")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req")
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
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.ReviewToProto(r), nil
}

func (h *Handler) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*cm.Review, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("reviews.UpdateReview.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "UpdateReview")
	}()

	reviewID := req.ReviewId
	if req == nil || reviewID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
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
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.ReviewToProto(r), nil
}

func (h *Handler) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("reviews.DeleteReview.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteReview")
	}()

	if req == nil || req.ReviewId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteReview(ctx, req.ReviewId); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, err
	}

	statusCode = codes.OK
	return &pb.EmptyResponse{}, nil
}
