package aes

import (
	"crypto/aes"
	cipher2 "crypto/cipher"
	"encoding/hex"
	"fmt"
)

func (cipher *aesCipher) Decrypt(ciphertext string) string {
	enc, _ := hex.DecodeString(ciphertext)
	encrypted := string(enc)

	block, err := aes.NewCipher(cipher.key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher2.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	plaintext, err := aesGCM.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
