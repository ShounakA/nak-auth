package models

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
