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

func ListCreateSeller(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sellersList := models.GetAllSellers()
		responseData, err := json.Marshal(sellersList)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)

	case http.MethodPost:
		var newSellerData = &models.Seller{}
		utils.ParseBody(r, newSellerData)

		seller, err := newSellerData.CreateSeller()
		if err != nil {
			log.Printf("[ERROR] Creating seller error: %v\n", err)
			http.Error(w, fmt.Sprintf("Creating seller error: %v", err), http.StatusBadRequest)
			return
		}

		responseData, err := json.Marshal(seller)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Encoding error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
	}
}

func GetUpdateDeleteSeller(w http.ResponseWriter, r *http.Request) {
	sellerId := mux.Vars(r)["id"]

	switch r.Method {
	case http.MethodGet:
		sellerDetails := models.GetUserByID(sellerId)

		responseData, err := json.Marshal(sellerDetails)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)

	case http.MethodPut:
		var updatedSeller = &models.Seller{}
		utils.ParseBody(r, updatedSeller)

		newUserData, err := models.UpdateSeller(sellerId, updatedSeller)
		if err != nil {
			log.Printf("[ERROR] Updating seller error: %v\n", err)
			http.Error(w, fmt.Sprintf("Updating seller error: %v", err), http.StatusBadRequest)
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

	case http.MethodDelete:
		seller := models.DeleteSeller(sellerId)
		responseData, err := json.Marshal(seller)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Encoding error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		w.Write(responseData)
	}
}
