package controllers

import (
	"context"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type ItemHandler interface {
	ListItem(w http.ResponseWriter, r *http.Request)
	CreateItem(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
}

type ItemCtrl struct {
	cli pb.ItemServiceClient
}

func NewItemCtrl() *ItemCtrl {
	reg, err := consul.NewRegistry("localhost:8500")
	addrs, err := reg.ServiceAddresses(context.Background(), "items")

	r := rand.Intn(len(addrs))
	conn, err := grpc.Dial(addrs[r], grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to item service: %v", err)
	}
	return &ItemCtrl{
		cli: pb.NewItemServiceClient(conn),
	}
}

func (ctrl *ItemCtrl) ListItem(w http.ResponseWriter, r *http.Request) {
	items, _ := ctrl.cli.ListItem(context.Background(), &pb.EmptyRequest{})

	utils.OkResponse(w, http.StatusOK, items)
}

func (ctrl *ItemCtrl) CreateItem(w http.ResponseWriter, r *http.Request) {
	NewItem := &pb.CreateItemRequest{}
	utils.ParseBody(r, NewItem)

	i, err := ctrl.cli.CreateItem(context.Background(), NewItem)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusCreated, i)
}

func (ctrl *ItemCtrl) GetItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	i, err := ctrl.cli.GetItemByID(context.Background(), &pb.GetItemByIDRequest{ItemId: itemID})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	utils.OkResponse(w, http.StatusOK, i)
}

func (ctrl *ItemCtrl) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	var newData = &pb.UpdateItemRequest{ItemId: itemID}
	utils.ParseBody(r, newData)

	reqUserId := r.Context().Value("reqUserId").(uint64)
	if reqUserId != newData.UserId {
		utils.ErrResponse(w, http.StatusForbidden, ErrNoPermissions)
		return
	}

	i, err := ctrl.cli.UpdateItem(context.Background(), newData)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusOK, i)
}

func (ctrl *ItemCtrl) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	reqUserId := r.Context().Value("reqUserId").(uint64)

	_, err = ctrl.cli.DeleteItem(context.Background(), &pb.DeleteItemRequest{ItemId: itemID, ReqUserId: reqUserId})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
