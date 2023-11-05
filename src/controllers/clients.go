package controllers

import (
	"encoding/json"
	"nak-auth/models"
	srv "nak-auth/services"
	"nak-auth/templates"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ClientController struct {
	cs srv.IClientService
}

func NewClientController(clientService srv.IClientService) *ClientController {
	return &ClientController{cs: clientService}
}

func (*ClientController) Path() string {
	return "/clients"
}

func (c *ClientController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		clients, dbErr := c.cs.GetAll()
		if dbErr != nil {
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}
		client_json := models.ListOfClientsToListOfClientJson(clients)
		if r.Header.Get("Hx-Request") == "true" {
			templates.WriteClientsFragment(w, clients)
		} else {
			json.NewEncoder(w).Encode(client_json)
		}
		return
	case "POST":
		var usrBody models.ClientJson
		if (r.Header.Get("Content-Type") == "application/json") && (r.Header.Get("Content-Type") == "application/json; charset=utf-8") {
			err := json.NewDecoder(r.Body).Decode(&usrBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else if (r.Header.Get("Content-Type") == "application/x-www-form-urlencoded") || (r.Header.Get("Content-Type") == "application/x-www-form-urlencoded; charset=utf-8") {
			usrBody.Name = r.FormValue("name")
			usrBody.GrantType = r.FormValue("grant_type")
			usrBody.Secret = r.FormValue("secret")
			usrBody.RedirectURI = r.FormValue("redirect_uri")
			usrBody.Scopes = r.Form["scope"]
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(usrBody.Name) == "" {
			http.Error(w, "Name is required.", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(usrBody.GrantType) == "" {
			http.Error(w, "Grant Type is required.", http.StatusBadRequest)
			return
		}

		//Create the user
		err := c.cs.Create(usrBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if r.Header.Get("Hx-Request") == "true" {
			templates.WriteClientFragment(w)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
		return
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}

}

type ClientByIdController struct {
	cs srv.IClientService
}

func NewClientByIdController(cs srv.IClientService) *ClientByIdController {
	return &ClientByIdController{cs: cs}
}

func (*ClientByIdController) Path() string {
	return "/clients/{id}"
}

func (c *ClientByIdController) WriteResponse(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if ok {
			client, err := c.cs.GetByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			client_json := client.From()
			json.NewEncoder(w).Encode(client_json)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		return
	case "DELETE":
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if ok {
			err := c.cs.Delete(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if r.Header.Get("Hx-Request") == "true" {
				templates.WriteClientFragment(w)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}
