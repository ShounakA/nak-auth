package templates

import (
	"html/template"
	"io"
)

type LoginPageData struct {
	Title    string
	Redirect string
	Issuer   string
}

func WriteLoginPage(w io.Writer, data LoginPageData) error {
	data.Title = "Login"
	loginPage := template.Must(template.ParseFiles("templates/pages/login.html", "templates/components/loginForm.html"))
	return loginPage.Execute(w, data)
}

func WriteHomePage(w io.Writer) error {
	homePage := template.Must(template.ParseFiles("templates/pages/home.html"))
	return homePage.Execute(w, nil)
}
