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

func (h *Handler) GetAllCategories(ctx context.Context, req *pb.EmptyRequest) (*pb.ListCategoriesResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.GetAllCategories.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetAllCategories")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req")
	}

	if c, err := h.ctrl.GetAllCategories(ctx); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	} else {
		statusCode = codes.OK
		return &pb.ListCategoriesResponse{Categories: model.CategoriesToProto(c)}, nil
	}
}

func (h *Handler) GetCategoryByID(ctx context.Context, req *pb.GetCategoryByIDRequest) (*cm.Category, error) {
	return &cm.Category{}, nil
}

func (h *Handler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*cm.Category, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.CreateCategory.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "CreateCategory")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Error(statusCode, "nil req")
	}

	c, err := h.ctrl.CreateCategory(ctx, model.CategoryFromProto(&cm.Category{
		Name:             req.Name,
		Description:      req.Description,
		ParentCategoryId: req.ParentCategoryId,
	}))
	if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Error(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.CategoryToProto(c), nil
}

func (h *Handler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*cm.Category, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.UpdateCategory.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "UpdateCategory")
	}()

	categoryID := req.CategoryId
	if req == nil || categoryID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Error(statusCode, "nil req or empty id")
	}

	c, err := h.ctrl.UpdateCategory(ctx, categoryID, model.CategoryFromProto(&cm.Category{
		Name:             req.Name,
		Description:      req.Description,
		ParentCategoryId: req.ParentCategoryId,
	}))
	if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Error(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.CategoryToProto(c), nil
}

func (h *Handler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.DeleteCategory.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteCategory")
	}()

	categoryID := req.CategoryId
	if req == nil || categoryID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Error(statusCode, "nil req or empty id")
	}

	if err := h.ctrl.DeleteItem(ctx, categoryID); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Error(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.EmptyResponse{}, nil
}
