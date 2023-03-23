package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

type User struct {
	ID     int    `json:"id" sql:"id"`
	Name   string `json:"name" sql:"name"`
	Secret string `json:"secret" sql:"secret"`
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (*UserController) Path() string {
	return "/users"
}

func (c *UserController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var users []User
		result := c.db.Find(&users)
		if result.Error != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
		return
	case "POST":
		var usrBody User
		err := json.NewDecoder(r.Body).Decode(&usrBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Hash the password
		h := sha256.New()
		h.Write([]byte(usrBody.Secret))
		hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
		newUser := User{
			Name:   usrBody.Name,
			Secret: hashSecret,
		}

		// Write to the pscale
		result := c.db.Create(&newUser)
		if result.Error != nil {
			http.Error(w, "internal server error: failed to create user", http.StatusInternalServerError)
			return
		}

		// Respond with created
		json.NewEncoder(w).Encode("Created")
		return
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

type UserByIdController struct {
	db *gorm.DB
}

func NewUserByIdController(db *gorm.DB) *UserByIdController {
	return &UserByIdController{db: db}
}

func (*UserByIdController) Path() string {
	return "/users/{id}"
}

func (c *UserByIdController) WriteResponse(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if ok {
			marks, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "bad request: invalid ID provided", http.StatusBadRequest)
				return
			}
			var user User
			result := c.db.First(&user, User{ID: marks})
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		return
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}
