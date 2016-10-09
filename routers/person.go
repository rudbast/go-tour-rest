package routers

import (
	"github.com/gorilla/mux"
	"github.com/rudbast/go-tour-rest/controllers"
	"github.com/urfave/negroni"
)

// Set person related-routes.
func SetPersonRoutes(router *mux.Router) *mux.Router {
	router.Handle("/persons",
		negroni.New(
			negroni.HandlerFunc(controllers.Persons),
		)).Methods("GET")

	return router
}
