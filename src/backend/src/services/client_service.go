package services

import (
	"crypto/sha256"
	"encoding/base64"

	"gorm.io/gorm"
)

type IClientService interface {
	GetAll() ([]Client, error)
	GetByID(id string) (Client, error)
	Create(newClient ClientJson) error
	Delete(id string) error
}

type ClientService struct {
	db *gorm.DB
}

type Client struct {
	Name      string `sql:"name" gorm:"primaryKey"`
	Secret    string `sql:"secret"`
	GrantType string `sql:"grant_type"`

	RedirectURI string  `sql:"redirect_uri"`
	Scopes      []Scope `sql:"scopes" gorm:"many2many:client_scopes"`
}

type Scope struct {
	Name string `json:"name" sql:"name" gorm:"primaryKey"`
}

type ClientJson struct {
	Name      string `json:"name"`
	Secret    string `json:"secret"`
	GrantType string `json:"grant_type"`

	RedirectURI string   `json:"redirect_uri"`
	Scopes      []string `json:"scope"`
}

type from[To any] interface {
	From() To
}

func ListOfClientsToListOfClientJson(clients []Client) []ClientJson {
	var client_json = []ClientJson{}
	for i := 0; i < len(clients); i++ {
		cJson := clients[i].From()
		client_json = append(client_json, cJson)
	}
	return client_json
}

func (scope Scope) From() string {
	return scope.Name
}

func (client Client) From() ClientJson {
	scopes := []string{}
	for i := 0; i < len(client.Scopes); i++ {
		scopes = append(scopes, client.Scopes[i].Name)
	}
	client_json := ClientJson{
		Name:        client.Name,
		GrantType:   client.GrantType,
		Secret:      client.Secret,
		RedirectURI: client.RedirectURI,
		Scopes:      scopes,
	}
	return client_json
}

func (cJson ClientJson) From() Client {
	scopes := []Scope{}
	for i := 0; i < len(cJson.Scopes); i++ {
		scope := Scope{
			Name: cJson.Scopes[i],
		}
		scopes = append(scopes, scope)
	}
	client := Client{
		Name:        cJson.Name,
		GrantType:   cJson.GrantType,
		Secret:      cJson.Secret,
		RedirectURI: cJson.RedirectURI,
		Scopes:      scopes,
	}
	return client
}

func (Scope) CreateTable() string {
	return "scopes"
}

func (Client) CreateTable() string {
	return "clients"
}

func NewClientService(db *gorm.DB) *ClientService {
	return &ClientService{db: db}
}

func (cs *ClientService) GetAll() ([]Client, error) {
	var clients []Client
	var dbError error = nil
	result := cs.db.Model(&Client{}).Preload("Scopes").Find(&clients)
	if result.Error != nil {
		dbError = result.Error
	}
	return clients, dbError
}

func (cs *ClientService) GetByID(id string) (Client, error) {
	var client Client
	var dbErr error = nil
	result := cs.db.Model(&Client{}).Preload("Scopes").First(&client, Client{Name: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return client, dbErr
}

func (cs *ClientService) Create(newClient ClientJson) error {

	var dbErr error = nil
	// Hash the password
	h := sha256.New()
	h.Write([]byte(newClient.Secret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))

	var scopes = []Scope{}
	for i := 0; i < len(newClient.Scopes); i++ {
		scope := Scope{
			Name: newClient.Scopes[i],
		}
		scopes = append(scopes, scope)
	}
	newClientRow := Client{
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
	var client Client
	var dbErr error = nil
	result := cs.db.Delete(&client, Client{Name: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}
