package models

import "strings"

type Scope struct {
	Name string `json:"name" sql:"name" gorm:"primaryKey"`
}

// # Convert from Scope to string
func (scope Scope) From() string {
	return scope.Name
}

// # Creates 'scopes' table
func (Scope) CreateTable() string {
	return "scopes"
}

// # Checks if two scopes are equal
func (scope Scope) Equals(other interface{}) bool {
	if other == nil {
		return false
	}
	return scope.Name == other.(Scope).Name
}

func ParseScopes(scope string) []Scope {
	scopes := []Scope{}
	for _, s := range strings.Split(scope, " ") {
		scopes = append(scopes, Scope{Name: s})
	}
	return scopes
}
