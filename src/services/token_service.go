package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // in seconds
	RefreshToken string `json:"refresh_token,omitempty"`
}

type ITokenService interface {
	CreateRefreshToken(clientId string) (string, error)
	CreateAccessToken(clientId string) (AccessToken, error)
	CreateAccessTokenWithAuthorization(clientId, userName, authorization_code string) (AccessToken, error)
	CreateAccessTokenFromRefreshToken(refreshToken string) (AccessToken, error)
}

type TokenService struct {
	tokenSigningKey string
}

func NewTokenService(db *gorm.DB) *TokenService {
	signKey := os.Getenv("TOKEN_SIGNING_KEY")
	//TODO sign by client secret not the key
	return &TokenService{tokenSigningKey: signKey}
}

func (s *TokenService) CreateAccessToken(clientId string) (AccessToken, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	expireAt := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = expireAt
	claims["iat"] = time.Now().Unix()
	claims["client_id"] = clientId

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.tokenSigningKey))
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expireAt - time.Now().Unix(),
	}, nil
}

func (s *TokenService) CreateAccessTokenWithAuthorization(clientId, userName, authorization_code string) (AccessToken, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	expireAt := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = expireAt
	claims["iat"] = time.Now().Unix()
	claims["client_id"] = clientId
	claims["user"] = userName
	claims["authorization_code"] = authorization_code

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.tokenSigningKey))
	if err != nil {
		return AccessToken{}, err
	}
	refreshToken, err := s.CreateRefreshToken(clientId)
	if err != nil {
		return AccessToken{}, err
	}
	return AccessToken{
		AccessToken:  tokenString,
		TokenType:    "Bearer",
		ExpiresIn:    expireAt - time.Now().Unix(),
		RefreshToken: refreshToken,
	}, nil
}

func (s *TokenService) CreateRefreshToken(clientId string) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	expireAt := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = expireAt
	claims["iat"] = time.Now().Unix()
	claims["client_id"] = clientId

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(s.tokenSigningKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *TokenService) CreateAccessTokenFromRefreshToken(refreshToken string) (AccessToken, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the signing key
		return []byte(s.tokenSigningKey), nil
	})
	if err != nil {
		return AccessToken{}, err
	}
	if !token.Valid {
		return AccessToken{}, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AccessToken{}, errors.New("invalid token claims")
	}
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		return AccessToken{}, errors.New("token has expired")
	}
	accessToken, err := s.CreateAccessToken(claims["client_id"].(string))
	if err != nil {
		return AccessToken{}, err
	}
	return accessToken, nil

}
