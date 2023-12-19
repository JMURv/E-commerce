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

func ListFavorites(w http.ResponseWriter, r *http.Request) {
	reqUserID := r.Context().Value("reqUserId").(uint)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	favoriteID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse favoriteID: %v", err), http.StatusInternalServerError)
		return
	}

	favorite, err := models.GetFavoriteByID(uint(favoriteID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot find favorite: %v", err), http.StatusNotFound)
		return
	}

	reqUserID := r.Context().Value("reqUserId").(uint)
	if favorite.UserID != reqUserID {
		http.Error(w, "You have no permission to do it", http.StatusUnauthorized)
		return
	}

	err = models.DeleteFavorite(uint(favoriteID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Deleting favorite error: %v", err), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(favorite)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(responseData)
}
