package grpc

import (
	"context"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/favorite"
	ctrl "github.com/JMURv/e-commerce/favorites/internal/controller/favorite"
	metrics "github.com/JMURv/e-commerce/favorites/internal/metrics/prometheus"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Handler struct {
	pb.FavoriteServiceServer
	ctrl *ctrl.Controller
}

func New(ctrl *ctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetAllUserFavorites(ctx context.Context, req *pb.GetAllUserFavoritesRequest) (*pb.ListFavoritesResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("favorites.GetAllUserFavorites.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetAllUserFavorites")
	}()

	if req == nil || req.UserId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	if favs, err := h.ctrl.GetAllUserFavorites(ctx, req.UserId); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	} else {
		statusCode = codes.OK
		return &pb.ListFavoritesResponse{Favorites: model.FavoritesToProto(favs)}, nil

	}
}

func (h *Handler) GetFavoriteByID(ctx context.Context, req *pb.FavoriteIDRequest) (*cm.Favorite, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("favorites.GetFavoriteByID.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "GetFavoriteByID")
	}()

	if req == nil || req.FavoriteId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	f, err := h.ctrl.GetFavoriteByID(ctx, req.FavoriteId)
	if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.FavoriteToProto(f), nil
}

func (h *Handler) CreateFavorite(ctx context.Context, req *pb.CreateFavoriteRequest) (*cm.Favorite, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("favorites.CreateFavorite.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "CreateFavorite")
	}()

	if req == nil || req.UserId == 0 || req.ItemId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty ids")
	}

	f, err := h.ctrl.CreateFavorite(ctx, &model.Favorite{
		UserID: req.UserId,
		ItemID: req.ItemId,
	})
	if err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	}

	statusCode = codes.OK
	return model.FavoriteToProto(f), nil
}

func (h *Handler) DeleteFavorite(ctx context.Context, req *pb.DeleteFavoriteIDRequest) (*pb.EmptyResponse, error) {
	var statusCode codes.Code
	start := time.Now()

	span := opentracing.GlobalTracer().StartSpan("favorites.DeleteFavorite.handler")
	ctx = opentracing.ContextWithSpan(ctx, span)
	defer func() {
		span.Finish()
		metrics.ObserveRequest(time.Since(start), int(statusCode), "DeleteFavorite")
	}()

	if req == nil || req.FavoriteId == 0 || req.UserId == 0 {
		statusCode = codes.InvalidArgument
		return nil, status.Errorf(statusCode, "nil req or empty id")
	}

	if err := h.ctrl.DeleteFavorite(ctx, req.FavoriteId); err != nil {
		statusCode = codes.Internal
		span.SetTag("error", true)
		return nil, status.Errorf(statusCode, err.Error())
	} else {
		statusCode = codes.OK
		return &pb.EmptyResponse{}, nil
	}
}
