package services

import (
	"crypto/sha256"
	"encoding/base64"

	"gorm.io/gorm"
)

type LoginService struct {
	db *gorm.DB
}

type ILoginService interface {
	Login(username string, secret string) (bool, int)
}

func NewLoginService(db *gorm.DB) *LoginService {
	return &LoginService{db: db}
}

type Login struct {
	Name   string `json:"username"`
	Secret string `json:"secret"`
}

func (ls *LoginService) Login(username string, secret string) (bool, int) {
	var user User
	success := false
	userId := -1
	h := sha256.New()
	h.Write([]byte(secret))
	hashSecret := base64.URLEncoding.EncodeToString(h.Sum(nil))
	result := ls.db.Model(&User{}).First(&user, User{Name: username, Secret: hashSecret})
	if result.Error != nil {
		success = false
	} else {
		success = true
		userId = user.ID
	}
	return success, userId
}
