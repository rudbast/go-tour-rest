package routers

import (
	"github.com/gorilla/mux"
	"github.com/rudbast/go-tour-rest/controllers"
	"github.com/rudbast/go-tour-rest/core/auth"
	"github.com/urfave/negroni"
)

// Set authentication related routes.
func SetAuthenticationRoutes(router *mux.Router) {
	authBackend := auth.GetAuthBackend()

	router.HandleFunc("/login", controllers.Login).Methods("POST")

	router.Handle("/token/refresh",
		negroni.New(
			negroni.HandlerFunc(authBackend.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.RefreshToken),
		)).Methods("GET")

	router.Handle("/logout",
		negroni.New(
			negroni.HandlerFunc(authBackend.RequireTokenAuthentication),
			negroni.HandlerFunc(controllers.Logout),
		)).Methods("GET")
}
