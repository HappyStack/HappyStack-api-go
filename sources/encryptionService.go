package main

type EncryptionService interface {
	comparePasswords(ciphered string, plain string) bool
}