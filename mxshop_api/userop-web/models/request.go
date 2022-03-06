package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityID uint
	jwt.RegisteredClaims
}
