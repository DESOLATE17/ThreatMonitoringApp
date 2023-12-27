package models

import (
	"time"
)

type User struct {
	UserId           int    `gorm:"primaryKey" json:"userId"`
	Login            string `json:"login" binding:"required,max=64"`
	IsAdmin          bool   `json:"isAdmin"`
	Name             string `json:"name,omitempty"`
	Password         string `json:"password,omitempty" binding:"required,min=8,max=64"`
	RegistrationDate time.Time
}

type UserLogin struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type UserSignUp struct {
	Login    string `json:"login" binding:"required,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
