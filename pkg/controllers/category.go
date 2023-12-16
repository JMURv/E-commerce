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

func ListCategory(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetAllCategories()
	if err != nil {
		http.Error(w, "Cannot get categories", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		http.Error(w, "Cannot parse categoryID", http.StatusInternalServerError)
		return
	}

	category, err := models.GetCategoryByID(uint(categoryID))
	if err != nil {
		http.Error(w, "Invalid categoryID", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(category)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	newCategory := &models.Category{}
	utils.ParseBody(r, newCategory)

	category, err := newCategory.CreateCategory()
	if err != nil {
		http.Error(w, fmt.Sprintf("Creating category error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(category)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		http.Error(w, "Cannot parse categoryID", http.StatusInternalServerError)
		return
	}

	newData := &models.Category{}
	utils.ParseBody(r, newData)

	newCategoryData, err := models.UpdateCategory(uint(categoryID), newData)
	if err != nil {
		http.Error(w, "Cannot update category", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(newCategoryData)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		http.Error(w, "Cannot parse categoryID", http.StatusInternalServerError)
		return
	}

	deletedCategory, err := models.DeleteCategory(uint(categoryID))
	if err != nil {
		http.Error(w, "Cannot delete category", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(deletedCategory)
	if err != nil {
		http.Error(w, "Encoding error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(response)
}
