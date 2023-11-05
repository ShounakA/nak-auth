package services

import (
	"nak-auth/models"
	"os"

	"gorm.io/gorm"
)

type IClientService interface {
	GetAll() ([]models.Client, error)
	GetByID(id string) (models.Client, error)
	Create(newClient models.ClientJson) error
	Delete(id string) error
	Authenticated(clientId string, clientSecret string) bool
}

type ClientService struct {
	db        *gorm.DB
	clientId  string
	token_svc ITokenService
}

// Create a new client service. This service is responsible for managing clients.
func NewClientService(db *gorm.DB, tkn_svc ITokenService) *ClientService {
	name := os.Getenv("NAK_AUTH_CLIENT_ID")
	secret := os.Getenv("NAK_AUTH_CLIENT_SECRET")
	db.Create(&models.Client{Name: name, Secret: secret, GrantType: "client_credentials", RedirectURI: "http://localhost:8080", Scopes: []models.Scope{{Name: "internal"}}})
	return &ClientService{db: db, token_svc: tkn_svc, clientId: name}
}

func (cs *ClientService) GetAll() ([]models.Client, error) {
	var clients []models.Client
	var dbError error = nil
	result := cs.db.Model(&models.Client{}).Preload("Scopes").Where("name != ?", cs.clientId).Find(&clients)
	if result.Error != nil {
		dbError = result.Error
	}
	return clients, dbError
}

func (cs *ClientService) GetByID(id string) (models.Client, error) {
	var client models.Client
	var dbErr error = nil
	result := cs.db.Model(&models.Client{}).Preload("Scopes").First(&client, models.Client{Name: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return client, dbErr
}

func (cs *ClientService) Create(newClient models.ClientJson) error {

	var dbErr error = nil

	hashSecret := cs.token_svc.GenerateSecret(newClient.Name)
	var scopes = []models.Scope{}
	for i := 0; i < len(newClient.Scopes); i++ {
		scope := models.Scope{
			Name: newClient.Scopes[i],
		}
		scopes = append(scopes, scope)
	}
	newClientRow := models.Client{
		Name:        newClient.Name,
		Secret:      hashSecret,
		RedirectURI: newClient.RedirectURI,
		GrantType:   newClient.GrantType,
		Scopes:      scopes,
	}
	// Write to the pscale
	result := cs.db.Create(&newClientRow)
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}

func (cs *ClientService) Delete(id string) error {
	var client models.Client
	var dbErr error = nil
	result := cs.db.Delete(&client, models.Client{Name: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}

func (cs *ClientService) Authenticated(clientId string, clientSecret string) bool {
	var client models.Client
	result := cs.db.Model(&models.Client{}).First(&client, models.Client{Name: clientId})
	if result.Error != nil {
		return false
	}
	return client.Secret == clientSecret
}
