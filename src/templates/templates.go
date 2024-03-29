package templates

import (
	"html/template"
	"io"
	"nak-auth/models"
)

type LoginPageData struct {
	Title    string
	Redirect string
	Issuer   string
}

type ClientsPageData struct {
	Title   string
	Clients []models.Client
}

func WriteLoginPage(w io.Writer, data LoginPageData) error {
	data.Title = "Login"
	loginPage := template.Must(template.ParseFiles("templates/pages/login.html", "templates/components.html"))
	return loginPage.Execute(w, data)
}

func WriteHomePage(w io.Writer) error {
	homePage := template.Must(template.ParseFiles("templates/pages/home.html"))
	return homePage.Execute(w, nil)
}

func WriteClientsPage(w io.Writer, clients []models.Client) error {
	clientsPage := template.Must(template.ParseFiles("templates/pages/clients.html", "templates/components.html"))
	data := ClientsPageData{Title: "Clients", Clients: clients}
	return clientsPage.Execute(w, data)
}

func WriteClientsFragment(w io.Writer, clients []models.Client) error {
	clientsPage := template.Must(template.ParseFiles("templates/components.html"))
	return clientsPage.ExecuteTemplate(w, "clientRange", ClientsPageData{Title: "Clients", Clients: clients})
}

func WriteClientFragment(w io.Writer) error {
	clientsPage := template.Must(template.ParseFiles("templates/components.html"))
	return clientsPage.ExecuteTemplate(w, "clientList", nil)
}
