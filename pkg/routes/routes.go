package routes

import (
	"e-commerce/pkg/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

var RegisterItemsRoutes = func(router *mux.Router) {
	router.HandleFunc("/items", controllers.ListCreateItems).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/items/{id}", controllers.GetUpdateDeleteItem).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
}
