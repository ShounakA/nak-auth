package services

import (
	"crypto/sha256"
	"encoding/base64"
	"nak-auth/models"
	"testing"
)

func TestAuthenticateUser_Success(t *testing.T) {
	db := setupDB()

	service := LoginService{db: db}

	h := sha256.New()
	h.Write([]byte("Secret"))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
	registerUser := models.User{Name: "Test", Secret: hashSecret}

	db.Create(&registerUser)

	isAuthenticated, userId, token, err := service.AuthenticateUser("Test", "Secret")
	if err != nil {
		t.Errorf("Error was not expected while authenticating user: %s", err)
	}
	if !isAuthenticated {
		t.Errorf("Expected user to be authenticated")
	}
	if userId != 1 {
		t.Errorf("Expected user id to be %d, got %d", registerUser.ID, userId)
	}
	if token.AccessToken == "" {
		t.Errorf("Expected access token to be set")
	}
}

func TestAuthenticateUser_Failure(t *testing.T) {
	db := setupDB()

	service := LoginService{db: db}

	isAuthenticated, userId, token, err := service.AuthenticateUser("Test", "Secret")
	if err != nil {
		t.Errorf("Error was not expected while authenticating user: %s", err)
	}
	if isAuthenticated {
		t.Errorf("Expected user to not be authenticated")
	}
	if userId != -1 {
		t.Errorf("Expected user id to be %d, got %d", -1, userId)
	}
	if token.AccessToken != "" {
		t.Errorf("Expected access token to be set")
	}
}
