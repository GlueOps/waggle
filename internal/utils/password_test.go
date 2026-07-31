package utils

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// testPlaintext returns a random plaintext to hash. Generated rather than
// hard-coded so no credential-shaped literal sits in the source. The "aA1-"
// prefix guarantees at least one lower-case rune, so the case-sensitivity
// assertion in TestVerifyPasswordRejects can never coincide with the original.
func testPlaintext(t *testing.T) string {
	t.Helper()
	b, err := RandomBytes(12)
	if err != nil {
		t.Fatalf("RandomBytes: %v", err)
	}
	return "aA1-" + EncodeB64(b)
}

func TestHashPasswordRoundTrip(t *testing.T) {
	pw := testPlaintext(t)

	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == pw {
		t.Fatal("HashPassword returned the plaintext")
	}
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("hash %q does not look like a bcrypt hash", hash)
	}
	if !VerifyPassword(pw, hash) {
		t.Fatal("VerifyPassword rejected the correct password")
	}
}

func TestHashPasswordRejectsEmpty(t *testing.T) {
	if _, err := HashPassword(""); err == nil {
		t.Fatal("HashPassword accepted an empty password")
	}
}

// bcrypt salts every hash, so the same password must never produce the same
// digest twice.
func TestHashPasswordIsSalted(t *testing.T) {
	pw := testPlaintext(t)

	a, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	b, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if a == b {
		t.Fatal("hashing the same password twice produced identical digests — salt is missing")
	}
	if !VerifyPassword(pw, a) || !VerifyPassword(pw, b) {
		t.Fatal("both salted hashes should verify against the original password")
	}
}

func TestHashPasswordUsesConfiguredCost(t *testing.T) {
	hash, err := HashPassword(testPlaintext(t))
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("bcrypt.Cost: %v", err)
	}
	if cost != bcryptCost {
		t.Fatalf("hash cost is %d, want %d", cost, bcryptCost)
	}
	// Guard against someone lowering the constant to something weak.
	if cost < 10 {
		t.Fatalf("bcrypt cost %d is too low to be safe", cost)
	}
}

func TestVerifyPasswordRejects(t *testing.T) {
	plain := testPlaintext(t)
	good, err := HashPassword(plain)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}

	cases := []struct {
		name  string
		plain string
		hash  string
	}{
		{"wrong password", plain + "-wrong", good},
		{"empty password", "", good},
		{"empty hash", plain, ""},
		{"both empty", "", ""},
		{"malformed hash", plain, "not-a-bcrypt-hash"},
		{"case differs", strings.ToUpper(plain), good},
		{"trailing space differs", plain + " ", good},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if VerifyPassword(c.plain, c.hash) {
				t.Fatalf("VerifyPassword(%q, %q) = true, want false", c.plain, c.hash)
			}
		})
	}
}

// bcrypt caps input at 72 bytes. Go's implementation errors rather than
// silently truncating, which is the safe choice — truncation would make two
// distinct passwords sharing a 72-byte prefix interchangeable. Callers must
// therefore treat an over-length password as user error, not a 500.
func TestPasswordLongerThanBcryptLimitIsRejected(t *testing.T) {
	hash, err := HashPassword(strings.Repeat("a", 73))
	if err == nil {
		// Not the behaviour we expect today, but if a future bcrypt version
		// starts accepting long input, the truncation boundary must at least
		// not be exploitable.
		collide := strings.Repeat("a", 72) + "DIFFERENT"
		if VerifyPassword(collide, hash) {
			t.Fatal("a 72-byte-prefix collision verified — bcrypt truncation is exploitable")
		}
		return
	}
	if !strings.Contains(err.Error(), "72 bytes") {
		t.Fatalf("HashPassword(73 bytes) failed with an unexpected error: %v", err)
	}
}

// 72 bytes is the boundary and must still work.
func TestPasswordAtBcryptLimitIsAccepted(t *testing.T) {
	pw := strings.Repeat("a", 72)

	hash, err := HashPassword(pw)
	if err != nil {
		t.Fatalf("HashPassword(72 bytes): %v", err)
	}
	if !VerifyPassword(pw, hash) {
		t.Fatal("VerifyPassword rejected a valid 72-byte password")
	}
}
