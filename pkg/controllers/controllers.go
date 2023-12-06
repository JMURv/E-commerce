package controllers

import (
	"e-commerce/pkg/models"
	"e-commerce/pkg/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func ListCreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		return
		//var categories []models.Category
		//result := db.Find(&categories)
	}
}

func ListCreateItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		newItems := models.GetAllItems()
		res, _ := json.Marshal(newItems)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
		return
	} else if r.Method == "POST" {
		CreateItem := &models.Item{}
		utils.ParseBody(r, CreateItem)
		item := CreateItem.CreateItem()
		res, _ := json.Marshal(item)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(res)
	}
}

func GetUpdateDeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		vars := mux.Vars(r)
		itemID := vars["id"]

		id, err := strconv.ParseInt(itemID, 0, 0)
		if err != nil {
			log.Println("Parse error", err.Error())
		}

		itemDetails, _ := models.GetItemByID(id)
		res, _ := json.Marshal(itemDetails)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	} else if r.Method == http.MethodPut {
		var updateItem = &models.Item{}
		utils.ParseBody(r, updateItem)

		vars := mux.Vars(r)
		itemID := vars["id"]

		id, err := strconv.ParseInt(itemID, 0, 0)
		if err != nil {
			log.Println("Parse error", err.Error())
		}

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

	} else if r.Method == http.MethodDelete {
		vars := mux.Vars(r)
		itemID := vars["id"]

		id, err := strconv.ParseInt(itemID, 0, 0)
		if err != nil {
			log.Println("Parse error", err.Error())
		}

		item := models.DeleteItem(id)
		res, _ := json.Marshal(item)

		w.WriteHeader(http.StatusNoContent)
		w.Write(res)
	}
}
