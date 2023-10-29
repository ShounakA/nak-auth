package pages

import (
	"nak-auth/services"
	"nak-auth/templates"
	"net/http"
)

type HomeController struct {
	uSvc services.ILoginService
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (*HomeController) Path() string {
	return "/"
}

func (l *HomeController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.WriteHomePage(w)
	}
}
