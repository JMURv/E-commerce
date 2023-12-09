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

func ListCreateItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		itemsList := models.GetAllItems()
		responseData, err := json.Marshal(itemsList)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	case http.MethodPost:
		NewItem := &models.Item{}
		utils.ParseBody(r, NewItem)

		item, err := NewItem.CreateItem()
		if err != nil {
			log.Printf("[ERROR] Creating item error: %v\n", err)
			http.Error(w, fmt.Sprintf("Creating item error: %v", err), http.StatusBadRequest)
			return
		}

		responseData, err := json.Marshal(item)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
	}
}

func GetUpdateDeleteItem(w http.ResponseWriter, r *http.Request) {
	itemId := mux.Vars(r)["id"]

	switch r.Method {
	case http.MethodGet:
		itemDetails := models.GetItemByID(itemId)
		responseData, err := json.Marshal(itemDetails)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)

	case http.MethodPut:
		var updateItem = &models.Item{}
		utils.ParseBody(r, updateItem)

		itemToUpdate := models.GetItemByID(itemId)
		updatedItem, err := itemToUpdate.UpdateItem(updateItem)
		if err != nil {
			log.Printf("[ERROR] Updating item error: %v\n", err)
			http.Error(w, fmt.Sprintf("Updating item error: %v", err), http.StatusBadRequest)
			return
		}

		responseData, err := json.Marshal(updatedItem)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Encoding error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(responseData)

	case http.MethodDelete:
		item := models.DeleteItem(itemId)
		responseData, err := json.Marshal(item)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Encoding error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Write(responseData)
	}
}
