package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test vakid encode/decode data.
func TestEncryptDecryptValidData(t *testing.T) {
	key := "password"
	value := "!#@!MDF@#FMxcsa3we*"

	t.Run("encrypt/Decrypt a valid data", func(t *testing.T) {
		encodedVal, err := EncryptAES([]byte(value), key)
		assert.NoError(t, err)

		decodedVal, err := DecryptAES([]byte(encodedVal), key)
		assert.NoError(t, err)
		assert.EqualValues(t, value, string(decodedVal))
	})

	t.Run("encrypt/Decrypt an invalid data", func(t *testing.T) {
		encodedVal, err := EncryptAES([]byte(value), key)
		assert.NoError(t, err)

		decodedVal, err := DecryptAES([]byte(encodedVal), key)

		assert.NoError(t, err)
		assert.NotEqualValues(t, value+"md5", decodedVal)
	})

	t.Run("encrypt/Decrypt an invalid data wrong key", func(t *testing.T) {
		encodedVal, err := EncryptAES([]byte(value), key)
		assert.NoError(t, err)

		_, err = DecryptAES([]byte(encodedVal), key+"md5")
		assert.Error(t, err)
	})
}
