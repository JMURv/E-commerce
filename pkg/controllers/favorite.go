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

}

func ListUserFavorites(w http.ResponseWriter, r *http.Request) {

}

func CreateFavorites(w http.ResponseWriter, r *http.Request) {
	newFavorite := &models.Favorite{}
	utils.ParseBody(r, newFavorite)

	favorite, err := newFavorite.CreateFavorite()
	if err != nil {
		http.Error(w, fmt.Sprintf("Creating review error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(favorite)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	favoriteID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "Cannot parse categoryID", http.StatusInternalServerError)
		return
	}

	favorite, err := models.GetFavoriteByID(uint(favoriteID))
	if err != nil {
		http.Error(w, "Cannot get review error", http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(favorite)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func UpdateFavorites(w http.ResponseWriter, r *http.Request) {

}

func DeleteFavorites(w http.ResponseWriter, r *http.Request) {

}
