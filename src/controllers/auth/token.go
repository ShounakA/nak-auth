package auth

import (
	"encoding/base64"
	"encoding/json"
	"nak-auth/services"
	"net/http"
	"strings"
)

type TokenController struct {
	user_svc   services.IUserService
	client_svc services.IClientService
	token_svc  services.ITokenService
}

type AccessBody struct {
	GrantType         string `json:"grant_type"`
	AuthorizationCode string `json:"authorization_code",omitempty`
	CodeVerifier      string `json:"code_verifier",omitempty`
	RedirectUri       string `json:"redirect_uri",omitempty`
	RefreshToken      string `json:"refresh_token",omitempty`
	//Scope			 string `json:"scope",omitempty`
}

func NewTokenController(client_service services.IClientService, token_service services.ITokenService) *TokenController {
	return &TokenController{client_svc: client_service, token_svc: token_service}
}

func (*TokenController) Path() string {
	return "/oauth/token"
}

func (l *TokenController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		clientId, clientSecret, err := parseClientCredentials(auth)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !l.client_svc.Authenticated(clientId, clientSecret) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accessBody := AccessBody{
			GrantType:         r.Form.Get("grant_type"),
			AuthorizationCode: r.Form.Get("authorization_code"),
			RedirectUri:       r.Form.Get("redirect_uri"),
			RefreshToken:      r.Form.Get("refresh_token"),
			CodeVerifier:      r.Form.Get("code_verifier"),
			// TODO Scope
		}
		client, err := l.client_svc.GetByID(clientId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate the redirect uri requested is allowed for the client
		if client.RedirectURI != accessBody.RedirectUri {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the grant type requested is allowed for the client
		if client.GrantType != accessBody.GrantType {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// TODO Validate Scope

		switch accessBody.GrantType {
		case "authorization_code":
			if accessBody.CodeVerifier == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			user, dberr := l.user_svc.VerifyAuthorizationCode(accessBody.AuthorizationCode, accessBody.CodeVerifier, clientId)
			if dberr != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			accessToken, err := l.token_svc.CreateAccessTokenWithAuthorization(clientId, clientSecret, user.Name, accessBody.AuthorizationCode)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(accessToken)
			break
		case "refresh_token":
			token, err := l.token_svc.CreateAccessTokenFromRefreshToken(clientId, clientSecret, accessBody.RefreshToken)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			json.NewEncoder(w).Encode(token)
			break
		case "client_credentials":
			accessToken, err := l.token_svc.CreateAccessToken(clientId, clientSecret, "service")
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

func parseClientCredentials(authHeader string) (string, string, error) {
	// Decode the username and password from the Authorization header
	decoded, err := base64.StdEncoding.DecodeString(authHeader[len("Basic "):])
	if err != nil {
		return "", "", err
	}
	credentials := strings.Split(string(decoded), ":")
	return credentials[0], credentials[1], nil
}
