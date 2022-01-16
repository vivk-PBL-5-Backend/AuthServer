package aes

import (
	"crypto/aes"
	cipher2 "crypto/cipher"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

func (cipher *aesCipher) Decrypt(ciphertext string) string {
	ciphertextDecode, err := hex.DecodeString(ciphertext)
	if err != nil {
		log.Fatal(err.Error())
	}

	block, err := aes.NewCipher(cipher.key)
	if err != nil {
		log.Fatal(err.Error())
	}

	plaintextBytes := make([]byte, cipher.blockSize)
	modeDecrypter := cipher2.NewCBCDecrypter(block, cipher.iv)
	modeDecrypter.CryptBlocks(plaintextBytes, ciphertextDecode)
	return string(plaintextBytes)
}
