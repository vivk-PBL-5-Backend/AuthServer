package cipheradapter

type ICipher interface {
	Encrypt(plaintext string) string
	Decrypt(ciphertext string) string
}
