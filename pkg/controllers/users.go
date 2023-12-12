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

func ListCreateUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		usersList := models.GetAllUsers()
		responseData, err := json.Marshal(usersList)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)

	case http.MethodPost:
		var newUserData = &models.User{}
		utils.ParseBody(r, newUserData)

		user, err := newUserData.CreateUser()
		if err != nil {
			log.Printf("[ERROR] Creating user error: %v\n", err)
			http.Error(w, fmt.Sprintf("Creating user error: %v", err), http.StatusBadRequest)
			return
		}

		responseData, err := json.Marshal(user)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Encoding error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]
	userDetails := models.GetUserByID(userId)

	responseData, err := json.Marshal(userDetails)
	if err != nil {
		log.Printf("[ERROR] Encoding error: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
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
	w.WriteHeader(http.StatusNoContent)
	w.Write(responseData)
}
