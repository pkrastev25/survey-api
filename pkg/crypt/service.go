package crypt

import "golang.org/x/crypto/bcrypt"

type CryptService struct {
}

func NewCryptService() CryptService {
	return CryptService{}
}

func (service CryptService) GeneratePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
