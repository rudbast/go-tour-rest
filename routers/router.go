package routers

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	router := mux.NewRouter()

	SetAuthenticationRoutes(router)
	// router = SetPersonRoutes(router)

	return router
}
