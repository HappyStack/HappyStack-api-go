package main

import "golang.org/x/crypto/bcrypt"

type BCryptEncryptionService struct{}

func NewBCryptEncryptionService() *BCryptEncryptionService {
	return &BCryptEncryptionService{}
}

func (es *BCryptEncryptionService) hashAndSalt(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	return string(hash), err
}

func (es *BCryptEncryptionService) comparePasswords(ciphered string, plain string) bool {
	byteHash := []byte(ciphered)
	bytePlain := []byte(plain)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		return false
	}
	return true
}
