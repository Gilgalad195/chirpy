package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() failed: %v", err)
	}
	if token == "" {
		t.Errorf("MakeJWT() returned empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	token, _ := MakeJWT(userID, secret, time.Hour)

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() failed with valid token: %v", err)
	}
	if gotID != userID {
		t.Errorf("ValidateJWT() got userID = %v, want %v", gotID, userID)
	}

	_, err = ValidateJWT(token, "wrong_secret")
	if err == nil {
		t.Error("ValidateJWT() should fail with wrong secret")
	}

	expiredToken, _ := MakeJWT(userID, secret, -time.Hour)
	_, err = ValidateJWT(expiredToken, secret)
	if err == nil {
		t.Error("ValidateJWT() shoud fail with expired token")
	}
}

func TestGetBearerToken(t *testing.T) {
	req, err := http.NewRequest("GET", "https://api.example.com/data", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer token123")

	token, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("failed to get token: %v", err)
	}

	if token != "token123" {
		t.Errorf("expected token %q, got %q", "token123", token)
	}
	fmt.Printf("token retrieved: %v", token)

}
