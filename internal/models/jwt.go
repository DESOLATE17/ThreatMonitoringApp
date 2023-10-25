package models

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	jwt.StandardClaims
	UserId  int
	IsAdmin bool
}

type Role int

const (
	Client Role = iota // 0
	Admin              // 1
)
