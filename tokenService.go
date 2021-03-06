package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey = "asupersecretekeythatshouldntbehardcodedhere"
)

type ITokenService interface {
	CreateToken(user User) (string, error)
	ValidateToken(token string) (interface{}, error)
}

type TokenService struct{}

var TS ITokenService

// Creates a new Token for a User
func (t *TokenService) CreateToken(user User) (string, error) {
	if (User{}) == user {
		return "", fmt.Errorf("Cannot create token for empty user")
	}
	safeUser := UserSafe{}
	safeUser.ID = user.ID
	safeUser.Email = user.Email
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": safeUser,
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// TODO: Remove hardcoded signed string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Error signing token: ", err)
		return "", err
	}
	return tokenString, nil
}

func (t *TokenService) ValidateToken(tokenString string) (interface{}, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims["sub"], nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, fmt.Errorf("Malformed Token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return false, fmt.Errorf("Token is Expired")
		} else {
			return false, fmt.Errorf("Couldn't handle the token")
		}
	} else {
		return false, fmt.Errorf("Error handling token: %v", err)
	}
}
