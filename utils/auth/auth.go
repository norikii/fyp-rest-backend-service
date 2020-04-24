package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const (

)

type JWTToken struct {
	UserID interface{}
	Name string
	Email string
	IsAdmin bool
	*jwt.StandardClaims
}

func CreateJWTToken(userID interface{}, name string, email string, isAdmin bool) (string, error) {
	//tokenTTL, err := strconv.Atoi(os.Getenv("JWT_TOKEN_TTL"))
	//if err != nil {
	//	return "", fmt.Errorf("unable to get token ttl from .env: %v", err)
	//}
	tokenStruct := JWTToken{
		UserID:         userID,
		Name:           name,
		Email:          email,
		IsAdmin:		false,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
		},
	}

	if isAdmin {
		tokenStruct.IsAdmin = true
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

func IsValidJWTToken(tokenString string) (isValidToken bool, isAdmin bool ,err error) {
	tokenStruct := &JWTToken{}

	token, err := jwt.ParseWithClaims(tokenString, tokenStruct, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return false, false, fmt.Errorf("not valid JWT token: %v", err)
	}

	claims := token.Claims.(*JWTToken)
	if claims.IsAdmin == true {
		return true, true, nil
	}

	return true, false, nil
}
