package auth

import (
	"monkiato/guardian/internal/models"
	"testing"
	"time"
)

var (
	secretKey      = "my-secret-key"
	expirationTime = time.Now().Add(time.Second * 5).Unix()
	user           = models.SessionUser{Name: "user", Lastname: "test", Username: "user_test"}
)

func TestNewJwtToken(t *testing.T) {
	token := NewJwtToken(&user, secretKey, expirationTime)
	if token == nil {
		t.Fatalf("unexpected null instance")
	}

	if token.secretKey != secretKey {
		t.Errorf("secretKey doesn't match")
	}

	if token.Token == nil {
		t.Fatalf("unexpected null token")
	}
}

func TestJwtToken_ToString(t *testing.T) {
	token := NewJwtToken(&user, secretKey, expirationTime)
	strToken, err := token.ToString()
	if err != nil {
		t.Fatalf("unexpected error during token string conversion")
	}
	if strToken == "" {
		t.Fatalf("unexpected empty string token")
	}
}

func TestNewJwtTokenFromString(t *testing.T) {
	// create from user first to obtain the string token
	originToken := NewJwtToken(&user, secretKey, expirationTime)
	strToken, _ := originToken.ToString()

	token, err := NewJwtTokenFromString(strToken, secretKey)
	if err != nil {
		t.Fatalf("unexpected error during JwtToken creation from an existing string token")
	}

	if token == nil {
		t.Fatalf("unexpected null instance")
	}

	if token.secretKey != secretKey {
		t.Errorf("secretKey doesn't match")
	}

	if token.Token == nil {
		t.Fatalf("unexpected null token")
	}
}

func TestNewJwtTokenFromStringError(t *testing.T) {
	token, err := NewJwtTokenFromString("fake-token", secretKey)
	if err == nil {
		t.Fatalf("unexpected null error using fake string token")
	}

	if token != nil {
		t.Fatalf("unexpected valid token obtained from fake token")
	}
}

func TestNewJwtTokenFromStringSignatureError(t *testing.T) {
	// create from user first to obtain the string token
	originToken := NewJwtToken(&user, secretKey, expirationTime)
	strToken, _ := originToken.ToString()

	token, err := NewJwtTokenFromString(strToken, "wrong-secret-key")
	if err == nil {
		t.Fatalf("unexpected null error using wrong secret key")
	}

	if token != nil {
		t.Fatalf("unexpected valid token obtained from wrong secret key")
	}
}
