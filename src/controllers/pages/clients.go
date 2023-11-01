package pages

import (
	"nak-auth/services"
	"nak-auth/templates"
	"net/http"
)

type ClientsPage struct {
	cs services.IClientService
}

func NewClientsPage(clientSvc services.IClientService) *ClientsPage {
	return &ClientsPage{cs: clientSvc}
}

func (*ClientsPage) Path() string {
	return "/clients"
}

func (l *ClientsPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clients, error := l.cs.GetAll()
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	templates.WriteClientsPage(w, clients)
}
