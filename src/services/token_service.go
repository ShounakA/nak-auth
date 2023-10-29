package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // in seconds
}

type ITokenService interface {
	// CreateBasicToken(userId uint) (string, error)
	// CreateRefreshToken(userId uint) (string, error)
	CreateAccessToken(clientId, clientSecret string) (AccessToken, error)
	// ValidateAccessToken(token string) (bool, error)
}

type TokenService struct {
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{}
}

func (s *TokenService) CreateAccessToken(clientId, clientSecret string) (AccessToken, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	expireAt := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = expireAt
	claims["iat"] = time.Now().Unix()

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte("mySecretKey"))
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expireAt - time.Now().Unix(),
	}, nil
}
