package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type JWTToken struct {
	UserID interface{}
	Name string
	Email string
	*jwt.StandardClaims
}

func CreateJWTToken(userID interface{}, name string, email string) (string, error) {
	tokenStruct := JWTToken{
		UserID:         userID,
		Name:           name,
		Email:          email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}

	// generate JWT token
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenStruct)

	// creates complete signed token
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", fmt.Errorf("unable to sign the token: %v", err)
	}

	return tokenString, err
}

func IsValidJWTToken(tokenString string, ) (bool, error) {
	tokenStruct := &JWTToken{}

	_, err := jwt.ParseWithClaims(tokenString, tokenStruct, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return false, fmt.Errorf("not valid JWT token: %v", err)
	}

	return true, nil
}
