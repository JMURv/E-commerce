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
	router.HandleFunc("/users", controllers.ListCreateUsers).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/users/{id}", controllers.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", auth.AuthMiddleware(controllers.UpdateUser)).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", auth.AuthMiddleware(controllers.DeleteUser)).Methods(http.MethodDelete)
}

var RegisterSellersRoutes = func(router *mux.Router) {
	router.HandleFunc("/sellers", controllers.ListCreateSeller).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/sellers/{id}", controllers.GetUpdateDeleteSeller).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
}

var RegisterItemsRoutes = func(router *mux.Router) {
	router.HandleFunc("/items", controllers.ListCreateItems).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/items/{id}", controllers.GetUpdateDeleteItem).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
}

var RegisterCategoriesRoutes = func(router *mux.Router) {
	router.HandleFunc("/categories", controllers.ListCreateCategory).Methods(http.MethodGet, http.MethodPost)
}

//var RegisterReviewsRoutes = func(router *mux.Router) {
//	router.HandleFunc("/items/{itemID}/reviews", controllers.ListCreateReview).Methods(http.MethodGet, http.MethodPost)
//	router.HandleFunc("/reviews/{id}", controllers.GetUpdateDeleteReview).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
//}
//
//var RegisterCartsRoutes = func(router *mux.Router) {
//	router.HandleFunc("/carts", controllers.ListCarts).Methods(http.MethodGet, http.MethodPost)
//	router.HandleFunc("/carts/{id}", controllers.GetUpdateDeleteCart).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
//}
//
//var RegisterOrdersRoutes = func(router *mux.Router) {
//	router.HandleFunc("/orders", controllers.ListCreateOrder).Methods(http.MethodGet, http.MethodPost)
//	router.HandleFunc("/orders/{id}", controllers.GetUpdateDeleteOrder).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
//}
//
//var RegisterOrderItemsRoutes = func(router *mux.Router) {
//	router.HandleFunc("/orders/{orderID}/items", controllers.ListCreateOrderItem).Methods(http.MethodGet, http.MethodPost)
//	router.HandleFunc("/order_items/{id}", controllers.GetUpdateDeleteOrderItem).Methods(http.MethodGet, http.MethodPut, http.MethodDelete)
//}
