package main

import (
	"net/http"

	"github.com/rudbast/go-tour-rest/config"
	"github.com/rudbast/go-tour-rest/routers"
	"github.com/rudbast/go-tour-rest/util"
	"github.com/urfave/negroni"
)

func init() {
	// Initialize configurations.
	if err := config.Init(); err != nil {
		panic(err.Error())
	}

	// Initialize database.
	if err := util.InitDatabase(); err != nil {
		panic(err.Error())
	}
}

func main() {
	router := routers.InitRoutes()
	n := negroni.Classic()
	n.UseHandler(router)
	http.ListenAndServe(":8080", n)
}
