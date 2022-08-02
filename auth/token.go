package auth

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken() (string, error) {
	buffer, err := ioutil.ReadFile("private.pem")
	if err != nil {
		log.Fatal("Cannot read privatekey")
	}

	block, _ := pem.Decode(buffer)

	if err != nil {
		log.Fatal("Cannot parse private.pem file")
	}
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		log.Fatal("Cannot parse private key")
	}

	claims := jwt.RegisteredClaims{}

	jwt := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtString, err := jwt.SignedString(rsaPrivateKey)

	if err != nil {
		return "", err
	}

	return jwtString, nil

}
