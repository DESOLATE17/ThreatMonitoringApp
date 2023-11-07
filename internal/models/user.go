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

type UserLogin struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserSignUp struct {
	Login    string `json:"login" binding:"required,max=64"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
