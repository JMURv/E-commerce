package grpc

import (
	"context"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	metrics "github.com/JMURv/e-commerce/items/internal/metrics/prometheus"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (h *Handler) ListTags(ctx context.Context, req *pb.EmptyRequest) (*pb.ListTagsResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.ListTags.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "ListTags")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req")
	}

	t, err := h.ctrl.ListTags(ctx)
	if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.ListTagsResponse{Tags: model.TagsToProto(t)}, nil
}

func (h *Handler) CreateTag(ctx context.Context, req *pb.TagRequest) (*cm.Tag, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.CreateTag.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "CreateTag")
	}()

	if req == nil || req.Name == "" {
		statusCode = codes.InvalidArgument
		return nil, status.Error(statusCode, "nil req or empty name")
	}

	t, err := h.ctrl.CreateTag(ctx, model.TagFromProto(&cm.Tag{
		Name: req.Name,
	}))
	if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.TagToProto(t), nil
}

func (h *Handler) DeleteTag(ctx context.Context, req *pb.TagRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.DeleteTag.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteTag")
	}()

	if req == nil || req.Name == "" {
		statusCode = codes.InvalidArgument
		return nil, status.Error(statusCode, "nil req or empty name")
	}

	if err := h.ctrl.DeleteTag(ctx, req.Name); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.EmptyResponse{}, nil
}
