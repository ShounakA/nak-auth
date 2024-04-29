package services

import (
	"crypto/sha256"
	"encoding/base64"
	"nak-auth/db"
	"nak-auth/models"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

type LoginService struct {
	db              *gorm.DB
	fact            db.ILibSqlClientFactory
	tkn_svc         ITokenService
	nakAuthClientId string
	nakAuthSecret   string
	sessionName     string
	store           *sessions.CookieStore
}

type ILoginService interface {
	AuthenticateUser(username string, secret string) (bool, int, AccessToken, error)
	AuthenticateUserWithIdToken(accessToken string) (bool, int, AccessToken, error)
	ClientIsAuthenticated(r *http.Request) bool
	SaveSession(w http.ResponseWriter, r *http.Request, token AccessToken) error
}

// Creates a new login service. This services is responsible for authenticating users and creating sessions.
// It is also responsible for verifying that a client is authenticated.
// param db: The database connection
// param tkn_svc: The token service
func NewLoginService(fact db.ILibSqlClientFactory, tkn_svc ITokenService) *LoginService {
	sessionKey := os.Getenv("TOKEN_SIGNING_KEY")
	nakAuthClientId := os.Getenv("NAK_AUTH_CLIENT_ID")
	nakAuthSecret := os.Getenv("NAK_AUTH_CLIENT_SECRET")

	sessionName := "nak-auth-session"
	store := sessions.NewCookieStore([]byte(sessionKey))
	return &LoginService{fact: fact, sessionName: sessionName, store: store, tkn_svc: tkn_svc, nakAuthClientId: nakAuthClientId, nakAuthSecret: nakAuthSecret}
}

func (ls *LoginService) AuthenticateUser(username, secret string) (bool, int, AccessToken, error) {
	var user models.User
	var token AccessToken
	var userId int = -1
	var err error

	dbConn, err := ls.fact.CreateClient()
	if err != nil {
		return false, -1, AccessToken{}, err
	}
	success := false
	h := sha256.New()
	h.Write([]byte(secret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
	result := dbConn.Model(&models.User{}).First(&user, models.User{Name: username, Secret: hashSecret})
	if result.Error != nil {
		success = false
	} else {
		success = true
		token, err = ls.tkn_svc.CreateAccessToken(ls.nakAuthClientId, ls.nakAuthSecret, user.Name)
		userId = user.ID
	}
	return success, userId, token, err
}

func (ls *LoginService) AuthenticateUserWithIdToken(accessToken string) (bool, int, AccessToken, error) {
	var user models.User
	var token AccessToken
	var userId int = -1
	var err error
	success := false

	claims, err := ls.tkn_svc.VerifyNakAuthIdToken(accessToken, ls.nakAuthClientId, ls.nakAuthSecret)
	if err != nil {
		return success, userId, token, err
	}
	username := claims.(jwt.MapClaims)["sub"].(string)
	result := ls.db.Model(&models.User{}).First(&user, models.User{Name: username})
	if result.Error != nil {
		success = false
	} else {
		success = true
		userId = user.ID
	}
	return success, userId, token, err
}

func (ls *LoginService) ClientIsAuthenticated(r *http.Request) bool {
	session, _ := ls.store.Get(r, ls.sessionName)
	tokenValue := session.Values["token"]
	var token string
	if tokenValue != nil {
		token = tokenValue.(string)
	} else {
		token = ""
	}
	_, err := ls.tkn_svc.VerifyNakAuthAccessToken(token)
	if err != nil {
		return false
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return false
	}
	return true
}

func (ls *LoginService) SaveSession(w http.ResponseWriter, r *http.Request, token AccessToken) error {
	session, _ := ls.store.Get(r, ls.sessionName)
	session.Values["authenticated"] = true
	session.Values["token"] = token.AccessToken
	session.Options.MaxAge = int(token.ExpiresIn)
	return session.Save(r, w)
}
