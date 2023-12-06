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

func ListCreateCategory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		categoriesList := models.GetAllCategories()
		responseData, _ := json.Marshal(categoriesList)

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
		responseData, _ := json.Marshal(category)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseData)
	}
}

func ListCreateItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		itemsList := models.GetAllItems()
		res, _ := json.Marshal(itemsList)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	case http.MethodPost:
		NewItem := &models.Item{}
		utils.ParseBody(r, NewItem)

		item, err := NewItem.CreateItem()
		if err != nil {
			log.Printf("[ERROR] Creating item error: %v\n", err)
			http.Error(w, fmt.Sprintf("Creating item error: %v", err), http.StatusBadRequest)
			return
		}

		res, _ := json.Marshal(item)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}

func GetUpdateDeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 0, 0)
	if err != nil {
		log.Println("Parse error", err.Error())
	}

	switch r.Method {
	case http.MethodGet:
		itemDetails, _ := models.GetItemByID(id)
		res, _ := json.Marshal(itemDetails)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	case http.MethodPut:
		var updateItem = &models.Item{}
		utils.ParseBody(r, updateItem)

		itemDetails, db := models.GetItemByID(id)

		if updateItem.Name != "" {
			itemDetails.Name = updateItem.Name
		}

		if updateItem.Description != "" {
			itemDetails.Description = updateItem.Description
		}

		if updateItem.Price != 0 {
			itemDetails.Price = updateItem.Price
		}

		if updateItem.CategoryID != itemDetails.CategoryID {
			itemDetails.CategoryID = updateItem.CategoryID
		}

		db.Save(&itemDetails)
		res, _ := json.Marshal(itemDetails)
		w.WriteHeader(http.StatusOK)
		w.Write(res)

	case http.MethodDelete:
		item := models.DeleteItem(id)
		res, _ := json.Marshal(item)

		w.WriteHeader(http.StatusNoContent)
		w.Write(res)
	}
}
