package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var AesKey []byte

func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return "", errors.New("加密失败" + err.Error())
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("加密失败" + err.Error())
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.New("加密失败" + err.Error())
	}

	cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)

	result := append(nonce, cipherText...)
	return base64.StdEncoding.EncodeToString(result), nil

}
func Decrypt(cipherText string) (string, error) {

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", errors.New("解密失败" + err.Error())
	}

	block, err := aes.NewCipher(AesKey)
	if err != nil {
		return "", errors.New("解密失败" + err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.New("解密失败" + err.Error())
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("解密失败 cipher太短了")
	}

	nonce := data[:nonceSize]
	cipherData := data[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", errors.New("解密失败" + err.Error())
	}
	return string(plainText), nil

}
