package controllers

import (
	"e-commerce/pkg/auth"
	"e-commerce/pkg/models"
	"e-commerce/pkg/utils"
	"fmt"
	"net/http"
)

type Login struct {
	Email string `json:"email"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginData := &Login{}
	utils.ParseBody(r, loginData)

	userDetails := models.GetUserByEmail(loginData.Email)
	if userDetails.Username == "" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(int(userDetails.ID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate token: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send request: %v", err), http.StatusInternalServerError)
		return
	}
}
