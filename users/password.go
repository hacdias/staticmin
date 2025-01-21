package users

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// randomPasswordBytesCount is chosen to fit in a base64 string without padding
const randomPasswordBytesCount = 9

// HashPwd hashes a password.
func HashPwd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPwd checks if a password is correct.
func CheckPwd(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandomPwd() (string, error) {
	randomPasswordBytes := make([]byte, randomPasswordBytesCount)
	var _, err = rand.Read(randomPasswordBytes)
	if err != nil {
		return "", err
	}

	// This is done purely to make the password human-readable
	var randomPasswordString = base64.URLEncoding.EncodeToString(randomPasswordBytes)
	return randomPasswordString, nil
}
