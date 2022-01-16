package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	log "github.com/sirupsen/logrus"
)

func (cipher *rsaCipher) Decrypt(ciphertext string) string {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, cipher.privateKey, []byte(ciphertext), nil)
	if err != nil {
		log.Error(err.Error())
	}
	return string(plaintext)
}
