package models

import "testing"

func setup() Client {
	scopes := []Scope{}
	scopes = append(scopes, Scope{Name: "testScope1"})
	scopes = append(scopes, Scope{Name: "testScope2"})
	expectedClient := Client{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI", Scopes: scopes}
	return expectedClient
}

func TestCreateTable_Client(t *testing.T) {
	actualClient := Client{}

	if actualClient.CreateTable() != "clients" {
		t.Errorf("Expected 'scopes', got: %s", actualClient.CreateTable())
	}
}

func TestEquals_Client(t *testing.T) {
	scopes := []Scope{}
	scopes = append(scopes, Scope{Name: "testScope1"})
	scopes = append(scopes, Scope{Name: "testScope2"})
	actualClient := Client{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI", Scopes: scopes}

	if !actualClient.Equals(actualClient) {
		t.Errorf("Expected clients to be equal")
	}
}
