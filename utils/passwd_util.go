package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashAndSaltPwd takes a password and uses hashing algorithm to produce hashed password
func HashAndSaltPwd(pwd string) (string, error) {
	bytePwd := []byte(pwd)

	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("unable to generate hash from the password: %v", err)
	}

	return string(hash), nil
}

// IsValidPassword
func IsValidPassword(hashedPwd string, enteredPwd string) (bool, error) {
	byteEnteredPwd := []byte(enteredPwd)
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, byteEnteredPwd)
	if err != nil {
		return false, fmt.Errorf("password and hash does not match: %v", err)
	}

	return true, nil
}


