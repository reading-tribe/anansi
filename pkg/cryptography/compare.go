package cryptography

import (
	"golang.org/x/crypto/bcrypt"
)

func ComparePasswords(hashedPwd string, plainPwd []byte) error {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return err
	}

	return nil
}
