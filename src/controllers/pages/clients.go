package pages

import (
	"nak-auth/services"
	"nak-auth/templates"
	"net/http"
)

type ClientsPage struct {
	cs    services.IClientService
	login services.ILoginService
}

func NewClientsPage(clientSvc services.IClientService, loginSvc services.ILoginService) *ClientsPage {
	return &ClientsPage{cs: clientSvc, login: loginSvc}
}

func (*ClientsPage) Path() string {
	return "/clients"
}

func (l *ClientsPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !l.login.ClientIsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	clients, error := l.cs.GetAll()
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	templates.WriteClientsPage(w, clients)
}
