package auth

import (
	"github.com/rudbast/go-tour-rest/models"
	"net/http"
)

type Auth interface {
	// Authenticate user.
	Authenticate(user *models.User) (int, bool)
	// Generate token for authenticated user.
	GenerateToken(userId int) (string, error)
	// Initialize backend.
	InitBackend() error
	// Middleware function for routes that needed authentication.
	RequireTokenAuthentication(rw http.ResponseWriter, rq *http.Request, next http.HandlerFunc)
}
