package auth

import (
	"encoding/base64"
	"fmt"
	"nak-auth/services"
	"nak-auth/templates"
	"net/http"
	"strings"
)

type LoginController struct {
	login_svc services.ILoginService
	user_svc  services.IUserService
}

func NewLoginController(login_service services.ILoginService, user_service services.IUserService) *LoginController {
	return &LoginController{login_svc: login_service, user_svc: user_service}
}

func (*LoginController) Path() string {
	return "/login"
}

func (l *LoginController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var code string = ""
		auth := r.Header.Get("Authorization")

		if auth == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode the username and password from the Authorization header
		decoded, err := base64.StdEncoding.DecodeString(auth[len("Basic "):])
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		credentials := strings.Split(string(decoded), ":")
		username := credentials[0]
		password := credentials[1]

		// Authenticate the user with the username and password
		succ, userId, token, err := l.login_svc.Login(username, password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !succ {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		mode := r.URL.Query().Get("mode")
		redirect_uri := r.URL.Query().Get("redirect_uri")
		requestRedirect := fmt.Sprintf("%s?access_token=%s", redirect_uri, token.AccessToken)
		if mode == "authorize" {
			client_id := r.URL.Query().Get("client_id")
			challenge := r.URL.Query().Get("code_challenge")
			err = l.user_svc.AddAuthorizedClient(userId, client_id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			code, err = l.user_svc.AddAuthorizationCode(userId, challenge)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			requestRedirect = fmt.Sprintf("%s?code=%s", redirect_uri, code)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.Redirect(w, r, requestRedirect, http.StatusSeeOther)
	case "GET":
		redirect_uri := r.URL.Query().Get("redirect_uri")
		issuer := r.URL.Query().Get("issuer")
		if redirect_uri == "" {
			redirect_uri = "http://localhost:8080"
		}
		if issuer == "" {
			issuer = "nak-auth"
		}
		templates.WriteLoginPage(w, templates.LoginPageData{Redirect: redirect_uri, Issuer: issuer})
	}
}
