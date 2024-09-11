package models

import (
	"github.com/golang-jwt/jwt"
)

type Token struct {
	Value string
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwtKey = []byte("ordnung")

func (token *Token) GenerateToken(email string) *Token {
	claims := &Claims{
		Email: email,
	}
	createdToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenValue, err := createdToken.SignedString(jwtKey)
	if err != nil {
		return token
	}
	token.Value = tokenValue
	return token
}

func (token *Token) ValidToken(email string) bool {
	newToken := new(Token)
	newToken.GenerateToken(email)
	return token.Value == newToken.Value
}