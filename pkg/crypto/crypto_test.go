package crypto

import (
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	password := "strong-master-password"
	salt := []byte("kylrix-ecosystem-default-salt-!!")
	data := "this-is-a-secret-message"

	key := DeriveKey(password, salt)

	encrypted, err := Encrypt(data, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != data {
		t.Errorf("Decrypted data mismatch: expected %v, got %v", data, decrypted)
	}
}

func TestComplexDataEncryption(t *testing.T) {
	password := "another-password"
	salt := []byte("random-salt-12345678901234567890")
	key := DeriveKey(password, salt)

	data := map[string]interface{}{
		"id":    123,
		"label": "test-data",
		"meta":  []string{"a", "b", "c"},
	}

	encrypted, err := Encrypt(data, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Verify unmarshaling worked
	decMap := decrypted.(map[string]interface{})
	if decMap["id"].(float64) != 123 { // JSON unmarshals numbers to float64 by default
		t.Errorf("Decrypted ID mismatch")
	}
}
