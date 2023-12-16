package routes

import (
	"e-commerce/pkg/auth"
	"e-commerce/pkg/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

var RegisterAuthRoutes = func(router *mux.Router) {
	router.HandleFunc("/login", controllers.LoginHandler).Methods(http.MethodPost)
}

var RegisterUsersRoutes = func(router *mux.Router) {
	router.HandleFunc("/users", controllers.ListCreateUser).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/users/{id}", controllers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", auth.AuthMiddleware(controllers.UpdateUser)).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", auth.AuthMiddleware(controllers.DeleteUser)).Methods(http.MethodDelete)
}

var RegisterSellersRoutes = func(router *mux.Router) {
	router.HandleFunc("/sellers", controllers.ListSeller).Methods(http.MethodGet)
	router.HandleFunc("/sellers", controllers.CreateSeller).Methods(http.MethodPost)

	router.HandleFunc("/sellers/{id}", controllers.GetSeller).Methods(http.MethodGet)
	router.HandleFunc("/sellers/{id}", auth.AuthMiddleware(controllers.UpdateSeller)).Methods(http.MethodPut)
	router.HandleFunc("/sellers/{id}", auth.AuthMiddleware(controllers.DeleteSeller)).Methods(http.MethodDelete)

	router.HandleFunc("/sellers/{id}/items", controllers.ListSellerItems).Methods(http.MethodGet)

	router.HandleFunc("/sellers/{id}/items/{itemId}", auth.AuthMiddleware(controllers.LinkItemToSeller)).Methods(http.MethodPost)

	router.HandleFunc("/sellers/{id}/items/{sellerItemId}", auth.AuthMiddleware(controllers.GetSellerItem)).Methods(http.MethodGet)
	router.HandleFunc("/sellers/{id}/items/{sellerItemId}", auth.AuthMiddleware(controllers.UpdateSellerItem)).Methods(http.MethodPut)
	router.HandleFunc("/sellers/{id}/items/{sellerItemId}", auth.AuthMiddleware(controllers.DeleteSellerItem)).Methods(http.MethodDelete)
}

var RegisterItemsRoutes = func(router *mux.Router) {
	router.HandleFunc("/items", controllers.ListItem).Methods(http.MethodGet)
	router.HandleFunc("/items", auth.AuthMiddleware(controllers.CreateItem)).Methods(http.MethodPost)

	router.HandleFunc("/items/{id}", controllers.GetItem).Methods(http.MethodGet)
	router.HandleFunc("/items/{id}", auth.AuthMiddleware(controllers.UpdateItem)).Methods(http.MethodPut)
	router.HandleFunc("/items/{id}", auth.AuthMiddleware(controllers.DeleteItem)).Methods(http.MethodDelete)
}

var RegisterCategoriesRoutes = func(router *mux.Router) {
	router.HandleFunc("/categories", controllers.ListCategory).Methods(http.MethodGet)
	router.HandleFunc("/categories", auth.AuthMiddleware(controllers.CreateCategory)).Methods(http.MethodPost)

	router.HandleFunc("/categories{id}", controllers.GetCategory).Methods(http.MethodGet)
	router.HandleFunc("/categories{id}", auth.AuthMiddleware(controllers.UpdateCategory)).Methods(http.MethodPut)
	router.HandleFunc("/categories{id}", auth.AuthMiddleware(controllers.DeleteCategory)).Methods(http.MethodDelete)
}

var RegisterReviewsRoutes = func(router *mux.Router) {
	router.HandleFunc("/reviews", auth.AuthMiddleware(controllers.CreateReview)).Methods(http.MethodPost)

	router.HandleFunc("/reviews/{id}", controllers.GetReview).Methods(http.MethodGet)
	router.HandleFunc("/reviews/{id}", auth.AuthMiddleware(controllers.UpdateReview)).Methods(http.MethodPut)
	router.HandleFunc("/reviews/{id}", auth.AuthMiddleware(controllers.DeleteReview)).Methods(http.MethodDelete)
}

var RegisterFavoriteRoutes = func(router *mux.Router) {
	router.HandleFunc("/favorites", auth.AuthMiddleware(controllers.ListUserFavorites)).Methods(http.MethodGet)
	router.HandleFunc("/favorites", auth.AuthMiddleware(controllers.CreateFavorites)).Methods(http.MethodPost)

	router.HandleFunc("/favorites/{id}", auth.AuthMiddleware(controllers.GetFavorites)).Methods(http.MethodGet)
	router.HandleFunc("/favorites/{id}", auth.AuthMiddleware(controllers.UpdateFavorites)).Methods(http.MethodPut)
	router.HandleFunc("/favorites/{id}", auth.AuthMiddleware(controllers.DeleteFavorites)).Methods(http.MethodDelete)
}

var RegisterCartsRoutes = func(router *mux.Router) {
	router.HandleFunc("/carts", auth.AuthMiddleware(controllers.GetCart)).Methods(http.MethodGet)
	router.HandleFunc("/carts/{id}", auth.AuthMiddleware(controllers.UpdateCart)).Methods(http.MethodPut)
	router.HandleFunc("/carts/{id}", auth.AuthMiddleware(controllers.DeleteCart)).Methods(http.MethodDelete)

	router.HandleFunc("/carts/items", auth.AuthMiddleware(controllers.GetCartItems)).Methods(http.MethodGet)
	router.HandleFunc("/carts/items/{itemId}", auth.AuthMiddleware(controllers.AddOrDeleteCartItem)).Methods(http.MethodPut)
	router.HandleFunc("/carts/items/{itemId}", auth.AuthMiddleware(controllers.DeleteCartItem)).Methods(http.MethodDelete)
}

var RegisterOrdersRoutes = func(router *mux.Router) {
	router.HandleFunc("/orders", auth.AuthMiddleware(controllers.ListOrder)).Methods(http.MethodGet)
	router.HandleFunc("/orders", auth.AuthMiddleware(controllers.CreateOrder)).Methods(http.MethodPost)

	router.HandleFunc("/orders/{id}", auth.AuthMiddleware(controllers.GetOrder)).Methods(http.MethodGet)
	router.HandleFunc("/orders/{id}", auth.AuthMiddleware(controllers.UpdateOrder)).Methods(http.MethodPut)
	router.HandleFunc("/orders/{id}", auth.AuthMiddleware(controllers.DeleteOrder)).Methods(http.MethodDelete)
}
