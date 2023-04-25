package services

import (
	"crypto/sha256"
	"encoding/base64"

	"gorm.io/gorm"
)

type IUserService interface {
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	Create(newUser User) error
	Delete(id int) error
}

type UserService struct {
	db *gorm.DB
}

type User struct {
	ID                int    `json:"id" sql:"id"`
	Name              string `json:"name" sql:"name"`
	Secret            string `json:"secret" sql:"secret"`
	AuthorizedClients string `json:"-" sql:"authorized_clients" gorm:"many2many:user_clients"`
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (cs *UserService) GetAll() ([]User, error) {
	var users []User
	var dbError error = nil
	result := cs.db.Model(&User{}).Preload("AuthorizedClients").Find(&users)
	if result.Error != nil {
		dbError = result.Error
	}
	return users, dbError
}

func (cs *UserService) GetByID(id int) (User, error) {
	var user User
	var dbErr error = nil
	result := cs.db.Model(&User{}).Preload("AuthorizedClients").First(&user, User{ID: id})
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

func (cs *UserService) Delete(id int) error {
	var user User
	var dbErr error = nil
	result := cs.db.Delete(&user, User{ID: id})
	if result.Error != nil {
		dbErr = result.Error
	}
	return dbErr
}
