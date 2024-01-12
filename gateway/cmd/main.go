package main

import (
	"github.com/JMURv/market/gateway/pkg/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	routes.RegisterAuthRoutes(router)

	routes.RegisterUsersRoutes(router)

	routes.RegisterItemsRoutes(router)
	routes.RegisterCategoriesRoutes(router)

	routes.RegisterFavoriteRoutes(router)
	routes.RegisterReviewsRoutes(router)

	routes.RegisterTagsRoutes(router)

	routes.RegisterWSRoutes(router)
	routes.RegisterRoomRoutes(router)
	routes.RegisterMessageRoutes(router)

	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, router)))
}
