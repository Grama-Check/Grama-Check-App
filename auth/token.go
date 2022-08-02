package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var rsaPrivateKey *rsa.PrivateKey

func init() {

	f, err := os.Open("private.pem")
	if err != nil {
		log.Fatal("Cannot read privatekey")
	}
	defer f.Close()
	buf := make([]byte, 0)
	io.ReadFull(f, buf)
	block, _ := pem.Decode([]byte(buf))
	rsaPrivateKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)

}

func generateToken() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute))}

	jwt := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtString, err := jwt.SignedString(rsaPrivateKey)

	if err != nil {
		return "", err
	}

	return jwtString, nil

}
