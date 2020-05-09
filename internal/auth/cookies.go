package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/securecookie"
	"monkiato/guardian/internal/models"
)

//DefaultCookieName default name used to authentication cookies
const DefaultCookieName string = "guardian-auth"

// Handler JWT and cookie handling
type Handler struct {
	domainName string
	secretKey  string
	cookieName string
}

// NewHandler create new Handler instance
func NewHandler(domainName string, secretKey string) *Handler {
	return &Handler{
		domainName: domainName,
		secretKey:  secretKey,
		cookieName: getCookieName(),
	}
}

// ReadCookie extract cookie info from http request
func (handler *Handler) ReadCookie(req *http.Request) (*http.Cookie, *JwtToken, error) {
	// get cookie data
	cookie, err := req.Cookie(handler.cookieName)
	if err != nil {
		fmt.Println("cookie not found")
		return nil, nil, err
	}

	token, err := handler.GetToken(cookie)
	if err != nil {
		fmt.Println("error reading token from cookie")
		return nil, nil, err
	}

	jwtToken, err := NewJwtTokenFromString(token, handler.secretKey)
	if err != nil {
		fmt.Println("error creating JwtToken from string token")
		return nil, nil, err
	}

	return cookie, jwtToken, nil
}

// CreateCookie create new cookie based on user data, expirationTime is also a requirement for the cookie
func (handler *Handler) CreateCookie(user *models.SessionUser, expirationTime int64) (*http.Cookie, string, error) {
	// create JWT token
	token := NewJwtToken(user, handler.secretKey, expirationTime)
	strToken, err := token.ToString()
	if err != nil {
		return nil, "", err
	}

	// create cookie
	cookies := map[string]string{
		"token": strToken,
	}
	encoded, err := handler.getSecureCookie().Encode(handler.cookieName, cookies)
	if err != nil {
		return nil, "", errors.New("Couldn't create cookies. " + err.Error())
	}

	cookie := &http.Cookie{
		Name:   handler.cookieName,
		Value:  encoded,
		Path:   "/",
		Domain: handler.domainName,
	}
	return cookie, strToken, nil
}

// GetToken extract token as string from cookie
func (handler *Handler) GetToken(cookie *http.Cookie) (string, error) {
	value := make(map[string]string)
	// decode cookie
	if err := handler.getSecureCookie().Decode(handler.cookieName, cookie.Value, &value); err != nil {
		fmt.Println("unable to decode secure cookie")
		return "", err
	}
	return value["token"], nil
}

func (handler *Handler) getSecureCookie() *securecookie.SecureCookie {
	var hashKey = []byte(handler.secretKey)
	return securecookie.New(hashKey, nil)
}

// getCookieName get cookie name used to lookup information
func getCookieName() string {
	cookieName, found := os.LookupEnv("COOKIE_NAME")
	if !found {
		cookieName = DefaultCookieName
	}
	return cookieName
}
