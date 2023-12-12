package main

import (
	"bytes"
	"e-commerce/pkg/controllers"
	"e-commerce/pkg/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateItem(t *testing.T) {
	newItem := models.Item{
		Name:        "TestProduct",
		Description: "A test product",
		Price:       29.99,
		Tags:        []models.Tag{{Name: "Tag1"}, {Name: "Tag2"}},
	}

	jsonPayload, err := json.Marshal(newItem)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := createTestRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdItem models.Item
	err = json.Unmarshal(rr.Body.Bytes(), &createdItem)
	if err != nil {
		t.Fatal(err)
	}

	if createdItem.Name != newItem.Name {
		t.Errorf("Expected item name %s, got %s", newItem.Name, createdItem.Name)
	}

	if len(createdItem.Tags) != 2 {
		t.Errorf("Expected 2 item tags got %v", len(createdItem.Tags))
	}
}

func createTestRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/items", controllers.CreateItem).Methods("POST")
	return r
}
