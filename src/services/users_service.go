package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"gorm.io/gorm"
)

type IUserService interface {
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	AddAuthorizedClient(id int, cliend_id string) error
	AddAuthorizationCode(id int, challenge string) (string, error)
	VerifyAuthorizationCode(code, code_verifier, clientId string) (bool, error)
	Create(newUser User) error
	Delete(id int) error
}

type UserService struct {
	db *gorm.DB
}

type User struct {
	ID      int      `json:"id" sql:"id" gorm:"primaryKey"`
	Name    string   `json:"name" sql:"name"`
	Secret  string   `json:"secret" sql:"secret"`
	Clients []Client `json:"-" sql:"authorizedClients" gorm:"many2many:user_clients"`
	Codes   []Code   `json:"-" sql:"codes" gorm:"foreignKey:UserRefer"`
}

type Code struct {
	Secret    string
	Challenge string
	UserRefer uint
}

func (User) CreateTable() string {
	return "users"
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (cs *UserService) GetAll() ([]User, error) {
	var users []User
	var dbError error = nil
	result := cs.db.Model(&User{}).Find(&users)
	if result.Error != nil {
		dbError = result.Error
	}
	return users, dbError
}

func (cs *UserService) GetByID(id int) (User, error) {
	var user User
	var dbErr error = nil
	result := cs.db.Model(&User{}).First(&user, User{ID: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return user, dbErr
}

func (cs *UserService) Create(newUser User) error {

	var dbErr error = nil
	// Hash the password
	h := sha256.New()
	h.Write([]byte(newUser.Secret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))

	newUserRow := User{
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
	var client Client
	var dbErr error = nil
	result := cs.db.Model(&Client{}).Preload("Scopes").First(&client, Client{Name: cliend_id})
	if result.Error != nil {
		dbErr = result.Error
	}
	err := cs.db.Model(&User{}).Where("id = ?", true).Association("Clients").Append(&client)
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
	dbErr = cs.db.Model(&User{}).Where("id = ?", true).Association("Codes").Append(Code{Secret: auth_code, Challenge: challenge})
	if dbErr != nil {
		return code, dbErr
	}
	code = auth_code
	return code, nil
}

func (cs *UserService) VerifyAuthorizationCode(auth_code, code_verifier, clientId string) (bool, error) {
	var user User
	hasAuthorization := false
	var dbErr error = nil
	result := cs.db.Model(&User{}).Preload("Codes").Preload("Clients").Where("codes.secret = ?", auth_code).Where("clients.name = ?", clientId).First(&user)
	if result.Error != nil {
		dbErr = result.Error
	}
	if user.ID == 0 {
		return false, dbErr
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
	return hasAuthorization, dbErr
}

func (cs *UserService) Delete(id int) error {
	var user User
	var dbErr error = nil
	result := cs.db.Delete(&user, User{ID: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}
