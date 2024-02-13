package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

var itemConn *grpc.ClientConn

func init() {
	//reg, err := consul.NewRegistry("localhost:8080")
	//addrs, err := consul.ServiceAddresses(context.Background(), "items")
	var err error
	itemConn, err = grpc.Dial("localhost:50080", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
	}
}

func ListItem(w http.ResponseWriter, r *http.Request) {
	cli := pb.NewItemServiceClient(itemConn)
	items, _ := cli.ListItem(context.Background(), &pb.EmptyRequest{})

	resp, err := json.Marshal(items)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf(ErrWhileEncoding, err))
		return
	}

	utils.OkResponse(w, http.StatusOK, resp)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	NewItem := &pb.CreateItemRequest{}
	utils.ParseBody(r, NewItem)

	cli := pb.NewItemServiceClient(itemConn)
	i, err := cli.CreateItem(context.Background(), NewItem)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Error creating item: %v", err))
		return
	}

	resp, err := json.Marshal(i)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf(ErrWhileEncoding, err))
		return
	}

	utils.OkResponse(w, http.StatusCreated, resp)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse itemID: %v", err))
		return
	}

	cli := pb.NewItemServiceClient(itemConn)
	i, err := cli.GetItemByID(context.Background(), &pb.GetItemByIDRequest{ItemId: itemID})
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error getting item: %v", err))
		return
	}

	resp, err := json.Marshal(i)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf(ErrWhileEncoding, err))
		return
	}

	utils.OkResponse(w, http.StatusOK, resp)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse itemID: %v", err))
		return
	}

	var newData = &pb.UpdateItemRequest{ItemId: itemID}
	utils.ParseBody(r, newData)

	reqUserId := r.Context().Value("reqUserId").(uint64)
	if reqUserId != newData.UserId {
		utils.ErrResponse(w, http.StatusForbidden, ErrNoPermissions)
		return
	}

	cli := pb.NewItemServiceClient(itemConn)
	i, err := cli.UpdateItem(context.Background(), newData)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error updating item: %v", err))
		return
	}

	resp, err := json.Marshal(i)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf(ErrWhileEncoding, err))
		return
	}

	utils.OkResponse(w, http.StatusOK, resp)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse itemID: %v", err))
		return
	}

	reqUserId := r.Context().Value("reqUserId").(uint64)

	cli := pb.NewItemServiceClient(itemConn)
	_, err = cli.DeleteItem(context.Background(), &pb.DeleteItemRequest{ItemId: itemID, ReqUserId: reqUserId})
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting item: %v", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
