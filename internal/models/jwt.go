package models

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	jwt.StandardClaims
	UserId  int  `json:"userId"`
	IsAdmin bool `json:"isAdmin"`
}

type Role int

const (
	Client Role = iota // 0
	Admin              // 1
)
