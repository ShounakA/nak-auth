package pages

import (
	"nak-auth/templates"
	"net/http"
)

type LoginPage struct{}

func NewLoginPage() *LoginPage {
	return &LoginPage{}
}

func (*LoginPage) Path() string {
	return "/login"
}

func (l *LoginPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
