package aes

import (
	"crypto/aes"
	cipher2 "crypto/cipher"
	"fmt"
	"log"
)

func (cipher *aesCipher) Encrypt(plaintext string) string {
	plaintextBytes := []byte(plaintext)

	block, err := aes.NewCipher(cipher.key)
	if err != nil {
		log.Fatal(err.Error())
	}

	aesGCM, err := cipher2.NewGCM(block)
	if err != nil {
		log.Fatal(err.Error())
	}

	nonce := keyConvert(cipher.iv, aesGCM.NonceSize())

	ciphertext := aesGCM.Seal(nonce, nonce, plaintextBytes, nil)
	return fmt.Sprintf("%x", ciphertext)
}
