package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"nak-auth/models"

	"gorm.io/gorm"
)

type IUserService interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (models.User, error)
	AddAuthorizedClient(id int, cliend_id string) error
	AddAuthorizationCode(id int, challenge string) (string, error)
	VerifyAuthorizationCode(code, code_verifier, clientId string) (models.User, error)
	Create(newUser models.User) error
	Delete(id int) error
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (cs *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	var dbError error = nil
	result := cs.db.Model(&models.User{}).Find(&users)
	if result.Error != nil {
		dbError = result.Error
	}
	return users, dbError
}

func (cs *UserService) GetByID(id int) (models.User, error) {
	var user models.User
	var dbErr error = nil
	result := cs.db.Model(&models.User{}).First(&user, models.User{ID: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return user, dbErr
}

func (cs *UserService) Create(newUser models.User) error {

	var dbErr error = nil
	// Hash the password
	h := sha256.New()
	h.Write([]byte(newUser.Secret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))

	newUserRow := models.User{
		Name:   newUser.Name,
		Secret: hashSecret,
	}
	// Write to the pscale
	result := cs.db.Create(&newUserRow)
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}

func (cs *UserService) AddAuthorizedClient(id int, cliend_id string) error {
	var client models.Client
	var dbErr error = nil
	result := cs.db.Model(&models.Client{}).Preload("Scopes").First(&client, models.Client{Name: cliend_id})
	if result.Error != nil {
		dbErr = result.Error
	}
	err := cs.db.Model(&models.User{}).Where("id = ?", true).Association("Clients").Append(&client)
	if err != nil {
		dbErr = err
	}
	return dbErr
}

func (cs *UserService) AddAuthorizationCode(id int, challenge string) (string, error) {
	var dbErr error = nil
	var code string = ""
	buf := make([]byte, 64)
	_, dbErr = rand.Read(buf)
	if dbErr != nil {
		return code, dbErr
	}
	auth_code := base64.StdEncoding.EncodeToString(buf)
	dbErr = cs.db.Model(&models.User{}).Where("id = ?", true).Association("Codes").Append(models.Code{Secret: auth_code, Challenge: challenge})
	if dbErr != nil {
		return code, dbErr
	}
	code = auth_code
	return code, nil
}

func (cs *UserService) VerifyAuthorizationCode(auth_code, code_verifier, clientId string) (models.User, error) {
	var user models.User
	hasAuthorization := false
	var dbErr error = nil
	result := cs.db.Model(&models.User{}).Preload("Codes").Preload("Clients").Where("codes.secret = ?", auth_code).Where("clients.name = ?", clientId).First(&user)
	if result.Error != nil {
		dbErr = result.Error
		return user, dbErr
	}
	for _, code := range user.Codes {
		if code.Secret == auth_code {
			h := sha256.New()
			h.Write([]byte(code_verifier))
			code_challenge := base64.URLEncoding.EncodeToString(h.Sum(nil))
			if code.Challenge == code_challenge {
				hasAuthorization = true
				break
			}
		}
	}
	if hasAuthorization == true {
		return user, nil
	}
	return models.User{}, errors.New("invalid authorization code")
}

func (cs *UserService) Delete(id int) error {
	var user models.User
	var dbErr error = nil
	result := cs.db.Delete(&user, models.User{ID: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}
