package utils_test

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/hrygo/echomind/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	keyHex, err := utils.GenerateRandomKey()
	assert.NoError(t, err)

	key, err := hex.DecodeString(keyHex)
	assert.NoError(t, err)

	plaintext := "This is a secret message."

	// Test successful encryption and decryption
	ciphertext, err := utils.Encrypt(plaintext, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext)

	decryptedText, err := utils.Decrypt(ciphertext, key)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decryptedText)

	// Test with different plaintext
	plaintext2 := "Another secret."
	ciphertext2, err := utils.Encrypt(plaintext2, key)
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext2)
	assert.NotEqual(t, ciphertext, ciphertext2) // Ciphertext should be different due to random nonce

	decryptedText2, err := utils.Decrypt(ciphertext2, key)
	assert.NoError(t, err)
	assert.Equal(t, plaintext2, decryptedText2)
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	keyHex, err := utils.GenerateRandomKey()
	assert.NoError(t, err)
	key, err := hex.DecodeString(keyHex)
	assert.NoError(t, err)

	// Test with invalid base64 string
	_, err = utils.Decrypt("invalid-base64", key)
	assert.Error(t, err)

	// Test with too short ciphertext (less than nonce size)
	_, err = utils.Decrypt(base64.StdEncoding.EncodeToString([]byte("short")), key)
	assert.Error(t, err)
}

func TestDecrypt_WrongKey(t *testing.T) {
	keyHex1, err := utils.GenerateRandomKey()
	assert.NoError(t, err)
	key1, err := hex.DecodeString(keyHex1)
	assert.NoError(t, err)

	keyHex2, err := utils.GenerateRandomKey()
	assert.NoError(t, err)
	key2, err := hex.DecodeString(keyHex2)
	assert.NoError(t, err)

	plaintext := "This is a secret message."
	ciphertext, err := utils.Encrypt(plaintext, key1)
	assert.NoError(t, err)

	// Attempt to decrypt with a different key
	_, err = utils.Decrypt(ciphertext, key2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cipher: message authentication failed")
}

func TestGenerateRandomKey(t *testing.T) {
	keyHex, err := utils.GenerateRandomKey()
	assert.NoError(t, err)
	assert.Len(t, keyHex, 64) // 32 bytes * 2 hex chars per byte

	key, err := hex.DecodeString(keyHex)
	assert.NoError(t, err)
	assert.Len(t, key, 32) // 32 bytes
}
