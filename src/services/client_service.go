package services

import (
	"nak-auth/db"
	"nak-auth/models"
	"os"
)

type IClientService interface {
	GetAll() ([]models.Client, error)
	GetByID(id string) (models.Client, error)
	Create(newClient models.ClientJson) error
	Delete(id string) error
	Authenticated(clientId string, clientSecret string) bool
}

type ClientService struct {
	clientId  string
	token_svc ITokenService
	fact      db.ILibSqlClientFactory
}

// Create a new client service. This service is responsible for managing clients.
func NewClientService(fact db.ILibSqlClientFactory, tkn_svc ITokenService) *ClientService {
	name := os.Getenv("NAK_AUTH_CLIENT_ID")
	secret := os.Getenv("NAK_AUTH_CLIENT_SECRET")
	db, _ := fact.CreateClient()
	db.Create(&models.Client{Name: name, Secret: secret, GrantType: "client_credentials", RedirectURI: "http://localhost:8080", Scopes: []models.Scope{{Name: "internal"}}})
	return &ClientService{fact: fact, token_svc: tkn_svc, clientId: name}
}

func (cs *ClientService) GetAll() ([]models.Client, error) {
	var clients []models.Client
	var dbError error = nil
	dbConn, err := cs.fact.CreateClient()
	if err != nil {
		return clients, err
	}
	result := dbConn.Model(&models.Client{}).Preload("Scopes").Where("name != ?", cs.clientId).Find(&clients)
	if result.Error != nil {
		dbError = result.Error
	}
	return clients, dbError
}

func (cs *ClientService) GetByID(id string) (models.Client, error) {
	var client models.Client
	var dbErr error = nil
	dbConn, err := cs.fact.CreateClient()
	if err != nil {
		return client, err
	}
	result := dbConn.Model(&models.Client{}).Preload("Scopes").First(&client, models.Client{Name: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return client, dbErr
}

func (cs *ClientService) Create(newClient models.ClientJson) error {

	var dbErr error = nil
	dbConn, err := cs.fact.CreateClient()
	if err != nil {
		return err
	}

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
	result := dbConn.Create(&newClientRow)
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}

func (cs *ClientService) Delete(id string) error {
	var client models.Client
	var dbErr error = nil
	dbConn, err := cs.fact.CreateClient()
	if err != nil {
		return err
	}
	result := dbConn.Delete(&client, models.Client{Name: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}

func (cs *ClientService) Authenticated(clientId string, clientSecret string) bool {
	var client models.Client
	dbConn, err := cs.fact.CreateClient()
	if err != nil {
		return false
	}
	result := dbConn.Model(&models.Client{}).First(&client, models.Client{Name: clientId})
	if result.Error != nil {
		return false
	}
	return client.Secret == clientSecret
}
