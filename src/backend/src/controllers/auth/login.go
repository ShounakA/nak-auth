package auth

import (
	"fmt"
	"nak-auth/services"
	"net/http"
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
		r.ParseForm()
		username := r.Form.Get("username")
		secret := r.Form.Get("secret")
		badRequest := username == "" && secret == ""
		if badRequest {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		loginBody := services.Login{Name: username, Secret: secret}
		succ, user_id := l.login_svc.Login(loginBody.Name, loginBody.Secret)
		if !succ {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			client_id := r.URL.Query().Get("client_id")
			err := l.user_svc.AddAuthorizedClient(user_id, client_id)
			if badRequest {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			code, err = l.user_svc.AddAuthorizationCode(user_id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			redirect_uri := r.URL.Query().Get("redirect_uri")
			http.Redirect(w, r, fmt.Sprintf("%s?code=%s", redirect_uri, code), http.StatusSeeOther)
		}
	case "GET":
		p := r.URL.Path
		if p == "/login" {
			p = "static/index.html"
		}
		http.ServeFile(w, r, p)
	}
}

type AssetController struct {
	uSvc services.ILoginService
}

func NewAssetController(uService services.ILoginService) *AssetController {
	return &AssetController{uSvc: uService}
}

func (*AssetController) Path() string {
	return "/assets/{file}"
}

func (l *AssetController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		p := "static" + r.URL.Path
		http.ServeFile(w, r, p)
	}
}
