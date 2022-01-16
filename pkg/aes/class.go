package aes

import (
	"crypto/aes"
)

type ICipher interface {
	Encrypt(plaintext string) string
	Decrypt(ciphertext string) string
}

type aesCipher struct {
	key       []byte
	iv        []byte
	blockSize int
}

func New(sourceKey []byte, sourceIV []byte) ICipher {
	blockSize := aes.BlockSize

	key := keyConvert(sourceKey, blockSize*2)
	iv := keyConvert(sourceIV, blockSize)

	ase := &aesCipher{
		key:       key,
		iv:        iv,
		blockSize: blockSize,
	}

	return ase
}
