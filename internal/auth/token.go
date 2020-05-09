package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"monkiato/guardian/internal/models"
)

/*
	How to use JwtToken

	mySecretKey := "my-secret-key"
	tokenExpirationTime := time.Hours * 24
	mySessionUser := models.SessionUser{...}

	// create a new one from user data
	jwtToken := auth.NewJwtToken(mySessionUser, mySecretKey, expirationTime)
	strToken := jwtToken.ToString()

	// create a new one from an existing token string, secretKey is required for encoding/decoding
	auth.NewJwtTokenFromString(strToken, mySecretKey)

	// check if token is valid:
	if jwtToken.Token == nil || jwtToken.Token.Valid {
		...
	}
 */

// JwtToken wrapper for jwt.Token handling, using a specific secretKey and own User model
type JwtToken struct {
	Token     *jwt.Token
	secretKey string
}

// NewJwtToken create JwtToken instance with the specified secretKey
func NewJwtToken(user *models.SessionUser, secretKey string, expirationTime int64) *JwtToken {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = &TokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
		*user,
	}

	return &JwtToken{
		Token: token,
		secretKey: secretKey,
	}
}

func NewJwtTokenFromString(strToken string, secretKey string) (*JwtToken, error) {
	jwtToken, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error parsing the token")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return &JwtToken{
		Token: jwtToken,
		secretKey: secretKey,
	}, nil
}

// ToString string representation for existing kwt.Token
func (t *JwtToken) ToString() (string, error) {
	tokenString, err := t.Token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", errors.New("unable to create token")
	}
	return tokenString, nil
}


// GetUser extract user data from token, it unmarshal a SessionUser object
func (t *JwtToken) GetUser() (*models.SessionUser, error) {
	var user models.SessionUser
	// decode token claim data (user data)
	mapstructure.Decode(t.Token.Claims, &user)
	fmt.Printf("Claim info. user: %s, name: %s, lastname: %s", user.Username, user.Name, user.Lastname)

	// return SessionUser
	return &user, nil
}
