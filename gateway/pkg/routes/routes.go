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

var RegisterAuthRoutes = func(router *mux.Router) {
	router.HandleFunc("/login", controllers.LoginHandler).Methods(http.MethodPost)
}

var RegisterUsersRoutes = func(router *mux.Router) {
	ctrl := controllers.NewUserCtrl()
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
	router.HandleFunc("/reviews", auth.AuthMiddleware(controllers.CreateReview)).Methods(http.MethodPost)

	router.HandleFunc("/reviews/{id}", controllers.GetReview).Methods(http.MethodGet)
	router.HandleFunc("/reviews/{id}", auth.AuthMiddleware(controllers.UpdateReview)).Methods(http.MethodPut)
	router.HandleFunc("/reviews/{id}", auth.AuthMiddleware(controllers.DeleteReview)).Methods(http.MethodDelete)
	// TODO: рассчитать общий рейтинг для юзера основываясь на его проданных товарах
}

var RegisterFavoriteRoutes = func(router *mux.Router) {
	router.HandleFunc("/favorites", auth.AuthMiddleware(controllers.ListFavorites)).Methods(http.MethodGet)
	router.HandleFunc("/favorites", auth.AuthMiddleware(controllers.CreateFavorites)).Methods(http.MethodPost)
	router.HandleFunc("/favorites/{id}", auth.AuthMiddleware(controllers.DeleteFavorite)).Methods(http.MethodDelete)
}

var RegisterTagsRoutes = func(router *mux.Router) {
	router.HandleFunc("/tags", controllers.ListTags).Methods(http.MethodGet)
	router.HandleFunc("/tags", auth.AuthMiddleware(controllers.CreateTag)).Methods(http.MethodPost)
	router.HandleFunc("/tags/{id}", auth.AuthMiddleware(controllers.DeleteTag)).Methods(http.MethodDelete)
}
