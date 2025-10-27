package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/zalando/go-keyring"
)

const (
	// KeyringService is the service name used for keyring storage
	KeyringService = "globus-connect-server"

	// KeyringUser is the username used for keyring storage
	KeyringUser = "encryption-key"

	// KeyVersion is the current encryption key version for rotation support
	KeyVersion = "v1"

	// EncryptionKeySize is the size of the AES-256 key in bytes
	EncryptionKeySize = 32

	// NonceSize is the size of the GCM nonce in bytes (12 bytes = 96 bits)
	NonceSize = 12
)

// EncryptedData represents encrypted data with versioning information.
type EncryptedData struct {
	// Version is the key version used for encryption (for key rotation)
	Version string

	// Nonce is the randomly generated nonce for GCM
	Nonce []byte

	// Ciphertext is the encrypted data + authentication tag
	Ciphertext []byte
}

// GetOrCreateEncryptionKey retrieves the encryption key from the system keyring,
// or creates a new one if it doesn't exist.
//
// The key is stored in the system keyring:
//   - macOS: Keychain
//   - Linux: Secret Service API (gnome-keyring, kwallet)
//   - Windows: Credential Manager
//
// If the keyring is not available, returns an error with instructions.
func GetOrCreateEncryptionKey() ([]byte, error) {
	// Try to get existing key
	keyString, err := keyring.Get(KeyringService, KeyringUser)
	if err == nil {
		// Decode existing key from base64
		key, err := base64.StdEncoding.DecodeString(keyString)
		if err != nil {
			return nil, fmt.Errorf("decode encryption key: %w", err)
		}

		if len(key) != EncryptionKeySize {
			return nil, fmt.Errorf("invalid encryption key size: %d (expected %d)", len(key), EncryptionKeySize)
		}

		return key, nil
	}

	// If key doesn't exist, create a new one
	if err == keyring.ErrNotFound {
		key, err := generateEncryptionKey()
		if err != nil {
			return nil, fmt.Errorf("generate encryption key: %w", err)
		}

		// Store in keyring (base64 encoded for safe storage)
		keyString := base64.StdEncoding.EncodeToString(key)
		if err := keyring.Set(KeyringService, KeyringUser, keyString); err != nil {
			return nil, fmt.Errorf("store encryption key in keyring: %w\n\n"+
				"Keyring storage is required for secure token encryption.\n"+
				"Please ensure your system keyring is available:\n"+
				"  - macOS: Keychain (built-in)\n"+
				"  - Linux: gnome-keyring or kwallet\n"+
				"  - Windows: Credential Manager (built-in)", err)
		}

		return key, nil
	}

	// Keyring not available or other error
	return nil, fmt.Errorf("access system keyring: %w\n\n"+
		"Keyring storage is required for secure token encryption.\n"+
		"Please ensure your system keyring is available:\n"+
		"  - macOS: Keychain (built-in)\n"+
		"  - Linux: Install gnome-keyring or kwallet\n"+
		"  - Windows: Credential Manager (built-in)", err)
}

// generateEncryptionKey generates a new 256-bit (32-byte) encryption key
// using cryptographically secure random number generation.
func generateEncryptionKey() ([]byte, error) {
	key := make([]byte, EncryptionKeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("generate random key: %w", err)
	}
	return key, nil
}

// Encrypt encrypts plaintext data using AES-256-GCM.
//
// AES-256-GCM provides:
//   - Confidentiality: Data is encrypted with AES-256
//   - Authenticity: GCM authentication tag prevents tampering
//   - Performance: Hardware-accelerated on modern CPUs
//
// Returns encrypted data with version and nonce for decryption.
func Encrypt(plaintext []byte) (*EncryptedData, error) {
	// Get or create encryption key
	key, err := GetOrCreateEncryptionKey()
	if err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	// Generate random nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	// Encrypt and authenticate
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return &EncryptedData{
		Version:    KeyVersion,
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}, nil
}

// Decrypt decrypts encrypted data using AES-256-GCM.
//
// Verifies the authentication tag to ensure data integrity and authenticity.
// Returns an error if the data has been tampered with.
func Decrypt(encrypted *EncryptedData) ([]byte, error) {
	// Validate input
	if encrypted == nil {
		return nil, fmt.Errorf("encrypted data is nil")
	}
	if len(encrypted.Nonce) != NonceSize {
		return nil, fmt.Errorf("invalid nonce size: %d (expected %d)", len(encrypted.Nonce), NonceSize)
	}

	// Get encryption key
	// TODO: For key rotation, we'd look up the key by version
	// For now, we only support the current version
	if encrypted.Version != KeyVersion {
		return nil, fmt.Errorf("unsupported encryption version: %s (current: %s)", encrypted.Version, KeyVersion)
	}

	key, err := GetOrCreateEncryptionKey()
	if err != nil {
		return nil, err
	}

	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	// Decrypt and verify authentication tag
	plaintext, err := gcm.Open(nil, encrypted.Nonce, encrypted.Ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt and verify: %w (data may be corrupted or tampered with)", err)
	}

	return plaintext, nil
}

// DeriveKeyFromPassphrase derives a 256-bit encryption key from a passphrase.
//
// This is a fallback method when the system keyring is not available.
// Uses SHA-256 to derive a key from the passphrase.
//
// WARNING: This is less secure than keyring storage because:
//   - The passphrase must be entered each time
//   - The passphrase may be visible in memory
//   - No protection if the system is compromised
//
// Only use this as a last resort when keyring is unavailable.
func DeriveKeyFromPassphrase(passphrase string) []byte {
	// Use SHA-256 to derive a 256-bit key
	hash := sha256.Sum256([]byte(passphrase))
	return hash[:]
}

// RotateEncryptionKey generates a new encryption key and re-encrypts all tokens.
//
// This is used for key rotation to maintain security hygiene.
// The old key is kept temporarily to decrypt existing tokens,
// then all tokens are re-encrypted with the new key.
//
// Key rotation should be performed periodically (e.g., every 90 days).
func RotateEncryptionKey() error {
	// TODO: Implement key rotation
	// Steps:
	// 1. Generate new key with new version
	// 2. Store new key in keyring with versioned identifier
	// 3. List all token files
	// 4. For each token:
	//    a. Decrypt with old key
	//    b. Re-encrypt with new key
	//    c. Save updated token
	// 5. Mark old key as deprecated (but keep for rollback)

	return fmt.Errorf("key rotation not yet implemented")
}

// ClearEncryptionKey removes the encryption key from the system keyring.
//
// WARNING: This will make all encrypted tokens unreadable!
// Only use this if you're sure you want to delete all encrypted data.
func ClearEncryptionKey() error {
	if err := keyring.Delete(KeyringService, KeyringUser); err != nil {
		if err == keyring.ErrNotFound {
			return nil // Already deleted
		}
		return fmt.Errorf("delete encryption key from keyring: %w", err)
	}
	return nil
}
