package main

import (
	"fmt"
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
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	signer.Claims = claims
	tokenString, err := signer.SignedString(SignKey)
	if err != nil {
		return Token{Token: ""}, err
	}
	return Token{Token: tokenString}, err
}

func (s *JWTAuthService) hasAuthorization(req Request) bool {
	httpReq, _ := req.(HttpRequest)
	authToken := httpReq.httpr.Header.Get("Authorization")
	return authToken != ""
}

func (s *JWTAuthService) isAuthorizedForUserId(userId int, req Request) bool {
	httpReq, _ := req.(HttpRequest)
	tokenString := httpReq.httpr.Header.Get("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate token is signed with the right method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return SignKey, nil
	})

	if err != nil {
		return false
	}

	// Verify that token is valid and that it belongs to the asked user.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenUserIDFloat, _ := claims["userId"].(float64)
		tokenUserID := int(tokenUserIDFloat)
		return tokenUserID == userId
	}

	return false
}
