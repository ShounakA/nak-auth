package pages

import (
	"nak-auth/templates"
	"net/http"
)

type HomePage struct{}

func NewHomePage() *HomePage {
	return &HomePage{}
}

func (*HomePage) Path() string {
	return "/"
}

func (l *HomePage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	templates.WriteHomePage(w)
}
