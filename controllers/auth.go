package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/rudbast/go-tour-rest/models"
	"github.com/rudbast/go-tour-rest/services"
)

// Login to system.
func Login(rw http.ResponseWriter, rq *http.Request) {
	var user models.User
	_ = json.NewDecoder(rq.Body).Decode(&user)

	responseStatus, token := services.Login(&user)
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(responseStatus)
	rw.Write(token)
}

// Refresh token from system.
func RefreshToken(rw http.ResponseWriter, rq *http.Request, next http.HandlerFunc) {
	var user models.User
	_ = json.NewDecoder(rq.Body).Decode(&user)

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(services.RefreshToken(&user))
}

// // Logout from system.
// func Logout(rw http.ResponseWriter, rq *http.Request, next http.HandlerFunc) {
// 	rw.Header().Set("Content-Type", "application/json")
// 	if err := services.Logout(rq); err != nil {
// 		rw.WriteHeader(http.StatusInternalServerError)
// 	} else {
// 		rw.WriteHeader(http.StatusOK)
// 	}
// }
