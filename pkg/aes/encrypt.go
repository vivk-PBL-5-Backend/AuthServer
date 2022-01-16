package aes

import (
	"crypto/aes"
	cipher2 "crypto/cipher"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

func (cipher *aesCipher) Encrypt(plaintext string) string {
	bytePlaintext := cipher.paddingPKCS5([]byte(plaintext), cipher.blockSize)
	ciphertext := make([]byte, len(bytePlaintext))

	block, err := aes.NewCipher(cipher.key)
	if err != nil {
		log.Fatal(err.Error())
	}

	modeEncrypter := cipher2.NewCBCEncrypter(block, cipher.iv)
	modeEncrypter.CryptBlocks(ciphertext, bytePlaintext)
	return hex.EncodeToString(ciphertext)
}
