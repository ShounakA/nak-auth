package controllers

import (
	"encoding/json"
	"net/http"
)

type HealthController struct{}

type Health struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (*HealthController) Path() string {
	return "/health"
}

func (*HealthController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		health := Health{
			Status:  "healthy",
			Version: "0.1.0",
		}
		json.NewEncoder(w).Encode(health)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}
