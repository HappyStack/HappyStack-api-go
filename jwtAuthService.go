package main

import (
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var VerifyKey, SignKey []byte

type JWTAuthService struct{}

func NewJWTAuthService() *JWTAuthService {
	return &JWTAuthService{}
}

const (
	privateKeyPath = "keys/app.rsa"
	publicKeyPath  = "keys/app.rsa.pub"
)

func (s *JWTAuthService) init() {
	sKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
	SignKey = sKey
	vKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
	VerifyKey = vKey
}

func (s *JWTAuthService) tokenFor(user User) (Token, error) {
	signer := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["usename"] = user.Username
	claims["userId"] = user.Id
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	signer.Claims = claims
	tokenString, err := signer.SignedString(SignKey)
	if err != nil {
		return Token{Token: ""}, err
	}
	return Token{Token: tokenString}, err
}
