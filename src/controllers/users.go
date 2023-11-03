package controllers

import (
	"encoding/json"
	"nak-auth/models"
	svc "nak-auth/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserController struct {
	uSvc svc.IUserService
}

func NewUserController(uService svc.IUserService) *UserController {
	return &UserController{uSvc: uService}
}

func (*UserController) Path() string {
	return "/users"
}

func (c *UserController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, dbErr := c.uSvc.GetAll()
		if dbErr != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
		return
	case "POST":
		var usrBody models.User
		err := json.NewDecoder(r.Body).Decode(&usrBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dbErr := c.uSvc.Create(usrBody)
		if dbErr != nil {
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
	uSvc svc.IUserService
}

func NewUserByIdController(uService svc.IUserService) *UserByIdController {
	return &UserByIdController{uSvc: uService}
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
			user, dbErr := c.uSvc.GetByID(marks)
			if dbErr != nil {
				http.Error(w, dbErr.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		return
	case "DELETE":
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if ok {
			marks, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "bad request: invalid ID provided", http.StatusBadRequest)
				return
			}
			if c.uSvc.Delete(marks) != nil {
				http.Error(w, "failed to delete", http.StatusInternalServerError)
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
