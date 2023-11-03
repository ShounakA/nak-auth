package services

import (
	"nak-auth/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Client{})

	return db
}

func TestGetAll(t *testing.T) {
	db := setupDB()
	service := ClientService{db: db}

	// Setup test data
	db.Create(&models.Client{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI"})

	clients, err := service.GetAll()
	if err != nil {
		t.Errorf("Error was not expected while getting all clients: %s", err)
	}

	if len(clients) != 1 {
		t.Errorf("Expected one client, got: %d", len(clients))
	}

	if clients[0].Name != "Test" {
		t.Errorf("Expected name to be 'Test', got: %s", clients[0].Name)
	}
}

func TestGetbyID_Success(t *testing.T) {
	db := setupDB()
	service := ClientService{db: db}
	scopes := []models.Scope{}
	scopes = append(scopes, models.Scope{Name: "testScope1"})
	scopes = append(scopes, models.Scope{Name: "testScope2"})

	scopes2 := []models.Scope{}
	scopes2 = append(scopes2, models.Scope{Name: "testScope1"})
	scopes2 = append(scopes2, models.Scope{Name: "testScope2"})
	expectedClient := models.Client{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI", Scopes: scopes}

	db.Create(&expectedClient)
	db.Create(&models.Client{Name: "Test2", Secret: "Secret2", GrantType: "GrantType2", RedirectURI: "RedirectURI2", Scopes: scopes2})

	client, err := service.GetByID("Test")
	if err != nil {
		t.Errorf("Error was not expected while getting client by id: %s", err)
	}
	if !client.Equals(expectedClient) {
		t.Errorf("Expected client to be 'Test', got: %s", client.Name)
	}

}
