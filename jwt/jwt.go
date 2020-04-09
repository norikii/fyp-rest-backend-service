package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/tatrasoft/fyp-rest-backend-service/model"
	"time"
)

var jwtSecretKey = []byte("jwt_secret_ket")

// creates jwt token when you signing in and signing out
func CreateJWT(email string) (response string, err error) {
	// for token expiration
	expirationTime := time.Now().Add(5 * time.Minute)
	// here and identity claim is being created
	claims := &model.Claim{
		Email:          email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// getting the auth token
	tokenToString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("unable to convert token: %v", err)
	}

	return tokenToString, nil
}

// verifies the identity claim by getting auth token and returning the email
func VerifyToken(tokenString string) (email string, err error) {
	claims := &model.Claim{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if token != nil {
		return claims.Email, nil
	}

	return "", err
}
