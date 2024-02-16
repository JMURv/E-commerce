package controllers

import (
	"github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/user"
	"github.com/JMURv/e-commerce/gateway/pkg/auth"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type UserTokenResponse struct {
	User  *common.User `json:"user"`
	Token string       `json:"token"`
}

type UserHandler interface {
	ListCreateUser(w http.ResponseWriter, r *http.Request)
	GetUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserCtrl struct {
	cli pb.UserServiceClient
}

func NewUserCtrl() *UserCtrl {
	reg, err := consul.NewRegistry("localhost:8500")
	addrs, err := reg.ServiceAddresses(context.Background(), "users")

	r := rand.Intn(len(addrs) + 1)
	conn, err := grpc.Dial(addrs[r], grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
	}
	return &UserCtrl{
		cli: pb.NewUserServiceClient(conn),
	}
}

func (ctrl *UserCtrl) ListCreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u, _ := ctrl.cli.ListUser(context.Background(), &pb.EmptyRequest{})
		utils.OkResponse(w, http.StatusOK, u.Users)

	case http.MethodPost:
		var userData = &pb.CreateUserRequest{}
		utils.ParseBody(r, userData)

		u, err := ctrl.cli.CreateUser(context.Background(), userData)
		if err != nil {
			log.Println(err.Error())
			utils.ErrResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		token, err := auth.GenerateToken(u.Id)
		if err != nil {
			log.Println(err.Error())
			utils.ErrResponse(w, http.StatusInternalServerError, ErrWhileGenToken)
			return
		}

		utils.OkResponse(w, http.StatusCreated, &UserTokenResponse{
			User:  u,
			Token: token,
		})
	}
}

func (ctrl *UserCtrl) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	u, err := ctrl.cli.GetUserByID(context.Background(), &pb.GetUserByIDRequest{UserId: userID})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	utils.OkResponse(w, http.StatusOK, u)
}

func (ctrl *UserCtrl) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
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

	u, err := ctrl.cli.UpdateUser(context.Background(), newData)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusOK, u)
}

func (ctrl *UserCtrl) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusBadRequest, ErrParseParam)
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint64)
	if reqUserID != userID {
		utils.ErrResponse(w, http.StatusForbidden, ErrNoPermissions)
		return
	}

	_, err = ctrl.cli.DeleteUser(context.Background(), &pb.DeleteUserRequest{
		UserId: reqUserID,
	})
	if err != nil {
		log.Println(err.Error())
		utils.ErrResponse(w, http.StatusNotFound, ErrNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
