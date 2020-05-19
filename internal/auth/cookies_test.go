package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	domainName = "testing.com"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler(domainName, secretKey)
	if handler == nil {
		t.Fatalf("unexpected null instance")
	}
}

func TestHandler_CreateCookie(t *testing.T) {
	handler := NewHandler(domainName, secretKey)
	cookie, token, err := handler.CreateCookie(&user, expirationTime)
	if err != nil {
		t.Fatalf("unexpected error trying to create a cookie")
	}
	if cookie == nil {
		t.Fatalf("unexpected null instance for cookie")
	}
	if token == "" {
		t.Fatalf("unexpected empty value for token")
	}
}

func TestHandler_GetToken(t *testing.T) {
	// preparation
	handler := NewHandler(domainName, secretKey)
	cookie, originalToken, _ := handler.CreateCookie(&user, expirationTime)
	if originalToken == "" {
		t.Fatalf("unexpected empty value for token")
	}

	// real method call we want to test
	token, err := handler.GetToken(cookie)
	if err != nil {
		t.Fatalf("unexpected error trying to get token from cookie")
	}
	if originalToken != token {
		t.Fatalf("mismatching token obtained from cookie")
	}
}

func TestHandler_ReadCookie(t *testing.T) {
	// preparation
	handler := NewHandler(domainName, secretKey)
	originalCookie, _, _ := handler.CreateCookie(&user, expirationTime)

	// content doesn't matter here more than the cookies coming in the request
	req := httptest.NewRequest(http.MethodGet, "http://localhost/test", nil)
	req.AddCookie(originalCookie)

	// real method call we want to test
	cookie, token, err := handler.ReadCookie(req)
	if err != nil {
		t.Fatalf("unexpected error trying to read cookie from request")
	}

	parsedUser, _ := token.GetUser()
	if user.Username != parsedUser.Username ||
		user.Name != parsedUser.Name ||
		user.Lastname != parsedUser.Lastname {
		t.Fatalf("mismatching user associated to token")
	}

	if originalCookie.Value != cookie.Value {
		t.Fatalf("mismatching cookie obtained from request")
	}
}
