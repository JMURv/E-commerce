package main

import (
	"e-commerce/pkg/routes"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	routes.RegisterItemsRoutes(router)
	routes.RegisterCategoriesRoutes(router)

	log.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", router)
}
