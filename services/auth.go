package services

import (
	"encoding/json"
	"net/http"

	"github.com/rudbast/go-tour-rest/api/parameters"
	"github.com/rudbast/go-tour-rest/core/auth"
	"github.com/rudbast/go-tour-rest/models"
)

func Login(user *models.User) (int, []byte) {
	authBackend := auth.GetAuthBackend()

	if userId, ok := authBackend.Authenticate(user); ok {
		if token, err := authBackend.GenerateToken(userId); err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(parameters.TokenAuthentication{Token: token})
			return http.StatusOK, response
		}
	} else {
		return http.StatusInternalServerError, []byte("")
	}
}

func RefreshToken(user *models.User) []byte {
	authBackend := auth.GetAuthBackend()

	token, err := authBackend.GenerateToken(user.Id)
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(parameters.TokenAuthentication{Token: token})
	if err != nil {
		panic(err)
	}

	return response
}

func Logout(req *http.Request) error {
	authBackend := auth.GetAuthBackend()

	tokenString := req.Header.Get("Authorization")

	return authBackend.InvalidateToken(tokenString, req)
}
