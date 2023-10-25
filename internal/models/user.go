package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserId           int `gorm:"primaryKey"`
	Login            string
	IsAdmin          bool
	Name             string
	PasswordHash     string
	RegistrationDate time.Time
}

func GetClientId() int {
	return 1
}

func GetAdminId() int {
	return 2
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTClaims struct {
	jwt.StandardClaims           // все что точно необходимо по RFC
	UserUUID           uuid.UUID `json:"user_uuid"`            // наши данные - uuid этого пользователя в базе данных
	Scopes             []string  `json:"scopes" json:"scopes"` // список доступов в нашей системе
}
