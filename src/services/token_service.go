package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // in seconds
	RefreshToken string `json:"refresh_token,omitempty"`
	IdToken      string `json:"id_token,omitempty"`
}

type ITokenService interface {
	CreateRefreshToken(clientId string) (string, error)
	CreateAccessToken(clientId, clientSecret, username string) (AccessToken, error)
	CreateAccessTokenWithAuthorization(clientId, clientSecret, userName, authorization_code string) (AccessToken, error)
	CreateAccessTokenFromRefreshToken(clientId, clientSecret, refreshToken string) (AccessToken, error)
	CreateIdToken(clientId, clientSecret, username string) (string, error)
	GenerateSecret(clientId string) string
	VerifyNakAuthAccessToken(token string) (jwt.Claims, error)
	VerifyNakAuthIdToken(id_token, clientId, clientSecret string) (jwt.Claims, error)
}

type TokenService struct {
	tokenSigningKey string
	nakAuthClientId string
	nakAuthSecret   string
}

// Create a new token service. This service is responsible for creating, verifying, and refreshing tokens.
func NewTokenService() *TokenService {
	signKey := os.Getenv("TOKEN_SIGNING_KEY")
	nakAuthSecret := os.Getenv("NAK_AUTH_CLIENT_SECRET")
	nakAuthClientId := os.Getenv("NAK_AUTH_CLIENT_ID")

	return &TokenService{tokenSigningKey: signKey, nakAuthSecret: nakAuthSecret, nakAuthClientId: nakAuthClientId}
}

func (s *TokenService) CreateAccessToken(clientId, clientSecret, username string) (AccessToken, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	expireAt := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = expireAt
	claims["iat"] = time.Now().Unix()
	claims["client_id"] = clientId
	claims["sub"] = username

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expireAt - time.Now().Unix(),
	}, nil
}

func (s *TokenService) CreateAccessTokenWithAuthorization(clientId, clientSecret, userName, authorization_code string) (AccessToken, error) {
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
	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return AccessToken{}, err
	}
	refreshToken, err := s.CreateRefreshToken(clientId)
	if err != nil {
		return AccessToken{}, err
	}
	idToken, err := s.CreateIdToken(clientId, clientSecret, userName)
	if err != nil {
		return AccessToken{}, err
	}
	return AccessToken{
		AccessToken:  tokenString,
		TokenType:    "Bearer",
		ExpiresIn:    expireAt - time.Now().Unix(),
		RefreshToken: refreshToken,
		IdToken:      idToken,
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

func (s *TokenService) CreateAccessTokenFromRefreshToken(clientId, clientSecret, refreshToken string) (AccessToken, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the signing key
		return []byte(clientSecret), nil
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
	if clientId != claims["client_id"].(string) {
		return AccessToken{}, errors.New("invalid client id")
	}
	accessToken, err := s.CreateAccessToken(claims["client_id"].(string), clientSecret, claims["sub"].(string))
	if err != nil {
		return AccessToken{}, err
	}
	return accessToken, nil
}

func (s *TokenService) GenerateSecret(clientId string) string {
	currentTime := time.Now().String()
	clientSecret := clientId + currentTime
	// Hash the password

	h := sha256.New()
	h.Write([]byte(clientSecret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashSecret
}

func (s *TokenService) CreateIdToken(clientId, clientSecret, username string) (string, error) {
	// Create a new JWT token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	expireAt := time.Now().Add(time.Hour * 5).Unix()
	claims["exp"] = expireAt
	claims["iat"] = time.Now().Unix()
	claims["client_id"] = clientId
	claims["sub"] = username
	claims["aud"] = clientId
	claims["typ"] = "id_token"

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// These should only be used for nak-auth to nak-auth communication

func (s *TokenService) GenerateNakAuthAccessToken(username string) (AccessToken, error) {
	return s.CreateAccessToken(s.nakAuthClientId, s.nakAuthSecret, username)
}

func (s *TokenService) VerifyNakAuthAccessToken(access_token string) (jwt.Claims, error) {
	token, err := jwt.Parse(access_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.nakAuthSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid access token")
}

func (s *TokenService) VerifyNakAuthIdToken(id_token, clientId, clientSecret string) (jwt.Claims, error) {
	token, err := jwt.Parse(id_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println(token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(clientSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["client_id"] != clientId {
			return nil, errors.New("invalid client id")
		}
		if claims["typ"] != "id_token" {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}
	return nil, fmt.Errorf("invalid id token")
}
