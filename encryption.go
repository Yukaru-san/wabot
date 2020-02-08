package wabot

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// EncryptData uses AES to encrypt using the enckey
func EncryptData(data []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(encrypKey)
	if err != nil {
		panic(err)
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(data))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext
}

// DecryptData decrypts a given data with the encKey
func DecryptData(data []byte) []byte {

	// Create the AES cipher
	block, err := aes.NewCipher(encrypKey)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(data) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := data[:aes.BlockSize]

	// Remove the IV from the ciphertext
	data = data[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(data, data)

	return data
}
