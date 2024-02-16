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

type CategoryHandler interface {
	ListCategory(w http.ResponseWriter, r *http.Request)
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateCategory(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
}

type CategoryCtrl struct {
	conn *grpc.ClientConn
}

func NewCategoryCtrl() *CategoryCtrl {
	reg, err := consul.NewRegistry("localhost:8500")
	addrs, err := reg.ServiceAddresses(context.Background(), "items")

	r := rand.Intn(len(addrs) + 1)
	conn, err := grpc.Dial(addrs[r], grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to category service: %v", err)
	}
	return &CategoryCtrl{
		conn: conn,
	}
}

func (ctrl *CategoryCtrl) ListCategory(w http.ResponseWriter, r *http.Request) {
	cli := pb.NewCategoryServiceClient(ctrl.conn)
	categories, _ := cli.GetAllCategories(context.Background(), &pb.EmptyRequest{})

	utils.OkResponse(w, http.StatusOK, categories)
}

func (ctrl *CategoryCtrl) GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	cli := pb.NewCategoryServiceClient(ctrl.conn)
	c, err := cli.GetCategoryByID(context.Background(), &pb.GetCategoryByIDRequest{CategoryId: categoryID})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	utils.OkResponse(w, http.StatusOK, c)
}

func (ctrl *CategoryCtrl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	newCategory := &pb.CreateCategoryRequest{}
	utils.ParseBody(r, newCategory)

	cli := pb.NewCategoryServiceClient(ctrl.conn)
	c, err := cli.CreateCategory(context.Background(), newCategory)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusCreated, c)
}

func (ctrl *CategoryCtrl) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	newData := &pb.UpdateCategoryRequest{CategoryId: categoryID}
	utils.ParseBody(r, newData)

	cli := pb.NewCategoryServiceClient(ctrl.conn)
	c, err := cli.UpdateCategory(context.Background(), newData)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusOK, c)
}

func (ctrl *CategoryCtrl) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	cli := pb.NewCategoryServiceClient(ctrl.conn)
	_, err = cli.DeleteCategory(context.Background(), &pb.DeleteCategoryRequest{CategoryId: categoryID})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrWhileDeletingObj)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
