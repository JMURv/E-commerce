package test

import (
	"bytes"
	"e-commerce/pkg/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateUser(t *testing.T) {
	newUser := models.User{
		Username: "TestUsername",
		Email:    "TestEmail@email.com",
	}

	jsonPayload, err := json.Marshal(newUser)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := createTestRouter()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser models.UserWithToken
	err = json.Unmarshal(rr.Body.Bytes(), &createdUser)
	if err != nil {
		t.Fatal(err)
	}

	if createdUser.User.Username != newUser.Username {
		t.Errorf("Expected username %s, got %s", newUser.Username, createdUser.User.Username)
	}
}
