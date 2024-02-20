package grpc

import (
	"context"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/favorite"
	ctrl "github.com/JMURv/e-commerce/favorites/internal/controller/favorite"
	"github.com/JMURv/e-commerce/favorites/pkg/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.FavoriteServiceServer
	ctrl *ctrl.Controller
}

func New(ctrl *ctrl.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetAllUserFavorites(ctx context.Context, req *pb.GetAllUserFavoritesRequest) (*pb.ListFavoritesResponse, error) {
	if req == nil || req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if favs, err := h.ctrl.GetAllUserFavorites(ctx, req.UserId); err != nil {
		return nil, err
	} else {
		return &pb.ListFavoritesResponse{Favorites: model.FavoritesToProto(favs)}, nil

	}
}

func (h *Handler) GetFavoriteByID(ctx context.Context, req *pb.FavoriteIDRequest) (*cm.Favorite, error) {
	if req == nil || req.FavoriteId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	f, err := h.ctrl.GetFavoriteByID(ctx, req.FavoriteId)
	if err != nil {
		return nil, err
	}

	return model.FavoriteToProto(f), nil
}

func (h *Handler) CreateFavorite(ctx context.Context, req *pb.CreateFavoriteRequest) (*cm.Favorite, error) {
	if req == nil || req.UserId == 0 || req.ItemId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty ids")
	}

	f, err := h.ctrl.CreateFavorite(ctx, &model.Favorite{
		UserID: req.UserId,
		ItemID: req.ItemId,
	})
	if err != nil {
		return nil, err
	}

	return model.FavoriteToProto(f), nil
}

func (h *Handler) DeleteFavorite(ctx context.Context, req *pb.DeleteFavoriteIDRequest) (*pb.EmptyResponse, error) {
	if req == nil || req.FavoriteId == 0 || req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}

	if err := h.ctrl.DeleteFavorite(ctx, req.FavoriteId); err != nil {
		return nil, err
	} else {
		return &pb.EmptyResponse{}, nil
	}
}
