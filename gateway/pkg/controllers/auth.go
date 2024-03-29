package controllers

import (
	"context"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/user"
	"github.com/JMURv/e-commerce/gateway/pkg/auth"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"net/http"
)

type Login struct {
	Email string `json:"email"`
}

// LoginHandler handles the login request and generates a token for the user.
func (ctrl *UserCtrl) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginData := &Login{}
	utils.ParseBody(r, loginData)

	u, err := ctrl.cli.GetUserByEmail(context.Background(), &pb.GetUserByEmailRequest{Email: loginData.Email})
	if err != nil || u.Username == "" {
		utils.ErrResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := auth.GenerateToken(u.Id)
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, ErrWhileGenToken)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
}
