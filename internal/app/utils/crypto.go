package utils

import (
	"crypto/aes"
	"crypto/cipher"
)

func Encrypt(msg string, key string) (string, error) {
	aesBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGcm.NonceSize())

	encryptMsg := aesGcm.Seal(nil, nonce, []byte(msg), nil)

	return string(encryptMsg), nil
}

func Decrypt(msg string, key string) (string, error) {
	aesBlock, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGcm.NonceSize())

	decryptMsg, err := aesGcm.Open(nil, nonce, []byte(msg), nil)
	if err != nil {
		return "", err
	}

	return string(decryptMsg), nil
}
