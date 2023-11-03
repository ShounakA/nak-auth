package models

import (
	"testing"
)

func TestFrom_Scope(t *testing.T) {
	scope := Scope{Name: "testScope"}

	if scope.From() != "testScope" {
		t.Errorf("Expected 'testScope', got: %s", scope.From())
	}
}

func TestCreateTable_Scope(t *testing.T) {
	scope := Scope{}

	if scope.CreateTable() != "scopes" {
		t.Errorf("Expected 'scopes', got: %s", scope.CreateTable())
	}
}

func TestEquals_Scope(t *testing.T) {
	scope1 := Scope{Name: "testScope"}
	scope2 := Scope{Name: "testScope"}
	scope3 := Scope{Name: "differentScope"}

	if !scope1.Equals(scope2) {
		t.Errorf("Expected scopes to be equal")
	}

	if scope1.Equals(scope3) {
		t.Errorf("Expected scopes to be not equal")
	}

	if scope1.Equals(nil) {
		t.Errorf("Expected scope not to be equal to nil")
	}
}
