package controllers

import (
	"context"
	pb "github.com/JMURv/e-commerce/api/pb/review"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type ReviewHandler interface {
	GetReview(w http.ResponseWriter, r *http.Request)
	CreateReview(w http.ResponseWriter, r *http.Request)
	UpdateReview(w http.ResponseWriter, r *http.Request)
	DeleteReview(w http.ResponseWriter, r *http.Request)
}

type ReviewCtrl struct {
	cli pb.ReviewServiceClient
}

func NewReviewCtrl() *ReviewCtrl {
	reg, err := consul.NewRegistry("localhost:8500")
	addrs, err := reg.ServiceAddresses(context.Background(), "reviews")

	r := rand.Intn(len(addrs))
	conn, err := grpc.Dial(addrs[r], grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to reviews service: %v", err)
	}
	return &ReviewCtrl{
		cli: pb.NewReviewServiceClient(conn),
	}
}

func (ctrl *ReviewCtrl) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	rev, err := ctrl.cli.GetReviewByID(context.Background(), &pb.GetReviewByIDRequest{ReviewId: reviewID})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	utils.OkResponse(w, http.StatusOK, rev)
}

func (ctrl *ReviewCtrl) CreateReview(w http.ResponseWriter, r *http.Request) {
	newReview := &pb.CreateReviewRequest{}
	utils.ParseBody(r, newReview)

	rev, err := ctrl.cli.CreateReview(context.Background(), newReview)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusCreated, rev)

	// TODO: Отправить смс-ку в очередь сообщений для уведомлений
	//newNotification := models.Notification{
	//	Type:       "notification",
	//	UserID:     rev.UserId,
	//	ReceiverID: rev.ReviewedUserId,
	//	Message:    "new review",
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
	//go broadcast(uint(rev.UserId), uint(rev.ReviewedUserId), notificationBytes)
}

func (ctrl *ReviewCtrl) UpdateReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	newData := &pb.UpdateReviewRequest{ReviewId: reviewID}
	utils.ParseBody(r, newData)

	rev, err := ctrl.cli.UpdateReview(context.Background(), newData)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusOK, rev)
}

func (ctrl *ReviewCtrl) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	_, err = ctrl.cli.DeleteReview(context.Background(), &pb.DeleteReviewRequest{
		ReviewId: reviewID,
	})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
