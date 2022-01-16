package aes

import (
	"crypto/aes"
)

type aesCipher struct {
	key       []byte
	iv        []byte
	blockSize int
}

func New(sourceKey []byte, sourceIV []byte) *aesCipher {
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
