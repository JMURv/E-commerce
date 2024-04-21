package grpc

import (
	"context"
	"errors"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	controller "github.com/JMURv/e-commerce/items/internal/controller/item"
	metrics "github.com/JMURv/e-commerce/items/internal/metrics/prometheus"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Handler struct {
	pb.ItemServiceServer
	pb.CategoryServiceServer
	pb.TagServiceServer
	ctrl *controller.Controller
}

func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) ListUserItemsByID(ctx context.Context, req *pb.ListUserItemsByIDRequest) (*pb.ListItemResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.ListUserItemsByID.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "ListUserItemsByID")
	}()

	userID := req.UserId
	if req == nil || userID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	i, err := h.ctrl.ListUserItemsByID(ctx, userID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		statusCode = codes.NotFound
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	} else if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.ListItemResponse{Items: model.ItemsToProto(i)}, nil
}

func (h *Handler) ListItem(ctx context.Context, req *pb.EmptyRequest) (*pb.ListItemResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.ListItem.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "ListItem")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req")
	}

	if i, err := h.ctrl.ListItem(ctx); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	} else {
		statusCode = codes.OK
		return &pb.ListItemResponse{Items: model.ItemsToProto(i)}, nil
	}
}

func (h *Handler) GetItemByID(ctx context.Context, req *pb.GetItemByIDRequest) (*cm.Item, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.GetItemByID.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetItemByID")
	}()

	itemID := req.ItemId
	if req == nil || itemID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	i, err := h.ctrl.GetItemByID(ctx, itemID)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		statusCode = codes.NotFound
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	} else if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.ItemToProto(i), nil
}

func (h *Handler) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*cm.Item, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.CreateItem.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "CreateItem")
	}()

	if req == nil {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req")
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
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.ItemToProto(i), nil
}

func (h *Handler) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*cm.Item, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.UpdateItem.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "UpdateItem")
	}()

	itemID := req.ItemId
	if req == nil || itemID == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
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
		span.SetTag("error", true)
		statusCode = codes.Internal
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.ItemToProto(i), nil
}

func (h *Handler) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("items.DeleteItem.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteItem")
	}()

	if req == nil || req.ItemId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Error(statusCode, "nil req or empty id")
	}

	if err := h.ctrl.DeleteItem(ctx, req.ItemId); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return &pb.EmptyResponse{}, nil
}
