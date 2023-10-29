package controllers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type CounterController struct {
	db *gorm.DB
}
type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by Counter to `Counter`
func (Counter) TableName() string {
	return "Counter"
}

type Counter struct {
	Id      int    `json:"id" sql:"id"`
	Name    string `json:"name" sql:"name"`
	Message string `json:"message" sql:"message"`
	Value   string `json:"value" sql:"value"`
}

func NewCounterController(db *gorm.DB) *CounterController {
	return &CounterController{db: db}
}

func (*CounterController) Path() string {
	return "/counter"
}

func (c *CounterController) WriteResponse(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var ctrs []Counter
		result := c.db.Find(&ctrs)
		if result.Error != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(ctrs)
		return
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}
