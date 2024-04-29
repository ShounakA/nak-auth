package services

import (
	"nak-auth/models"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockTokenService struct {
	tokenSigningKey string
	nakAuthClientId string
	nakAuthSecret   string
}

type MockLibSqlFactory struct {
	db *gorm.DB
}

func (l *MockLibSqlFactory) CreateClient() (*gorm.DB, error) {
	return l.db, nil
}

func (s *MockTokenService) GenerateSecret(clientId string) string {
	return "Secret"
}

func (s *MockTokenService) CreateAccessToken(clientId, clientSecret, username string) (AccessToken, error) {
	return AccessToken{AccessToken: "totally_legit_token"}, nil
}

func (s *MockTokenService) CreateAccessTokenWithAuthorization(clientId, clientSecret, userName, authorization_code string) (AccessToken, error) {
	return AccessToken{}, nil
}

func (s *MockTokenService) CreateRefreshToken(clientId string) (string, error) {
	return "totally_legit_token", nil
}

func (s *MockTokenService) CreateAccessTokenFromRefreshToken(clientId, clientSecret, refreshToken string) (AccessToken, error) {
	return AccessToken{}, nil
}

func (s *MockTokenService) VerifyNakAuthAccessToken(token string) (jwt.Claims, error) {
	return nil, nil
}

func (s *MockTokenService) VerifyNakAuthIdToken(id_token, clientId, clientSecret string) (jwt.Claims, error) {
	return nil, nil
}

func (s *MockTokenService) CreateIdToken(clientId, clientSecret, username string) (string, error) {
	return "", nil
}

func setupDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Client{}, &models.User{})

	return db
}

func TestGetAll_Success(t *testing.T) {
	db := setupDB()
	service := ClientService{fact: &MockLibSqlFactory{db: db}}

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

func TestGetAll_Empty(t *testing.T) {
	db := setupDB()
	service := ClientService{fact: &MockLibSqlFactory{db: db}}

	clients, err := service.GetAll()
	if err != nil {
		t.Errorf("Error was not expected while getting all clients: %s", err)
	}

	if len(clients) != 0 {
		t.Errorf("Expected no clients, got: %d", len(clients))
	}
}

func TestGetbyID_Success(t *testing.T) {
	db := setupDB()
	service := ClientService{fact: &MockLibSqlFactory{db: db}}
	scopes := []models.Scope{}
	scopes = append(scopes, models.Scope{Name: "testScope1"})
	scopes = append(scopes, models.Scope{Name: "testScope2"})
	expectedClient := models.Client{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI", Scopes: scopes}

	scopes2 := []models.Scope{}
	scopes2 = append(scopes2, models.Scope{Name: "testScope1"})
	scopes2 = append(scopes2, models.Scope{Name: "testScope2"})

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

func TestGetbyID_NotFound(t *testing.T) {
	db := setupDB()
	service := ClientService{fact: &MockLibSqlFactory{db: db}}

	_, err := service.GetByID("Test")
	if err == nil {
		t.Errorf("Error was expected while getting client by id")
	}
}

func TestCreate_Success(t *testing.T) {
	db := setupDB()
	service := ClientService{fact: &MockLibSqlFactory{db: db}, token_svc: &MockTokenService{tokenSigningKey: "Secret", nakAuthSecret: "Secret", nakAuthClientId: "Secret"}}

	scopes := []models.Scope{}
	scopes = append(scopes, models.Scope{Name: "testScope1"})
	scopes = append(scopes, models.Scope{Name: "testScope2"})
	expectedClient := models.Client{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI", Scopes: scopes}
	newClient := models.ClientJson{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI", Scopes: []string{"testScope1", "testScope2"}}

	err := service.Create(newClient)
	if err != nil {
		t.Errorf("Error was not expected while creating client: %s", err)
	}

	var clients []models.Client
	result := db.Preload("Scopes").Find(&clients)
	if result.Error != nil {
		t.Errorf("Error was not expected while getting all clients: %s", result.Error)
	}
	if len(clients) != 1 {
		t.Errorf("Expected one client, got: %d", len(clients))
	}
	actualClient := clients[0]

	if !actualClient.Equals(expectedClient) {
		t.Errorf("Expected client to be %s, got: %s", expectedClient, actualClient)
	}
}

func TestDelete_Success(t *testing.T) {
	db := setupDB()
	service := ClientService{fact: &MockLibSqlFactory{db: db}}
	scopes := []models.Scope{}
	scopes = append(scopes, models.Scope{Name: "testScope1"})
	scopes = append(scopes, models.Scope{Name: "testScope2"})
	expectedClient := models.Client{Name: "Test", Secret: "Secret", GrantType: "GrantType", RedirectURI: "RedirectURI", Scopes: scopes}

	db.Create(&expectedClient)

	err := service.Delete("Test")
	if err != nil {
		t.Errorf("Error was not expected while deleting client: %s", err)
	}
}
