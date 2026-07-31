package service

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestLooksLikeAPIKey(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"minted key", mustNewAPIKeyToken(t), true},
		{"bare prefix", apiKeyPrefix, true},
		{"jwt", "eyJhbGciOiJIUzI1NiJ9.eyJhIjoxfQ.sig", false},
		{"empty", "", false},
		{"prefix in the middle", "Bearer wgl_abc", false},
		{"wrong case", "WGL_abc", false},
		{"near miss", "wgl-abc", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := LooksLikeAPIKey(c.in); got != c.want {
				t.Fatalf("LooksLikeAPIKey(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func mustNewAPIKeyToken(t *testing.T) string {
	t.Helper()
	tok, err := newAPIKeyToken()
	if err != nil {
		t.Fatalf("newAPIKeyToken: %v", err)
	}
	return tok
}

func TestNewAPIKeyTokenShape(t *testing.T) {
	tok := mustNewAPIKeyToken(t)

	if !strings.HasPrefix(tok, apiKeyPrefix) {
		t.Fatalf("token %q lacks the %q prefix", tok, apiKeyPrefix)
	}
	body := strings.TrimPrefix(tok, apiKeyPrefix)

	// 32 random bytes, raw-url base64 => 43 characters.
	if len(body) != 43 {
		t.Fatalf("token body is %d chars, want 43 (32 bytes raw-url base64)", len(body))
	}
	raw, err := base64.RawURLEncoding.DecodeString(body)
	if err != nil {
		t.Fatalf("token body is not valid raw-url base64: %v", err)
	}
	if len(raw) != 32 {
		t.Fatalf("token carries %d bytes of entropy, want 32", len(raw))
	}

	// The Prefix column slices the token directly, so it must be long enough.
	if len(tok) < apiKeyDisplayLen {
		t.Fatalf("token is %d chars, shorter than apiKeyDisplayLen %d — Issue would panic slicing it",
			len(tok), apiKeyDisplayLen)
	}
}

func TestNewAPIKeyTokenIsURLSafe(t *testing.T) {
	// Raw-url base64 must never emit '+', '/' or '=' padding, which would
	// break the token when passed in headers or URLs.
	for i := 0; i < 50; i++ {
		tok := mustNewAPIKeyToken(t)
		if strings.ContainsAny(tok, "+/=") {
			t.Fatalf("token %q contains a URL-unsafe character", tok)
		}
	}
}

func TestNewAPIKeyTokensAreUnique(t *testing.T) {
	const n = 200
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		tok := mustNewAPIKeyToken(t)
		if _, dup := seen[tok]; dup {
			t.Fatalf("newAPIKeyToken produced a duplicate within %d calls", n)
		}
		seen[tok] = struct{}{}
	}
}

func TestHashAPIKey(t *testing.T) {
	tok := mustNewAPIKeyToken(t)
	h := hashAPIKey(tok)

	if len(h) != 64 {
		t.Fatalf("hash is %d chars, want 64 (hex sha256)", len(h))
	}
	if h != hashAPIKey(tok) {
		t.Fatal("hashAPIKey is not deterministic; lookup by hash would never match")
	}
	if strings.Contains(h, tok) || strings.Contains(h, strings.TrimPrefix(tok, apiKeyPrefix)) {
		t.Fatal("the hash leaks the plaintext token")
	}
	if h == hashAPIKey(tok+"x") {
		t.Fatal("hashAPIKey collided on different inputs")
	}
	for _, r := range h {
		if !strings.ContainsRune("0123456789abcdef", r) {
			t.Fatalf("hash %q contains non-hex rune %q", h, r)
		}
	}
}

// The Prefix stored for display must be recognisable but must not reveal
// enough of the token to be useful.
func TestAPIKeyDisplayPrefixIsSafe(t *testing.T) {
	tok := mustNewAPIKeyToken(t)
	prefix := tok[:apiKeyDisplayLen]

	if !strings.HasPrefix(prefix, apiKeyPrefix) {
		t.Fatalf("display prefix %q should start with %q so keys are identifiable", prefix, apiKeyPrefix)
	}
	secretChars := apiKeyDisplayLen - len(apiKeyPrefix)
	if secretChars > 12 {
		t.Fatalf("display prefix reveals %d characters of the secret; that is too many", secretChars)
	}
	if len(prefix) >= len(tok) {
		t.Fatal("the display prefix is the whole token")
	}
}
