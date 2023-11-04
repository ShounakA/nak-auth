package services

import (
	"crypto/sha256"
	"encoding/base64"
	"nak-auth/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type LoginService struct {
	db          *gorm.DB
	sessionName string
	store       *sessions.CookieStore
}

type ILoginService interface {
	AuthenticateUser(username string, secret string) (bool, int, AccessToken, error)
	ClientIsAuthenticated(r *http.Request) bool
	SaveSession(w http.ResponseWriter, r *http.Request) error
}

func NewLoginService(db *gorm.DB) *LoginService {
	sessionKey := os.Getenv("TOKEN_SIGNING_KEY")
	sessionName := "nak-auth-session"
	store := sessions.NewCookieStore([]byte(sessionKey))
	return &LoginService{db: db, sessionName: sessionName, store: store}
}

func (ls *LoginService) AuthenticateUser(username string, secret string) (bool, int, AccessToken, error) {
	var user models.User
	var token AccessToken
	var userId int = -1
	var err error
	success := false
	h := sha256.New()
	h.Write([]byte(secret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
	result := ls.db.Model(&models.User{}).First(&user, models.User{Name: username, Secret: hashSecret})
	if result.Error != nil {
		success = false
	} else {
		success = true
		token, err = createToken(username)
		userId = user.ID
	}
	return success, userId, token, err
}

func (ls *LoginService) ClientIsAuthenticated(r *http.Request) bool {
	session, _ := ls.store.Get(r, ls.sessionName)
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
}

func (ls *LoginService) SaveSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := ls.store.Get(r, ls.sessionName)
	session.Values["authenticated"] = true
	return session.Save(r, w)
}

// TODO move to Token Service
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
