package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/JMURv/e-commerce/gateway/pkg/models"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func ListFavorites(w http.ResponseWriter, r *http.Request) {
	reqUserID := r.Context().Value("reqUserId").(uint64)

	favorites, err := models.GetAllUserFavorites(reqUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot get favorites: %v", err), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(favorites)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	utils.OkResponse(w, http.StatusOK, response)
}

func CreateFavorites(w http.ResponseWriter, r *http.Request) {
	newFavorite := &models.Favorite{}
	utils.ParseBody(r, newFavorite)

	favorite, err := newFavorite.CreateFavorite()
	if err != nil {
		http.Error(w, fmt.Sprintf("Creating favorite error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(favorite)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	// Create notification for item's author and send message via WS
	newNotification := models.Notification{
		Type:       "notification",
		UserID:     favorite.UserID,
		ReceiverID: uint64(favorite.Item.UserID),
		Message:    "new favorite",
	}

	notification, err := newNotification.CreateNotification()
	if err != nil {
		log.Printf("Error while creating notification: %v", err)
	}

	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error while encoding notification message: %v", err)
	}

	go broadcast(uint(favorite.UserID), favorite.Item.UserID, notificationBytes)

	utils.OkResponse(w, http.StatusCreated, response)
}

func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	favoriteID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse favoriteID: %v", err), http.StatusInternalServerError)
		return
	}

	favorite, err := models.GetFavoriteByID(favoriteID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot find favorite: %v", err), http.StatusNotFound)
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint64)
	if favorite.UserID != reqUserID {
		http.Error(w, "You have no permission to do it", http.StatusUnauthorized)
		return
	}

	err = models.DeleteFavorite(favoriteID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Deleting favorite error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
