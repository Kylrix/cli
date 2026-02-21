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
	PBKDF2Iterations      = 600000
	PinPBKDF2Iterations   = 100000
	SessionPBKDF2Iterations = 5000
	SaltSize              = 32
	PinSaltSize           = 16
	SessionSaltSize       = 16
	IVSize                = 16 // Matching EcosystemSecurity.IV_SIZE
	KeySize               = 32 // 256 bits
)

// DeriveKey derives a key from a password and salt matching the EcosystemSecurity implementation
func DeriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, PBKDF2Iterations, KeySize, sha256.New)
}

// DerivePinKey derives a key from a PIN and salt for verifier storage
func DerivePinKey(pin string, salt []byte) []byte {
	return pbkdf2.Key([]byte(pin), salt, PinPBKDF2Iterations, KeySize, sha256.New)
}

// DeriveEphemeralKey derives a fast key for session piggybacking
func DeriveEphemeralKey(pin string, salt []byte) []byte {
	return pbkdf2.Key([]byte(pin), salt, SessionPBKDF2Iterations, KeySize, sha256.New)
}

// WrapKey wraps the MEK with an ephemeral key using AES-GCM
func WrapKey(rawKey []byte, ephemeralKey []byte) (string, error) {
	block, err := aes.NewCipher(ephemeralKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, IVSize)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, IVSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, rawKey, nil)
	combined := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(combined), nil
}

// UnwrapKey unwraps the MEK using the ephemeral key
func UnwrapKey(wrappedKeyBase64 string, ephemeralKey []byte) ([]byte, error) {
	combined, err := base64.StdEncoding.DecodeString(wrappedKeyBase64)
	if err != nil {
		return nil, err
	}

	if len(combined) < IVSize {
		return nil, errors.New("wrapped key too short")
	}

	block, err := aes.NewCipher(ephemeralKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(block, IVSize)
	if err != nil {
		return nil, err
	}

	nonce := combined[:IVSize]
	ciphertext := combined[IVSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

// ZeroBytes explicitly overwrites a byte slice with zeroes for security
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
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
