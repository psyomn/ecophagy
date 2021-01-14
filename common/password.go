package common

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, string) {
	saltBytes := make([]byte, 8)
	_, err := rand.Read(saltBytes)
	if err != nil {
		panic(err)
	}
	saltStr := string(saltBytes)
	saltedPassword := []byte(saltStr + password)

	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, 14)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword), saltStr
}

func EncryptPasswordWithSalt(password string, salt string) string {
	saltedPassword := []byte(salt + password)
	hashedPassword, err := bcrypt.GenerateFromPassword(saltedPassword, 14)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
