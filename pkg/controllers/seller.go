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

func ListSeller(w http.ResponseWriter, r *http.Request) {
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
}

func CreateSeller(w http.ResponseWriter, r *http.Request) {
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

func GetSeller(w http.ResponseWriter, r *http.Request) {
	sellerId := mux.Vars(r)["id"]
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
}

func UpdateSeller(w http.ResponseWriter, r *http.Request) {
	sellerId := mux.Vars(r)["id"]
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
}

func DeleteSeller(w http.ResponseWriter, r *http.Request) {
	sellerId := mux.Vars(r)["id"]
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

func ListSellerItems(w http.ResponseWriter, r *http.Request) {

}

func CreateSellerItem(w http.ResponseWriter, r *http.Request) {

}

func UpdateSellerItem(w http.ResponseWriter, r *http.Request) {

}

func DeleteSellerItem(w http.ResponseWriter, r *http.Request) {

}
