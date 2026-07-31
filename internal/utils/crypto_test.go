package utils

import (
	"bytes"
	"testing"
)

// aes256Key is a 32-byte key; AES-256 is what the callers use.
func aes256Key(t *testing.T) []byte {
	t.Helper()
	k, err := RandomBytes(32)
	if err != nil {
		t.Fatalf("RandomBytes: %v", err)
	}
	return k
}

func TestRandomBytesLength(t *testing.T) {
	for _, n := range []int{0, 1, 3, 16, 32, 64} {
		b, err := RandomBytes(n)
		if err != nil {
			t.Fatalf("RandomBytes(%d): %v", n, err)
		}
		if len(b) != n {
			t.Fatalf("RandomBytes(%d) returned %d bytes", n, len(b))
		}
	}
}

func TestRandomBytesDiffer(t *testing.T) {
	a, err := RandomBytes(32)
	if err != nil {
		t.Fatalf("RandomBytes: %v", err)
	}
	b, err := RandomBytes(32)
	if err != nil {
		t.Fatalf("RandomBytes: %v", err)
	}
	if bytes.Equal(a, b) {
		t.Fatal("two RandomBytes(32) calls returned identical output")
	}
}

func TestB64RoundTrip(t *testing.T) {
	for _, in := range [][]byte{
		{}, {0x00}, {0xff, 0xfe, 0xfd}, []byte("hello world"),
	} {
		out, err := DecodeB64(EncodeB64(in))
		if err != nil {
			t.Fatalf("DecodeB64(EncodeB64(%v)): %v", in, err)
		}
		if !bytes.Equal(out, in) {
			t.Fatalf("round trip changed %v into %v", in, out)
		}
	}
}

func TestDecodeB64RejectsGarbage(t *testing.T) {
	if _, err := DecodeB64("!!!not base64!!!"); err == nil {
		t.Fatal("DecodeB64 accepted invalid base64")
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := aes256Key(t)
	cases := []struct {
		name      string
		plaintext []byte
	}{
		{"empty", []byte{}},
		{"short", []byte("hi")},
		{"typical secret", []byte("proxmox-api-token-secret-value")},
		{"binary", []byte{0x00, 0x01, 0xff, 0x7f, 0x80}},
		{"long", bytes.Repeat([]byte("waggle"), 1000)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ct, iv, tag, err := EncryptAESGCM(c.plaintext, key)
			if err != nil {
				t.Fatalf("EncryptAESGCM: %v", err)
			}
			if len(tag) != gcmTagSize {
				t.Fatalf("tag is %d bytes, want %d", len(tag), gcmTagSize)
			}
			if len(ct) != len(c.plaintext) {
				t.Fatalf("ciphertext is %d bytes, want %d (GCM is a stream mode)", len(ct), len(c.plaintext))
			}
			if len(c.plaintext) > 0 && bytes.Equal(ct, c.plaintext) {
				t.Fatal("ciphertext equals plaintext")
			}

			got, err := DecryptAESGCM(ct, key, iv, tag)
			if err != nil {
				t.Fatalf("DecryptAESGCM: %v", err)
			}
			if !bytes.Equal(got, c.plaintext) {
				t.Fatalf("round trip changed %q into %q", c.plaintext, got)
			}
		})
	}
}

// Every encryption must use a fresh nonce — reusing one under the same key is
// a catastrophic GCM failure, so this is worth pinning.
func TestEncryptUsesFreshNonce(t *testing.T) {
	key := aes256Key(t)
	pt := []byte("same plaintext every time")

	ct1, iv1, _, err := EncryptAESGCM(pt, key)
	if err != nil {
		t.Fatalf("EncryptAESGCM: %v", err)
	}
	ct2, iv2, _, err := EncryptAESGCM(pt, key)
	if err != nil {
		t.Fatalf("EncryptAESGCM: %v", err)
	}
	if bytes.Equal(iv1, iv2) {
		t.Fatal("nonce reused across two encryptions with the same key")
	}
	if bytes.Equal(ct1, ct2) {
		t.Fatal("identical plaintext produced identical ciphertext — nonce is not being applied")
	}
}

func TestDecryptRejectsTamperedCiphertext(t *testing.T) {
	key := aes256Key(t)
	pt := []byte("authenticated data must not be malleable")

	ct, iv, tag, err := EncryptAESGCM(pt, key)
	if err != nil {
		t.Fatalf("EncryptAESGCM: %v", err)
	}

	t.Run("flipped ciphertext bit", func(t *testing.T) {
		bad := bytes.Clone(ct)
		bad[0] ^= 0x01
		if _, err := DecryptAESGCM(bad, key, iv, tag); err == nil {
			t.Fatal("decrypt accepted tampered ciphertext")
		}
	})

	t.Run("flipped tag bit", func(t *testing.T) {
		bad := bytes.Clone(tag)
		bad[0] ^= 0x01
		if _, err := DecryptAESGCM(ct, key, iv, bad); err == nil {
			t.Fatal("decrypt accepted tampered auth tag")
		}
	})

	t.Run("flipped nonce bit", func(t *testing.T) {
		bad := bytes.Clone(iv)
		bad[0] ^= 0x01
		if _, err := DecryptAESGCM(ct, key, bad, tag); err == nil {
			t.Fatal("decrypt accepted mismatched nonce")
		}
	})

	t.Run("wrong key", func(t *testing.T) {
		if _, err := DecryptAESGCM(ct, aes256Key(t), iv, tag); err == nil {
			t.Fatal("decrypt accepted the wrong key")
		}
	})
}

func TestDecryptRejectsWrongNonceSize(t *testing.T) {
	key := aes256Key(t)
	ct, iv, tag, err := EncryptAESGCM([]byte("x"), key)
	if err != nil {
		t.Fatalf("EncryptAESGCM: %v", err)
	}
	if _, err := DecryptAESGCM(ct, key, append(bytes.Clone(iv), 0x00), tag); err == nil {
		t.Fatal("decrypt accepted an over-long nonce")
	}
	if _, err := DecryptAESGCM(ct, key, iv[:len(iv)-1], tag); err == nil {
		t.Fatal("decrypt accepted a truncated nonce")
	}
}

func TestEncryptRejectsInvalidKeySize(t *testing.T) {
	// AES accepts 16/24/32-byte keys only.
	for _, n := range []int{0, 8, 15, 17, 31, 33, 64} {
		key, err := RandomBytes(n)
		if err != nil {
			t.Fatalf("RandomBytes: %v", err)
		}
		if _, _, _, err := EncryptAESGCM([]byte("x"), key); err == nil {
			t.Fatalf("EncryptAESGCM accepted a %d-byte key", n)
		}
	}
}
