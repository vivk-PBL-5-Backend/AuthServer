package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
)

func generateKeyPair(privateBytes []byte) (*rsa.PublicKey, *rsa.PrivateKey) {
	block, _ := pem.Decode(privateBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	return &privateKey.PublicKey, privateKey
}
