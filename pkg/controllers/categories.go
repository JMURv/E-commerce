package controllers

import (
	"e-commerce/pkg/models"
	"e-commerce/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ListCreateCategory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categoriesList := models.GetAllCategories()
		responseData, err := json.Marshal(categoriesList)
		if err != nil {
			log.Printf("[ERROR] Encoding error: %v\n", err)
			http.Error(w, "Encoding error", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseData)
	case http.MethodPost:
		NewCategory := &models.Category{}
		utils.ParseBody(r, NewCategory)

		category, err := NewCategory.CreateCategory()
		if err != nil {
			log.Printf("[ERROR] Creating category error: %v\n", err)
			http.Error(w, fmt.Sprintf("Creating category error: %v", err), http.StatusBadRequest)
			return
		}

		responseData, err := json.Marshal(category)
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
