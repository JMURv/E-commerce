package routes

import (
	"github.com/JMURv/e-commerce/gateway/pkg/auth"
	"github.com/JMURv/e-commerce/gateway/pkg/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

var RegisterWSRoutes = func(router *mux.Router) {
	router.HandleFunc("/ws", controllers.HandleWebSocket)
}

var RegisterRoomRoutes = func(router *mux.Router) {
	router.HandleFunc("/rooms", auth.AuthMiddleware(controllers.ListRoom)).Methods(http.MethodGet)
	router.HandleFunc("/rooms", auth.AuthMiddleware(controllers.CreateRoom)).Methods(http.MethodPost)

	router.HandleFunc("/rooms/{id}", auth.AuthMiddleware(controllers.DeleteRoom)).Methods(http.MethodDelete)
}

var RegisterMessageRoutes = func(router *mux.Router) {
	router.HandleFunc("/messages", auth.AuthMiddleware(controllers.ListMessage)).Methods(http.MethodGet)
	router.HandleFunc("/messages", auth.AuthMiddleware(controllers.CreateMessage)).Methods(http.MethodPost)

	router.HandleFunc("/messages/{id}", auth.AuthMiddleware(controllers.UpdateMessage)).Methods(http.MethodPut)
	router.HandleFunc("/messages/{id}", auth.AuthMiddleware(controllers.DeleteMessage)).Methods(http.MethodDelete)
}

var RegisterUsersRoutes = func(router *mux.Router) {
	ctrl := controllers.NewUserCtrl()
	router.HandleFunc("/login", ctrl.LoginHandler).Methods(http.MethodPost)

	router.HandleFunc("/users", ctrl.ListCreateUser).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/users/{id}", ctrl.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", auth.AuthMiddleware(ctrl.UpdateUser)).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", auth.AuthMiddleware(ctrl.DeleteUser)).Methods(http.MethodDelete)
}

var RegisterItemsRoutes = func(router *mux.Router) {
	ctrl := controllers.NewItemCtrl()
	router.HandleFunc("/items", ctrl.ListItem).Methods(http.MethodGet)
	router.HandleFunc("/items", auth.AuthMiddleware(ctrl.CreateItem)).Methods(http.MethodPost)

	router.HandleFunc("/items/{id}", ctrl.GetItem).Methods(http.MethodGet)
	router.HandleFunc("/items/{id}", auth.AuthMiddleware(ctrl.UpdateItem)).Methods(http.MethodPut)
	router.HandleFunc("/items/{id}", auth.AuthMiddleware(ctrl.DeleteItem)).Methods(http.MethodDelete)
}

var RegisterCategoriesRoutes = func(router *mux.Router) {
	ctrl := controllers.NewCategoryCtrl()
	router.HandleFunc("/categories", ctrl.ListCategory).Methods(http.MethodGet)
	router.HandleFunc("/categories", auth.AuthMiddleware(ctrl.CreateCategory)).Methods(http.MethodPost)

	router.HandleFunc("/categories{id}", ctrl.GetCategory).Methods(http.MethodGet)
	router.HandleFunc("/categories{id}", auth.AuthMiddleware(ctrl.UpdateCategory)).Methods(http.MethodPut)
	router.HandleFunc("/categories{id}", auth.AuthMiddleware(ctrl.DeleteCategory)).Methods(http.MethodDelete)
}

var RegisterReviewsRoutes = func(router *mux.Router) {
	ctrl := controllers.NewReviewCtrl()
	router.HandleFunc("/reviews", auth.AuthMiddleware(ctrl.CreateReview)).Methods(http.MethodPost)

	router.HandleFunc("/reviews/{id}", ctrl.GetReview).Methods(http.MethodGet)
	router.HandleFunc("/reviews/{id}", auth.AuthMiddleware(ctrl.UpdateReview)).Methods(http.MethodPut)
	router.HandleFunc("/reviews/{id}", auth.AuthMiddleware(ctrl.DeleteReview)).Methods(http.MethodDelete)
	// TODO: рассчитать общий рейтинг для юзера основываясь на его проданных товарах
}

var RegisterFavoriteRoutes = func(router *mux.Router) {
	ctrl := controllers.NewFavoriteCtrl()
	router.HandleFunc("/favorites", auth.AuthMiddleware(ctrl.ListFavorites)).Methods(http.MethodGet)
	router.HandleFunc("/favorites", auth.AuthMiddleware(ctrl.CreateFavorites)).Methods(http.MethodPost)
	router.HandleFunc("/favorites/{id}", auth.AuthMiddleware(ctrl.DeleteFavorite)).Methods(http.MethodDelete)
}

var RegisterTagsRoutes = func(router *mux.Router) {
	ctrl := controllers.NewTagCtrl()
	router.HandleFunc("/tags", ctrl.ListTags).Methods(http.MethodGet)
	router.HandleFunc("/tags", auth.AuthMiddleware(ctrl.CreateTag)).Methods(http.MethodPost)
	router.HandleFunc("/tags/{name}", auth.AuthMiddleware(ctrl.DeleteTag)).Methods(http.MethodDelete)
}
