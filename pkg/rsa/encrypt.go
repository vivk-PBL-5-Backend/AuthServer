package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	log "github.com/sirupsen/logrus"
)

func (cipher *rsaCipher) Encrypt(plaintext string) string {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, cipher.publicKey, []byte(plaintext), nil)
	if err != nil {
		log.Error(err.Error())
	}
	return string(ciphertext)
}
