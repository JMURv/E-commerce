package controllers

import (
	"context"
	pb "github.com/JMURv/e-commerce/api/pb/favorite"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type FavoriteHandler interface {
	ListFavorites(w http.ResponseWriter, r *http.Request)
	CreateFavorites(w http.ResponseWriter, r *http.Request)
	DeleteFavorite(w http.ResponseWriter, r *http.Request)
}

type FavoriteCtrl struct {
	cli pb.FavoriteServiceClient
}

func NewFavoriteCtrl() *FavoriteCtrl {
	reg, err := consul.NewRegistry("localhost:8500")
	addrs, err := reg.ServiceAddresses(context.Background(), "favorites")

	r := rand.Intn(len(addrs))
	conn, err := grpc.Dial(addrs[r], grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to favorites service: %v", err)
	}
	return &FavoriteCtrl{
		cli: pb.NewFavoriteServiceClient(conn),
	}
}

func (ctrl *FavoriteCtrl) ListFavorites(w http.ResponseWriter, r *http.Request) {
	reqUserID := r.Context().Value("reqUserId").(uint64)

	favs, _ := ctrl.cli.GetAllUserFavorites(context.Background(), &pb.GetAllUserFavoritesRequest{UserId: reqUserID})
	utils.OkResponse(w, http.StatusOK, favs)
}

func (ctrl *FavoriteCtrl) CreateFavorites(w http.ResponseWriter, r *http.Request) {
	reqUserID := r.Context().Value("reqUserId").(uint64)
	newFavorite := &pb.CreateFavoriteRequest{UserId: reqUserID}
	utils.ParseBody(r, newFavorite)

	f, err := ctrl.cli.CreateFavorite(context.Background(), newFavorite)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.OkResponse(w, http.StatusCreated, f)

	// TODO: Create kafka message
	//newNotification := models.Notification{
	//	Type:       "notification",
	//	UserID:     favorite.UserID,
	//	ReceiverID: uint64(favorite.Item.UserID),
	//	Message:    "new favorite",
	//}
	//
	//notification, err := newNotification.CreateNotification()
	//if err != nil {
	//	log.Printf("Error while creating notification: %v", err)
	//}
	//
	//notificationBytes, err := json.Marshal(notification)
	//if err != nil {
	//	log.Printf("Error while encoding notification message: %v", err)
	//}
	//
	//go broadcast(uint(favorite.UserID), favorite.Item.UserID, notificationBytes)
}

func (ctrl *FavoriteCtrl) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	favoriteID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint64)
	_, err = ctrl.cli.DeleteFavorite(context.Background(), &pb.DeleteFavoriteIDRequest{
		FavoriteId: favoriteID,
		UserId:     reqUserID,
	})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	utils.OkResponse(w, http.StatusOK, "success")
}
