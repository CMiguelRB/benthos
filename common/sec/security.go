package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
)

func Encrypt(input string) (output string, err error) {

	key, err := hex.DecodeString(os.Getenv("ENCRYPTION_KEY"))

	if err != nil {
		err = errors.New("Encryption key not found")
		return "", err
	}

	iv := make([]byte, 16)
	blockSize := aes.BlockSize
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}

	padding := (blockSize - len(input)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	intputBytes := append([]byte(input), padtext...)

	block, _ := aes.NewCipher(key)
	cipherText := make([]byte, len(intputBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, intputBytes)
	result := hex.EncodeToString(iv) + hex.EncodeToString(cipherText)
	return result, err
}

func Decrypt(input string) (output string, err error) {

	key, err := hex.DecodeString(os.Getenv("ENCRYPTION_KEY"))

	if err != nil {
		err = errors.New("Encryption key not found")
		return "", err
	}

	inputHex, err := hex.DecodeString(input)
	if err != nil {
		return "", err
	}

	inputBytes := []byte(inputHex)

	if len(inputBytes) < 16 {
		err = errors.New("Index out of bounds: decrypting input string")
		return "", err
	}

	cipherText := inputBytes[16:]

	iv := inputBytes[0:16]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	length := len(cipherText)
	if length == 0 {
		return "", nil
	}

	padlen := int(cipherText[length-1])
	if padlen > length || padlen == 0 {
		return "nil", nil
	}
	cipherText = cipherText[:length-padlen]

	return string(cipherText), err
}
