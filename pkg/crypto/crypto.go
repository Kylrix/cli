package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	PBKDF2Iterations = 600000
	SaltSize         = 32
	IVSize           = 16 // Matching EcosystemSecurity.IV_SIZE
	KeySize          = 32 // 256 bits
)

// DeriveKey derives a key from a password and salt matching the EcosystemSecurity implementation
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeySize, sha256.New)
}

// Encrypt encrypts data using AES-256-GCM, matching the Kylrix Ecosystem Security Protocol
func Encrypt(data interface{}, key []byte) (string, error) {
	// 1. JSON Stringify the data (matching TS implementation)
	plaintext, err := json.Marshal(data)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal data to JSON")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create cipher")
	}

	// Web Crypto's AES-GCM with 16-byte IV requires NewGCMWithNonceSize in Go
	gcm, err := cipher.NewGCMWithNonceSize(block, IVSize)
	if err != nil {
		return "", errors.Wrap(err, "failed to create GCM")
	}

	nonce := make([]byte, IVSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "failed to generate nonce")
	}

	// 2. Encrypt
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// 3. Combine IV + Ciphertext
	combined := append(nonce, ciphertext...)

	// 4. Base64 encode
	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt decrypts data using AES-256-GCM, matching the Kylrix Ecosystem Security Protocol
func Decrypt(encryptedBase64 string, key []byte) (interface{}, error) {
	// 1. Base64 decode
	combined, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode base64")
	}

	if len(combined) < IVSize {
		return nil, errors.New("encrypted data too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cipher")
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, IVSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create GCM")
	}

	// 2. Split IV and Ciphertext
	nonce := combined[:IVSize]
	ciphertext := combined[IVSize:]

	// 3. Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt (likely wrong key or corrupted data)")
	}

	// 4. JSON Unmarshal (matching TS implementation)
	var result interface{}
	if err := json.Unmarshal(plaintext, &result); err != nil {
		// If it's not valid JSON, it might be a raw string that was stringified
		// but since we stringified it, it should be valid JSON.
		return nil, errors.Wrap(err, "failed to unmarshal decrypted JSON")
	}

	return result, nil
}
