package models

import (
	"time"
)

type User struct {
	UserId           int    `gorm:"primaryKey"`
	Login            string `json:"login" binding:"required,max=64"`
	IsAdmin          bool
	Name             string `json:"name"`
	Password         string `json:"password" binding:"required,min=8,max=64"`
	RegistrationDate time.Time
}

func GetClientId() int {
	return 1
}

func GetAdminId() int {
	return 2
}
