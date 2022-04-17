package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type Encoder struct {
	SecretKey string
}

func (e *Encoder) Encrypt(value string) (string, error) {
	// получаем cipher.Block
	aesblock, err := aes.NewCipher([]byte(e.SecretKey))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce, err := GenerateRandom(aesgcm.NonceSize())
	if err != nil {
		return "", err
	}

	// зашифровываем
	dst := aesgcm.Seal(nonce, nonce, []byte(value), nil)

	return hex.EncodeToString(dst), nil
}

func (e *Encoder) Decrypt(value []byte) (string, error) {
	var decodedValue []byte
	decodedValue, _ = hex.DecodeString(string(value))

	// получаем cipher.Block
	aesblock, err := aes.NewCipher([]byte(e.SecretKey))
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	// создаём вектор инициализации
	nonce, cipherText := decodedValue[:aesgcm.NonceSize()], decodedValue[aesgcm.NonceSize():]

	// расшифровываем
	userID, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(userID), nil
}
