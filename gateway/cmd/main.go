package main

import (
	"github.com/JMURv/e-commerce/gateway/pkg/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

	router := mux.NewRouter()
	routes.RegisterUsersRoutes(router)

	routes.RegisterItemsRoutes(router)
	routes.RegisterCategoriesRoutes(router)

	routes.RegisterFavoriteRoutes(router)
	routes.RegisterReviewsRoutes(router)

	routes.RegisterTagsRoutes(router)

	routes.RegisterWSRoutes(router)
	routes.RegisterRoomRoutes(router)
	routes.RegisterMessageRoutes(router)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully...")
		os.Exit(0)
	}()
	log.Println("Server is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, router)))
}
