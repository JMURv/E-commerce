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

func GetReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse reviewID: %v", err), http.StatusInternalServerError)
		return
	}

	review, err := models.GetReviewByID(uint(reviewID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot get review error: %v", err), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(review)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func CreateReview(w http.ResponseWriter, r *http.Request) {
	newReview := &models.Review{}
	utils.ParseBody(r, newReview)

	review, err := newReview.CreateReview()
	if err != nil {
		http.Error(w, fmt.Sprintf("Creating review error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(review)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func UpdateReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse reviewID: %v", err), http.StatusInternalServerError)
		return
	}

	newData := &models.Review{}
	utils.ParseBody(r, newData)

	newReviewData, err := models.UpdateReview(uint(reviewID), newData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot update review: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(newReviewData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse reviewID: %v", err), http.StatusInternalServerError)
		return
	}

	deletedReview, err := models.DeleteReview(uint(reviewID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot delete review: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(deletedReview)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write(response)
}
