package main

import (
	"e-commerce/pkg/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
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
	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, router))
}
