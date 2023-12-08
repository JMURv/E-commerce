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
	routes.RegisterUsersRoutes(router)
	routes.RegisterSellersRoutes(router)

	routes.RegisterItemsRoutes(router)
	routes.RegisterCategoriesRoutes(router)

	//routes.RegisterCartsRoutes(router)
	//routes.RegisterOrdersRoutes(router)
	//routes.RegisterOrderItemsRoutes(router)

	//routes.RegisterReviewsRoutes(router)

	log.Println("Server is running on port 8000")
	http.ListenAndServe(":8000", router)
}
