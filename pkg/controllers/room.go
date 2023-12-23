package controllers

import (
	"e-commerce/pkg/models"
	"e-commerce/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func ListRoom(w http.ResponseWriter, r *http.Request) {
	reqUserID := r.Context().Value("reqUserId").(uint)

	rooms, err := models.GetUserRoomWithMessages(reqUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while getting user's rooms: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(rooms)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	utils.ResponseOk(w, http.StatusOK, response)
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		User1ID uint `json:"user1ID"`
		User2ID uint `json:"user2ID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	room, err := models.CreateRoom(requestData.User1ID, requestData.User2ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Creating room error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(room)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while encoding JSON: %v", err), http.StatusInternalServerError)
		return
	}

	utils.ResponseOk(w, http.StatusCreated, response)
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {

}
