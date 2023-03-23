package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ClientController struct {
	db *gorm.DB
}

type Client struct {
	ID   		string 		`json:"name" sql:"name" gorm:"primaryKey"`
	Secret 		string 		`json:"secret" sql:"secret"`
//	AccessToken []string 	`json:"access_tokens" sql:"access_tokens"`  
}

func (Client) CreateTable() string {
	return "clients"
}

func NewClientController(db *gorm.DB) *ClientController {
	return &ClientController{db: db}
}

func (*ClientController) Path() string {
	return "/clients"
}

func (c *ClientController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var clients []Client
		result := c.db.Find(&clients)
		if result.Error != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(clients)
		return
	case "POST":
		var usrBody Client
		err := json.NewDecoder(r.Body).Decode(&usrBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Hash the password
		h := sha256.New()
		h.Write([]byte(usrBody.Secret))
		hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
		newClient := Client{
			ID:   usrBody.ID,
			Secret: hashSecret,
		}

		// Write to the pscale
		result := c.db.Create(&newClient)
		if result.Error != nil {
			http.Error(w, "internal server error: failed to create client", http.StatusInternalServerError)
			return
		}

		// Respond with created
		json.NewEncoder(w).Encode("Created")
		return
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

type ClientByIdController struct {
	db *gorm.DB
}

func NewClientByIdController(db *gorm.DB) *ClientByIdController {
	return &ClientByIdController{db: db}
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
			var client Client
			result := c.db.First(&client, Client{ID: id})
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(client)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		return
	case "DELETE":
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if ok {
			var client Client
			result := c.db.Delete(&client, Client{ID: id})
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
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
