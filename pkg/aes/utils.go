package aes

import "bytes"

func (cipher *aesCipher) paddingPKCS5(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, paddingText...)
}

func keyConvert(iv []byte, blockSize int) []byte {
	newIV := make([]byte, blockSize)
	lenIV := len(iv)

	for i := 0; i < blockSize; i++ {
		newIV[i] = iv[i%lenIV]
	}

	return newIV
}
