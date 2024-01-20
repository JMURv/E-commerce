package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/JMURv/e-commerce/gateway/pkg/auth"
	"github.com/JMURv/e-commerce/gateway/pkg/config"
	"github.com/JMURv/e-commerce/gateway/pkg/models"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	pb "github.com/JMURv/protos/ecom/user"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

var userConn *grpc.ClientConn

func init() {
	var err error
	userConn, err = grpc.Dial(config.UserSVC, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
	}
}

func ListCreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := models.GetAllUsers()

		response, err := json.Marshal(users)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}
		utils.ResponseOk(w, http.StatusOK, response)

	case http.MethodPost:
		var userData = &pb.CreateUserRequest{}
		utils.ParseBody(r, userData)

		client := pb.NewUserServiceClient(userConn)
		u, err := client.CreateUser(context.Background(), userData)
		if err != nil {
			log.Printf("Error creating user: %v", err)
			return
		}

		token, err := auth.GenerateToken(u.User.Id)
		if err != nil {
			log.Printf("Error generating token: %v", err)
			return
		}

		response, err := json.Marshal(&pb.CreateUserResponse{
			User:  u.User,
			Token: token,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusInternalServerError)
			return
		}

		utils.ResponseOk(w, http.StatusCreated, response)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse userID: %v", err), http.StatusInternalServerError)
		return
	}

	client := pb.NewUserServiceClient(userConn)
	response, err := client.GetUserByID(context.Background(), &pb.GetUserByIDRequest{UserId: userID})
	if err != nil {
		log.Printf("Error getting user: %v", err)
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	utils.ResponseOk(w, http.StatusOK, jsonResponse)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse userID: %v", err), http.StatusInternalServerError)
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint64)
	if reqUserID != userID {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	newData := &pb.UpdateUserRequest{
		UserId: userID,
	}
	utils.ParseBody(r, newData)

	client := pb.NewUserServiceClient(userConn)
	updatedUser, err := client.UpdateUser(context.Background(), newData)
	if err != nil {
		log.Printf("Error updating review: %v", err)
		http.Error(w, fmt.Sprintf("Error updating review: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(updatedUser.User)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusInternalServerError)
		return
	}

	utils.ResponseOk(w, http.StatusOK, response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse userID: %v", err), http.StatusInternalServerError)
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint64)
	if reqUserID != userID {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	client := pb.NewUserServiceClient(userConn)
	_, err = client.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		UserId: reqUserID,
	})
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, fmt.Sprintf("Cannot delete user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
