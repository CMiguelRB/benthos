package sec

import (
	"benthos/config"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func Encrypt(input string) (output string, err error) {

	key, err := hex.DecodeString(config.Settings.Security.EncryptionKey)

	if err != nil {
		return "", err
	}

	iv := make([]byte, 16)
	blockSize := aes.BlockSize
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}

	padding := (blockSize - len(input)%blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	inputBytes := append([]byte(input), padText...)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	
	cipherText := make([]byte, len(inputBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, inputBytes)
	result := hex.EncodeToString(iv) + hex.EncodeToString(cipherText)
	return result, err
}

func Decrypt(input string) (output string, err error) {

	key, err := hex.DecodeString(config.Settings.Security.EncryptionKey)

	if err != nil {
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

	padLen := int(cipherText[length-1])
	if padLen > length || padLen == 0 {
		return "nil", nil
	}
	cipherText = cipherText[:length-padLen]

	return string(cipherText), err
}
