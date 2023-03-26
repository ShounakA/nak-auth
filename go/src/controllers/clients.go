package controllers

import (
	"encoding/json"
	srv "nak-auth/services"
	"net/http"

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
		client_json := srv.ListOfClientsToListOfClientJson(clients)
		json.NewEncoder(w).Encode(client_json)
		return
	case "POST":
		var usrBody srv.ClientJson
		err := json.NewDecoder(r.Body).Decode(&usrBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Create the user
		err = c.cs.Create(usrBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode("Created")
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
			json.NewEncoder(w).Encode("Deleted")
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}
