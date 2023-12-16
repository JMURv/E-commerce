package controllers

import (
	"e-commerce/pkg/models"
	"e-commerce/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ListCreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users := models.GetAllUsers()

		response, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case http.MethodPost:
		var userData = &models.User{}
		utils.ParseBody(r, userData)

		user, accessToken, err := userData.CreateUser()
		if err != nil {
			http.Error(w, fmt.Sprintf("Creating user error: %v", err), http.StatusBadRequest)
			return
		}

		userWithToken := struct {
			User  *models.User `json:"user"`
			Token string       `json:"token"`
		}{
			User:  user,
			Token: accessToken,
		}

		response, err := json.Marshal(userWithToken)
		if err != nil {
			http.Error(w, "Encoding error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	userDetails, err := models.GetUserByID(userId)
	if err != nil {
		http.Error(w, "Cannot get user error", http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(userDetails)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	reqUserId := r.Context().Value("reqUserId")
	if reqUserId != userId {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	var updatedUser = &models.User{}
	utils.ParseBody(r, updatedUser)

	newUserData, err := models.UpdateUser(userId, updatedUser)
	if err != nil {
		log.Printf("[ERROR] Updating user error: %v\n", err)
		http.Error(w, fmt.Sprintf("Updating user error: %v", err), http.StatusBadRequest)
		return
	}

	responseData, err := json.Marshal(newUserData)
	if err != nil {
		log.Printf("[ERROR] Encoding error: %v\n", err)
		http.Error(w, "Encoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	reqUserId := r.Context().Value("reqUserId")

	if reqUserId != userId {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	user := models.DeleteUser(userId)
	responseData, err := json.Marshal(user)
	if err != nil {
		log.Printf("[ERROR] Encoding error: %v\n", err)
		http.Error(w, "Encoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(responseData)
}
