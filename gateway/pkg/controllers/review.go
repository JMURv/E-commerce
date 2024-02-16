package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/review"
	"github.com/JMURv/e-commerce/gateway/pkg/config"
	"github.com/JMURv/e-commerce/gateway/pkg/models"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

var reviewConn *grpc.ClientConn

func init() {
	var err error
	reviewConn, err = grpc.Dial(config.ReviewServiceURL, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to review service: %v", err)
	}
}

func GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse reviewID: %v", err), http.StatusInternalServerError)
		return
	}

	client := pb.NewReviewServiceClient(reviewConn)
	response, err := client.GetReviewByID(context.Background(), &pb.GetReviewByIDRequest{ReviewId: reviewID})
	if err != nil {
		log.Printf("Error getting review: %v", err)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	utils.OkResponse(w, http.StatusOK, jsonResponse)
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	newReviewRequest := &pb.CreateReviewRequest{}
	utils.ParseBody(r, newReviewRequest)

	client := pb.NewReviewServiceClient(reviewConn)
	review, err := client.CreateReview(context.Background(), newReviewRequest)
	if err != nil {
		log.Printf("Error creating review: %v", err)
	}

	response, err := json.Marshal(review)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	// Create notification for item's author and send message via WS
	newNotification := models.Notification{
		Type:       "notification",
		UserID:     review.UserId,
		ReceiverID: review.ReviewedUserId,
		Message:    "new review",
	}

	notification, err := newNotification.CreateNotification()
	if err != nil {
		log.Printf("Error while creating notification: %v", err)
	}

	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error while encoding notification message: %v", err)
	}

	go broadcast(uint(review.UserId), uint(review.ReviewedUserId), notificationBytes)

	utils.OkResponse(w, http.StatusCreated, response)
}

func UpdateReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse reviewID: %v", err), http.StatusInternalServerError)
		return
	}

	newData := &pb.UpdateReviewRequest{
		ReviewId: reviewID,
	}
	utils.ParseBody(r, newData)

	client := pb.NewReviewServiceClient(reviewConn)
	updatedReview, err := client.UpdateReview(context.Background(), newData)
	if err != nil {
		log.Printf("Error updating review: %v", err)
		http.Error(w, fmt.Sprintf("Error updating review: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(updatedReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	utils.OkResponse(w, http.StatusOK, response)
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse reviewID: %v", err), http.StatusInternalServerError)
		return
	}

	client := pb.NewReviewServiceClient(reviewConn)
	_, err = client.DeleteReview(context.Background(), &pb.DeleteReviewRequest{
		ReviewId: reviewID,
	})
	if err != nil {
		log.Printf("Error deleting review: %v", err)
		http.Error(w, fmt.Sprintf("Cannot delete review: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
