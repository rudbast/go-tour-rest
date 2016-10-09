package controllers

import (
	"net/http"
)

// Get list of all persons.
func Persons(rw http.ResponseWriter, rq *http.Request, next http.HandlerFunc) {
	// TODO: Access db.
}
