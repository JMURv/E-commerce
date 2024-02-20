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
)

type TagHandler interface {
	ListTags(w http.ResponseWriter, r *http.Request)
	CreateTag(w http.ResponseWriter, r *http.Request)
	DeleteTag(w http.ResponseWriter, r *http.Request)
}

type TagCtrl struct {
	cli pb.TagServiceClient
}

func NewTagCtrl() *TagCtrl {
	reg, err := consul.NewRegistry("localhost:8500")
	addrs, err := reg.ServiceAddresses(context.Background(), "items")

	r := rand.Intn(len(addrs))
	conn, err := grpc.Dial(addrs[r], grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to items service: %v", err)
	}
	return &TagCtrl{
		cli: pb.NewTagServiceClient(conn),
	}
}

func (ctrl *TagCtrl) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, _ := ctrl.cli.ListTags(context.Background(), &pb.EmptyRequest{})

	utils.OkResponse(w, http.StatusOK, tags)
}

func (ctrl *TagCtrl) CreateTag(w http.ResponseWriter, r *http.Request) {
	newTag := &pb.TagRequest{}
	utils.ParseBody(r, newTag)

	tag, err := ctrl.cli.CreateTag(context.Background(), newTag)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusCreated, tag)
}

func (ctrl *TagCtrl) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagName := mux.Vars(r)["name"]
	if tagName == "" {
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	_, err := ctrl.cli.DeleteTag(context.Background(), &pb.TagRequest{
		Name: tagName,
	})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
