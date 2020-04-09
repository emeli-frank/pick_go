package http

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"net/http"
)

func (s server) Routes() http.Handler {
	standardMiddleWare := alice.New(s.recoverPanic, s.logRequest)
	fooMiddleWare := alice.New(s.checkJWT, s.authenticatedOnly)
	checkJWTMiddleWare := alice.New(s.checkJWT)

	r := mux.NewRouter()

	r.HandleFunc("/api/register", s.registerHandler).Methods("POST")
	r.HandleFunc("/api/login", s.loginHandler).Methods("POST")
	r.Handle("/api/user",
		fooMiddleWare.Then(http.HandlerFunc(s.userHandler)))

	r.HandleFunc("/api/products", s.productListHandler)
	r.Handle("/api/products/{id}", checkJWTMiddleWare.Then(http.HandlerFunc(s.productDetailHandler)))

	r.Handle("/api/cart-items", fooMiddleWare.Then(http.HandlerFunc(s.saveProductToCartHandler))).
		Methods("POST")
	r.Handle("/api/cart-items", fooMiddleWare.Then(http.HandlerFunc(s.deleteProductFromCartHandler))).
		Methods("DELETE")
	r.Handle("/api/cart-items", fooMiddleWare.Then(http.HandlerFunc(s.cartItemsHandler)))

	r.Handle("/api/order-history", fooMiddleWare.Then(http.HandlerFunc(s.saveToOrderHistoryHandler))).
		Methods("POST")
	r.Handle("/api/order-history", fooMiddleWare.Then(http.HandlerFunc(s.orderHistoryHandler)))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // todo:: adjust before production
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders: []string{"*"},
	})
	return c.Handler(standardMiddleWare.Then(r))
	//return cors.Default().Handler(globalMiddleware.Then(r))
}