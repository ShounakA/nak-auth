package services

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type LoginService struct {
	db *gorm.DB
}

type ILoginService interface {
	Login(username string, secret string) (bool, int, AccessToken, error)
}

func NewLoginService(db *gorm.DB) *LoginService {
	return &LoginService{db: db}
}

type Login struct {
	Name   string `json:"username"`
	Secret string `json:"secret"`
}

func (ls *LoginService) Login(username string, secret string) (bool, int, AccessToken, error) {
	var user User
	var token AccessToken
	var userId int = -1
	var err error
	success := false
	h := sha256.New()
	h.Write([]byte(secret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
	result := ls.db.Model(&User{}).First(&user, User{Name: username, Secret: hashSecret})
	if result.Error != nil {
		success = false
	} else {
		success = true
		token, err = createToken(username)
		userId = user.ID
	}
	return success, userId, token, err
}

func createToken(username string) (AccessToken, error) {
	// Set the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the user ID and expiration time
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   username,
	}

	// Create the JWT token with the claims and the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte("my_secret_key")
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		AccessToken: signedToken,
		ExpiresIn:   claims.ExpiresAt - time.Now().Unix(),
		TokenType:   "Bearer",
	}, nil
}
