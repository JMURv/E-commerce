package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/user"
	"github.com/JMURv/e-commerce/gateway/pkg/auth"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
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
	userConn, err = grpc.Dial("localhost:50075", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
	}
}

func ListCreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cli := pb.NewUserServiceClient(userConn)
		u, _ := cli.ListUser(context.Background(), &pb.EmptyRequest{})

		resp, err := json.Marshal(u.Users)
		if err != nil {
			utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf(ErrWhileEncoding, err))
			return
		}
		utils.OkResponse(w, http.StatusOK, resp)

	case http.MethodPost:
		var userData = &pb.CreateUserRequest{}
		utils.ParseBody(r, userData)

		cli := pb.NewUserServiceClient(userConn)
		u, err := cli.CreateUser(context.Background(), userData)
		if err != nil {
			utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
			return
		}

		token, err := auth.GenerateToken(u.Id)
		if err != nil {
			utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf(ErrWhileGenToken, err))
			return
		}

		resp, err := json.Marshal(struct {
			User  *common.User
			Token string
		}{
			User:  u,
			Token: token,
		})
		if err != nil {
			utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf(ErrWhileEncoding, err))
			return
		}

		utils.OkResponse(w, http.StatusCreated, resp)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse userID: %v", err))
		return
	}

	cli := pb.NewUserServiceClient(userConn)
	u, err := cli.GetUserByID(context.Background(), &pb.GetUserByIDRequest{UserId: userID})
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error getting user: %v", err))
		return
	}

	resp, err := json.Marshal(u)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf(ErrWhileEncoding, err))
		return
	}

	utils.OkResponse(w, http.StatusOK, resp)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse userID: %v", err))
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint64)
	if reqUserID != userID {
		utils.ErrResponse(w, http.StatusForbidden, ErrNoPermissions)
		return
	}

	newData := &pb.UpdateUserRequest{
		UserId: userID,
	}
	utils.ParseBody(r, newData)

	cli := pb.NewUserServiceClient(userConn)
	u, err := cli.UpdateUser(context.Background(), newData)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user: %v", err))
		return
	}

	resp, err := json.Marshal(u)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf(ErrWhileEncoding, err))
		return
	}

	utils.OkResponse(w, http.StatusOK, resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot parse userID: %v", err))
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint64)
	if reqUserID != userID {
		utils.ErrResponse(w, http.StatusForbidden, ErrNoPermissions)
		return
	}

	cli := pb.NewUserServiceClient(userConn)
	_, err = cli.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		UserId: reqUserID,
	})
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting user: %v", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
