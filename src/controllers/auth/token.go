package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"nak-auth/services"
	"net/http"
)

type AccessController struct {
	user_svc   services.IUserService
	client_svc services.IClientService
	token_svc  services.ITokenService
}

type AccessBody struct {
	GrantType         string `json:"grant_type"`
	ClientId          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	AuthorizationCode string `json:"authorization_code",omitempty`
	CodeVerifier      string `json:"code_verifier",omitempty`
	RedirectUri       string `json:"redirect_uri",omitempty`
	RefreshToken      string `json:"refresh_token",omitempty`
	//Scope			 string `json:"scope",omitempty`
}

func NewAccessController(client_service services.IClientService, token_service services.ITokenService) *AccessController {
	return &AccessController{client_svc: client_service, token_svc: token_service}
}

func (*AccessController) Path() string {
	return "/oauth/token"
}

func (l *AccessController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var accessBody AccessBody
		err := json.NewDecoder(r.Body).Decode(&accessBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch accessBody.GrantType {
		case "authorization_code":
			user, dberr := l.user_svc.VerifyAuthorizationCode(accessBody.AuthorizationCode, accessBody.CodeVerifier, accessBody.ClientId)
			if dberr != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			accessToken, err := l.token_svc.CreateAccessTokenWithAuthorization(accessBody.ClientId, accessBody.ClientSecret, user.Name, accessBody.AuthorizationCode)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(accessToken)
			break
		case "refresh_token":
			token, err := l.token_svc.CreateAccessTokenFromRefreshToken(accessBody.ClientId, accessBody.ClientSecret, accessBody.RefreshToken)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			json.NewEncoder(w).Encode(token)
			break
		case "client_credentials":
			client, err := l.client_svc.GetByID(accessBody.ClientId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if client.GrantType != accessBody.GrantType {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			h := sha256.New()
			h.Write([]byte(accessBody.ClientSecret))
			hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
			if client.Secret != hashSecret {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			accessToken, err := l.token_svc.CreateAccessToken(accessBody.ClientId, accessBody.ClientSecret)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(accessToken)
			break
		default:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}
