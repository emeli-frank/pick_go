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

	r := mux.NewRouter()


	r.HandleFunc("/api/register", s.registerHandler).Methods("POST")
	r.HandleFunc("/api/login", s.loginHandler).Methods("POST")
	r.Handle("/api/user",
		fooMiddleWare.Then(http.HandlerFunc(s.userHandler)))

	/*r.Handle(
		"/api/confirmation-token",
		fooMiddleWare.Then(http.HandlerFunc(s.createEmailTokenConfirmationHandler)),
	).Methods("POST")
	r.HandleFunc("/api/confirmation-token/{token}", s.emailConfirmationTokenHandler).
		Methods("DELETE")*/

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // todo:: adjust before production
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders: []string{"*"},
	})
	return c.Handler(standardMiddleWare.Then(r))
	//return cors.Default().Handler(globalMiddleware.Then(r))
}