package service

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/glueops/waggle/internal/utils"

	"github.com/google/uuid"
)

// testSigningKey returns a fresh random signing key. Generated rather than
// hard-coded so every run uses a distinct key and no secret-shaped literal
// sits in the source for scanners to flag.
func testSigningKey(t *testing.T) string {
	t.Helper()
	b, err := utils.RandomBytes(32)
	if err != nil {
		t.Fatalf("RandomBytes: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func newTestTokenService(t *testing.T) *TokenService {
	t.Helper()
	return newTokenServiceWith(t, testSigningKey(t), "waggle-test", time.Hour, 24*time.Hour)
}

func newTokenServiceWith(t *testing.T, key, issuer string, accessTTL, refreshTTL time.Duration) *TokenService {
	t.Helper()
	ts, err := NewTokenService(key, issuer, accessTTL, refreshTTL, "waggle")
	if err != nil {
		t.Fatalf("NewTokenService: %v", err)
	}
	return ts
}

func TestNewTokenServiceRequiresSecret(t *testing.T) {
	if _, err := NewTokenService("", "iss", time.Hour, time.Hour, "aud"); err == nil {
		t.Fatal("NewTokenService accepted an empty secret")
	}
}

func TestNewTokenServiceDefaultsAudience(t *testing.T) {
	ts, err := NewTokenService(testSigningKey(t), "iss", time.Hour, time.Hour, "")
	if err != nil {
		t.Fatalf("NewTokenService: %v", err)
	}
	if ts.audience != "waggle" {
		t.Fatalf("audience = %q, want the %q default", ts.audience, "waggle")
	}
}

func TestIssueAndVerifyAccess(t *testing.T) {
	ts := newTestTokenService(t)
	accountID, orgID, sessionID := uuid.New(), uuid.New(), uuid.New()

	tok, exp, err := ts.IssueAccess(accountID, orgID, sessionID)
	if err != nil {
		t.Fatalf("IssueAccess: %v", err)
	}
	if tok == "" {
		t.Fatal("IssueAccess returned an empty token")
	}
	if !exp.After(time.Now().UTC()) {
		t.Fatalf("expiry %s is not in the future", exp)
	}

	claims, err := ts.Verify(tok)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if claims.AccountID != accountID {
		t.Fatalf("AccountID = %s, want %s", claims.AccountID, accountID)
	}
	if claims.OrganizationID != orgID {
		t.Fatalf("OrganizationID = %s, want %s", claims.OrganizationID, orgID)
	}
	if claims.SessionID != sessionID {
		t.Fatalf("SessionID = %s, want %s", claims.SessionID, sessionID)
	}
	if claims.Subject != accountID.String() {
		t.Fatalf("Subject = %q, want %q", claims.Subject, accountID)
	}
}

func TestIssueAccessRequiresIDs(t *testing.T) {
	ts := newTestTokenService(t)

	if _, _, err := ts.IssueAccess(uuid.Nil, uuid.New(), uuid.New()); err == nil {
		t.Fatal("IssueAccess accepted a nil account id")
	}
	if _, _, err := ts.IssueAccess(uuid.New(), uuid.New(), uuid.Nil); err == nil {
		t.Fatal("IssueAccess accepted a nil session id")
	}
	// A nil org id is legitimate: "no org context selected".
	if _, _, err := ts.IssueAccess(uuid.New(), uuid.Nil, uuid.New()); err != nil {
		t.Fatalf("IssueAccess rejected a nil org id, which should mean no org context: %v", err)
	}
}

func TestVerifyRejectsBadTokens(t *testing.T) {
	ts := newTestTokenService(t)
	good, _, err := ts.IssueAccess(uuid.New(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("IssueAccess: %v", err)
	}

	t.Run("empty", func(t *testing.T) {
		if _, err := ts.Verify(""); err == nil {
			t.Fatal("Verify accepted an empty token")
		}
	})

	t.Run("garbage", func(t *testing.T) {
		if _, err := ts.Verify("not.a.jwt"); err == nil {
			t.Fatal("Verify accepted garbage")
		}
	})

	t.Run("tampered payload", func(t *testing.T) {
		parts := strings.Split(good, ".")
		if len(parts) != 3 {
			t.Fatalf("expected a 3-part JWT, got %d parts", len(parts))
		}
		tampered := parts[0] + "." + parts[1] + "x." + parts[2]
		if _, err := ts.Verify(tampered); err == nil {
			t.Fatal("Verify accepted a token with a tampered payload")
		}
	})

	t.Run("signed with a different key", func(t *testing.T) {
		other := newTokenServiceWith(t, testSigningKey(t), "waggle-test", time.Hour, time.Hour)
		foreign, _, err := other.IssueAccess(uuid.New(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("IssueAccess: %v", err)
		}
		if _, err := ts.Verify(foreign); err == nil {
			t.Fatal("Verify accepted a token signed with another key")
		}
	})

	t.Run("wrong issuer", func(t *testing.T) {
		other := newTokenServiceWith(t, testSigningKey(t), "someone-else", time.Hour, time.Hour)
		foreign, _, err := other.IssueAccess(uuid.New(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("IssueAccess: %v", err)
		}
		if _, err := ts.Verify(foreign); err == nil {
			t.Fatal("Verify accepted a token from a different issuer")
		}
	})

	t.Run("expired", func(t *testing.T) {
		expired := newTokenServiceWith(t, testSigningKey(t), "waggle-test", -time.Hour, time.Hour)
		tok, _, err := expired.IssueAccess(uuid.New(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("IssueAccess: %v", err)
		}
		if _, err := expired.Verify(tok); err == nil {
			t.Fatal("Verify accepted an expired token")
		}
	})

	t.Run("unsigned alg=none", func(t *testing.T) {
		// The classic JWT downgrade: "alg":"none" with an empty signature.
		// Assembled at runtime rather than embedded as a literal, so secret
		// scanners don't flag the test fixture as a leaked token.
		seg := func(s string) string {
			return base64.RawURLEncoding.EncodeToString([]byte(s))
		}
		none := seg(`{"alg":"none","typ":"JWT"}`) + "." +
			seg(`{"account_id":"`+uuid.New().String()+`"}`) + "."
		if _, err := ts.Verify(none); err == nil {
			t.Fatal("Verify accepted an alg=none token")
		}
	})
}

// The three token kinds use distinct audiences precisely so one cannot be
// replayed as another. This is the security property worth locking down.
func TestTokenAudiencesAreNotInterchangeable(t *testing.T) {
	ts := newTestTokenService(t)
	accountID, orgID, sessionID, emailID := uuid.New(), uuid.New(), uuid.New(), uuid.New()

	access, _, err := ts.IssueAccess(accountID, orgID, sessionID)
	if err != nil {
		t.Fatalf("IssueAccess: %v", err)
	}
	verify, _, err := ts.IssueEmailVerification(emailID)
	if err != nil {
		t.Fatalf("IssueEmailVerification: %v", err)
	}
	invite, _, err := ts.IssueInvite(accountID, orgID)
	if err != nil {
		t.Fatalf("IssueInvite: %v", err)
	}

	t.Run("access token is not an email-verify token", func(t *testing.T) {
		if _, err := ts.VerifyEmailVerification(access); err == nil {
			t.Fatal("an access token was accepted as an email-verification token")
		}
	})
	t.Run("access token is not an invite", func(t *testing.T) {
		if _, _, err := ts.VerifyInvite(access); err == nil {
			t.Fatal("an access token was accepted as an invite token")
		}
	})
	t.Run("email-verify token is not an access token", func(t *testing.T) {
		if _, err := ts.Verify(verify); err == nil {
			t.Fatal("an email-verification token was accepted as an access token")
		}
	})
	t.Run("email-verify token is not an invite", func(t *testing.T) {
		if _, _, err := ts.VerifyInvite(verify); err == nil {
			t.Fatal("an email-verification token was accepted as an invite token")
		}
	})
	t.Run("invite is not an access token", func(t *testing.T) {
		if _, err := ts.Verify(invite); err == nil {
			t.Fatal("an invite token was accepted as an access token")
		}
	})
	t.Run("invite is not an email-verify token", func(t *testing.T) {
		if _, err := ts.VerifyEmailVerification(invite); err == nil {
			t.Fatal("an invite token was accepted as an email-verification token")
		}
	})
}

func TestEmailVerificationRoundTrip(t *testing.T) {
	ts := newTestTokenService(t)
	emailID := uuid.New()

	tok, exp, err := ts.IssueEmailVerification(emailID)
	if err != nil {
		t.Fatalf("IssueEmailVerification: %v", err)
	}
	if got := time.Until(exp).Round(time.Minute); got != EmailVerifyTTL {
		t.Fatalf("expiry is %s away, want %s", got, EmailVerifyTTL)
	}

	got, err := ts.VerifyEmailVerification(tok)
	if err != nil {
		t.Fatalf("VerifyEmailVerification: %v", err)
	}
	if got != emailID {
		t.Fatalf("AccountEmailID = %s, want %s", got, emailID)
	}

	if _, _, err := ts.IssueEmailVerification(uuid.Nil); err == nil {
		t.Fatal("IssueEmailVerification accepted a nil id")
	}
}

func TestInviteRoundTrip(t *testing.T) {
	ts := newTestTokenService(t)
	accountID, orgID := uuid.New(), uuid.New()

	tok, exp, err := ts.IssueInvite(accountID, orgID)
	if err != nil {
		t.Fatalf("IssueInvite: %v", err)
	}
	if got := time.Until(exp).Round(time.Minute); got != InviteTTL {
		t.Fatalf("expiry is %s away, want %s", got, InviteTTL)
	}

	gotAccount, gotOrg, err := ts.VerifyInvite(tok)
	if err != nil {
		t.Fatalf("VerifyInvite: %v", err)
	}
	if gotAccount != accountID || gotOrg != orgID {
		t.Fatalf("VerifyInvite = (%s, %s), want (%s, %s)", gotAccount, gotOrg, accountID, orgID)
	}

	if _, _, err := ts.IssueInvite(uuid.Nil, orgID); err == nil {
		t.Fatal("IssueInvite accepted a nil account id")
	}
	if _, _, err := ts.IssueInvite(accountID, uuid.Nil); err == nil {
		t.Fatal("IssueInvite accepted a nil org id")
	}
}

func TestGenerateRefresh(t *testing.T) {
	ts := newTestTokenService(t)

	a, err := ts.GenerateRefresh()
	if err != nil {
		t.Fatalf("GenerateRefresh: %v", err)
	}
	if a.Plain == "" || a.Hashed == "" {
		t.Fatal("GenerateRefresh returned an empty token or hash")
	}
	if a.Plain == a.Hashed {
		t.Fatal("the stored hash equals the plaintext refresh token")
	}
	if a.Hashed != HashRefreshToken(a.Plain) {
		t.Fatal("Hashed does not match HashRefreshToken(Plain); lookup would never match")
	}
	if !a.ExpiresAt.After(time.Now().UTC()) {
		t.Fatalf("refresh expiry %s is not in the future", a.ExpiresAt)
	}

	b, err := ts.GenerateRefresh()
	if err != nil {
		t.Fatalf("GenerateRefresh: %v", err)
	}
	if a.Plain == b.Plain {
		t.Fatal("two refresh tokens were identical")
	}
}

func TestHashRefreshTokenIsStable(t *testing.T) {
	const tok = "some-opaque-refresh-token"
	if HashRefreshToken(tok) != HashRefreshToken(tok) {
		t.Fatal("HashRefreshToken is not deterministic; session lookup would break")
	}
	if HashRefreshToken(tok) == HashRefreshToken(tok+"x") {
		t.Fatal("HashRefreshToken collided on different inputs")
	}
	if strings.Contains(HashRefreshToken(tok), tok) {
		t.Fatal("the hash leaks the plaintext token")
	}
}
