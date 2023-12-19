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

func ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := models.GetAllTags()
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot get tags: %v", err), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	newTag := &models.Tag{}
	utils.ParseBody(r, newTag)

	tag, err := newTag.CreateTag()
	if err != nil {
		http.Error(w, fmt.Sprintf("Creating tag error: %v", err), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(tag)
	if err != nil {
		http.Error(w, fmt.Sprintf("Encoding error: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	tagID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot parse tagID: %v", err), http.StatusInternalServerError)
		return
	}

	err = models.DeleteTag(uint(tagID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Deleting favorite error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
