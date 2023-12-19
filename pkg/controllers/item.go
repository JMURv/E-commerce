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

func ListItem(w http.ResponseWriter, r *http.Request) {
	itemsList := models.GetAllItems()

	responseData, err := json.Marshal(itemsList)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	NewItem := &models.Item{}
	utils.ParseBody(r, NewItem)

	item, err := NewItem.CreateItem()
	if err != nil {
		http.Error(w, fmt.Sprintf("Creating item error: %v", err), http.StatusBadRequest)
		return
	}

	responseData, err := json.Marshal(item)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseData)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse itemId: %v", err), http.StatusInternalServerError)
		return
	}

	itemDetails := models.GetItemByID(uint(itemId))
	responseData, err := json.Marshal(itemDetails)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse itemId: %v", err), http.StatusInternalServerError)
		return
	}

	var updateItem = &models.Item{}
	utils.ParseBody(r, updateItem)

	itemToUpdate := models.GetItemByID(uint(itemId))

	reqUserId := r.Context().Value("reqUserId")
	if reqUserId != itemToUpdate.UserID {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	updatedItem, err := itemToUpdate.UpdateItem(updateItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("Updating item error: %v", err), http.StatusBadRequest)
		return
	}

	responseData, err := json.Marshal(updatedItem)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemId, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse itemId: %v", err), http.StatusInternalServerError)
		return
	}

	itemToDelete := models.GetItemByID(uint(itemId))
	reqUserId := r.Context().Value("reqUserId")
	if reqUserId != itemToDelete.UserID {
		http.Error(w, "you have no permissions", http.StatusForbidden)
		return
	}

	item, err := models.DeleteItem(uint(itemId))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot delete item error: %v", err), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(item)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(responseData)
}
