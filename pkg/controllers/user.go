package controllers

import (
	"e-commerce/pkg/models"
	"e-commerce/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

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

	userDetails, err := models.GetUserByID(uint(userID))
	if err != nil {
		http.Error(w, fmt.Sprintf("User does not exist: %v", err), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(userDetails)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	utils.ResponseOk(w, http.StatusOK, response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse userID: %v", err), http.StatusInternalServerError)
		return
	}

	reqUserID := r.Context().Value("reqUserId")
	if reqUserID != uint(userID) {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	var updatedUser = &models.User{}
	utils.ParseBody(r, updatedUser)

	newUserData, err := models.UpdateUser(uint(userID), updatedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Updating user error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(newUserData)
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

	reqUserID := r.Context().Value("reqUserId")
	if reqUserID != uint(userID) {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	err = models.DeleteUser(uint(userID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Deleting user error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
