package utils

import "golang.org/x/crypto/bcrypt"

func Password(secret string, pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret+pwd), bcrypt.DefaultCost)
	return string(hash), err
}

func ComparePassword(secret string, inputPassword, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(secret+inputPassword))
	if err != nil {
		return err
	}
	return nil
}
