package controllers

import (
	"e-commerce/pkg/models"
	"e-commerce/pkg/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
	sellerDetails := models.GetSellerByID(sellerId)

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
	sellerID := mux.Vars(r)["id"]
	seller := models.GetSellerWithItems(sellerID)

	if seller == nil {
		http.Error(w, "Seller not found", http.StatusNotFound)
		return
	}

	items := seller.SellerItems
	jsonResponse, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func LinkItemToSeller(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sellerID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "cannot parse sellerID", http.StatusInternalServerError)
		return
	}

	itemID, err := strconv.ParseUint(vars["itemId"], 10, 64)
	if err != nil {
		http.Error(w, "cannot parse itemID", http.StatusInternalServerError)
		return
	}

	var newItemToSeller = &models.SellerItem{
		SellerID: sellerID,
		ItemID:   itemID,
	}

	utils.ParseBody(r, newItemToSeller)
	newItemToSeller, err = newItemToSeller.LinkItemToSeller()
	if err != nil {
		http.Error(w, "cannot link item to seller", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(newItemToSeller)
	if err != nil {
		http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func UpdateSellerItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sellerID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "cannot parse sellerID", http.StatusInternalServerError)
		return
	}

	itemID, err := strconv.ParseUint(vars["itemId"], 10, 64)
	if err != nil {
		http.Error(w, "cannot parse itemID", http.StatusInternalServerError)
		return
	}

	var newItemToSeller = &models.SellerItem{
		SellerID: sellerID,
		ItemID:   itemID,
	}

	utils.ParseBody(r, newItemToSeller)
	newItemToSeller, err = newItemToSeller.LinkItemToSeller()
	if err != nil {
		http.Error(w, "cannot link item to seller", http.StatusInternalServerError)
		return
	}
}

func DeleteSellerItem(w http.ResponseWriter, r *http.Request) {

}
