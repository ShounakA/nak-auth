package pages

import (
	"nak-auth/services"
	"nak-auth/templates"
	"net/http"
)

type HomePage struct {
	login_svc services.ILoginService
}

func NewHomePage(login_service services.ILoginService) *HomePage {
	return &HomePage{
		login_svc: login_service,
	}
}

func (*HomePage) Path() string {
	return "/"
}

func (l *HomePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !l.login_svc.ClientIsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.WriteHomePage(w)
}
