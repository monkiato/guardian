package auth

import (
	"monkiato/guardian/internal/models"

	jwt "github.com/dgrijalva/jwt-go"
)

// TokenClaim this is the cliam object which gets parsed from the authorization header
type TokenClaim struct {
	*jwt.StandardClaims
	models.SessionUser
}

// ErrorMsg ...
// Custom error object
type ErrorMsg struct {
	Message string `json:"message"`
}
