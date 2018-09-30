package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"time"
)

var VerifyKey, SignKey []byte

type JWTAuthService struct { }

func NewJWTAuthService() *JWTAuthService {
	return &JWTAuthService{ }
}

type Token struct {
	Token string `json:"token"`
}

const (
	privateKeyPath = "../keys/app.rsa"
	publicKeyPath  = "../keys/app.rsa.pub"
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

func (s *JWTAuthService) tokenFor(user UserCredentials) (Token, error) {
	signer := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["usename"] = user.Username
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	signer.Claims = claims
	tokenString, err := signer.SignedString(SignKey)
	if err != nil {
		return Token{Token: ""}	, err
	}
	return Token{Token: tokenString}, err
}
