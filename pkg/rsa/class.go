package rsa

import "crypto/rsa"

type rsaCipher struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func New(privateBytes []byte) *rsaCipher {
	publicKey, privateKey := generateKeyPair(privateBytes)
	cipher := &rsaCipher{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
	return cipher
}
