package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const gcmTagSize = 16

func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("read random: %w", err)
	}
	return b, nil
}

func EncodeB64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeB64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func EncryptAESGCM(plaintext, key []byte) (ct, iv, tag []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("aes new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gcm: %w", err)
	}

	iv = make([]byte, gcm.NonceSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, nil, fmt.Errorf("read nonce: %w", err)
	}

	sealed := gcm.Seal(nil, iv, plaintext, nil)
	if len(sealed) < gcmTagSize {
		return nil, nil, nil, fmt.Errorf("sealed output too short")
	}

	ct = sealed[:len(sealed)-gcmTagSize]
	tag = sealed[len(sealed)-gcmTagSize:]
	return ct, iv, tag, nil
}

func DecryptAESGCM(ct, key, iv, tag []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}
	if len(iv) != gcm.NonceSize() {
		return nil, fmt.Errorf("nonce size mismatch: got %d want %d", len(iv), gcm.NonceSize())
	}

	sealed := make([]byte, 0, len(ct)+len(tag))
	sealed = append(sealed, ct...)
	sealed = append(sealed, tag...)

	pt, err := gcm.Open(nil, iv, sealed, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm open: %w", err)
	}
	return pt, nil
}
