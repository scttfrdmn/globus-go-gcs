package auth

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
	}{
		{
			name:      "empty data",
			plaintext: []byte{},
		},
		{
			name:      "small text",
			plaintext: []byte("hello world"),
		},
		{
			name:      "json data",
			plaintext: []byte(`{"access_token": "secret123", "refresh_token": "refresh456"}`),
		},
		{
			name:      "large data",
			plaintext: make([]byte, 10000), // 10KB of zeros
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt() error = %v", err)
			}

			// Verify encrypted data structure
			if encrypted.Version != KeyVersion {
				t.Errorf("Version = %v, want %v", encrypted.Version, KeyVersion)
			}
			if len(encrypted.Nonce) != NonceSize {
				t.Errorf("Nonce length = %d, want %d", len(encrypted.Nonce), NonceSize)
			}
			if len(encrypted.Ciphertext) == 0 {
				t.Error("Ciphertext is empty")
			}

			// Verify ciphertext is different from plaintext
			if len(tt.plaintext) > 0 && bytes.Equal(encrypted.Ciphertext[:len(tt.plaintext)], tt.plaintext) {
				t.Error("Ciphertext appears to be unencrypted")
			}

			// Decrypt
			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt() error = %v", err)
			}

			// Verify decrypted matches original
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Decrypted data doesn't match original\ngot:  %s\nwant: %s", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptDecrypt_RandomData(t *testing.T) {
	// Test with 100 rounds of random data
	for i := 0; i < 100; i++ {
		// Generate random plaintext (1-1000 bytes)
		size := 1 + (i * 10)
		if size > 1000 {
			size = 1000
		}
		plaintext := make([]byte, size)
		if _, err := rand.Read(plaintext); err != nil {
			t.Fatalf("Generate random data: %v", err)
		}

		// Encrypt and decrypt
		encrypted, err := Encrypt(plaintext)
		if err != nil {
			t.Fatalf("Round %d: Encrypt() error = %v", i, err)
		}

		decrypted, err := Decrypt(encrypted)
		if err != nil {
			t.Fatalf("Round %d: Decrypt() error = %v", i, err)
		}

		if !bytes.Equal(decrypted, plaintext) {
			t.Errorf("Round %d: Decrypted data doesn't match", i)
		}
	}
}

func TestDecrypt_InvalidData(t *testing.T) {
	tests := []struct {
		name      string
		encrypted *EncryptedData
		wantErr   bool
	}{
		{
			name:      "nil data",
			encrypted: nil,
			wantErr:   true,
		},
		{
			name: "invalid nonce size",
			encrypted: &EncryptedData{
				Version:    KeyVersion,
				Nonce:      []byte{1, 2, 3}, // Too short
				Ciphertext: []byte{1, 2, 3},
			},
			wantErr: true,
		},
		{
			name: "unsupported version",
			encrypted: &EncryptedData{
				Version:    "v99",
				Nonce:      make([]byte, NonceSize),
				Ciphertext: []byte{1, 2, 3},
			},
			wantErr: true,
		},
		{
			name: "corrupted ciphertext",
			encrypted: &EncryptedData{
				Version:    KeyVersion,
				Nonce:      make([]byte, NonceSize),
				Ciphertext: []byte{1, 2, 3, 4, 5}, // Invalid ciphertext
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncrypt_DifferentNonces(t *testing.T) {
	// Verify that encrypting the same plaintext multiple times
	// produces different ciphertexts (due to different nonces)
	plaintext := []byte("test data")

	encrypted1, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	encrypted2, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Nonces should be different
	if bytes.Equal(encrypted1.Nonce, encrypted2.Nonce) {
		t.Error("Encrypt() produced same nonce twice (nonce reuse vulnerability!)")
	}

	// Ciphertexts should be different
	if bytes.Equal(encrypted1.Ciphertext, encrypted2.Ciphertext) {
		t.Error("Encrypt() produced same ciphertext twice (should differ due to nonce)")
	}

	// But both should decrypt to the same plaintext
	decrypted1, err := Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	decrypted2, err := Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}

	if !bytes.Equal(decrypted1, plaintext) || !bytes.Equal(decrypted2, plaintext) {
		t.Error("Decrypted data doesn't match original")
	}
}

func TestEncrypt_AuthenticationTag(t *testing.T) {
	plaintext := []byte("authenticated data")

	// Encrypt
	encrypted, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Tamper with ciphertext (flip a bit)
	if len(encrypted.Ciphertext) > 0 {
		encrypted.Ciphertext[0] ^= 1
	}

	// Decryption should fail due to authentication failure
	_, err = Decrypt(encrypted)
	if err == nil {
		t.Error("Decrypt() should fail on tampered ciphertext, but succeeded")
	}
}

func TestDeriveKeyFromPassphrase(t *testing.T) {
	tests := []struct {
		name       string
		passphrase string
	}{
		{
			name:       "simple passphrase",
			passphrase: "password123",
		},
		{
			name:       "complex passphrase",
			passphrase: "MyS3cur3P@ssw0rd!2025",
		},
		{
			name:       "unicode passphrase",
			passphrase: "пароль密码パスワード",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := DeriveKeyFromPassphrase(tt.passphrase)

			// Verify key size
			if len(key) != EncryptionKeySize {
				t.Errorf("Key length = %d, want %d", len(key), EncryptionKeySize)
			}

			// Verify deterministic (same passphrase = same key)
			key2 := DeriveKeyFromPassphrase(tt.passphrase)
			if !bytes.Equal(key, key2) {
				t.Error("DeriveKeyFromPassphrase() not deterministic")
			}

			// Verify different passphrases produce different keys
			differentKey := DeriveKeyFromPassphrase(tt.passphrase + "x")
			if bytes.Equal(key, differentKey) {
				t.Error("Different passphrases produced same key")
			}
		})
	}
}

func BenchmarkEncrypt(b *testing.B) {
	plaintext := []byte(`{
		"access_token": "very_long_access_token_string_here",
		"refresh_token": "very_long_refresh_token_string_here",
		"expires_at": "2025-12-31T23:59:59Z",
		"scopes": ["scope1", "scope2", "scope3"]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Encrypt(plaintext)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	plaintext := []byte(`{
		"access_token": "very_long_access_token_string_here",
		"refresh_token": "very_long_refresh_token_string_here",
		"expires_at": "2025-12-31T23:59:59Z",
		"scopes": ["scope1", "scope2", "scope3"]
	}`)

	encrypted, err := Encrypt(plaintext)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Decrypt(encrypted)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	plaintext := []byte(`{
		"access_token": "very_long_access_token_string_here",
		"refresh_token": "very_long_refresh_token_string_here",
		"expires_at": "2025-12-31T23:59:59Z",
		"scopes": ["scope1", "scope2", "scope3"]
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypted, err := Encrypt(plaintext)
		if err != nil {
			b.Fatal(err)
		}

		_, err = Decrypt(encrypted)
		if err != nil {
			b.Fatal(err)
		}
	}
}
