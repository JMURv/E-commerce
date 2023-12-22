package test

import (
	"e-commerce/pkg/controllers"
	"github.com/gorilla/mux"
)

func createTestRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/users", controllers.ListCreateUser).Methods("POST")
	r.HandleFunc("/items", controllers.CreateItem).Methods("POST")
	return r
}
