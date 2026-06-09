package utils

import (
	"errors"
	"strings"
)

// NormalizeEmail lower-cases and trims an email address. The local part's
// "+tag" suffix is intentionally preserved — providers vary on whether
// tags are equivalent to the base address, so we store what the user gave.
func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// NormalizeDomain lower-cases and trims a domain, also stripping a leading
// "@" or "." that may sneak in from user input or copy-paste.
func NormalizeDomain(s string) string {
	d := strings.ToLower(strings.TrimSpace(s))
	d = strings.TrimPrefix(d, "@")
	d = strings.TrimPrefix(d, ".")
	d = strings.TrimSuffix(d, ".")
	return d
}

// Slugify produces a URL-safe slug from a free-form name plus a short
// random suffix to avoid collisions across signups (we'd rather not
// surface "name already taken" errors during org creation).
func Slugify(name string) string {
	var b strings.Builder
	prevDash := true
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		case prevDash:
			// skip
		default:
			b.WriteRune('-')
			prevDash = true
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		slug = "org"
	}
	suffix, err := RandomBytes(3)
	if err == nil {
		slug += "-" + base32Lower(suffix)
	}
	return slug
}

func base32Lower(b []byte) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz234567"
	out := make([]byte, 0, len(b)*8/5+1)
	var bits, n uint
	for _, x := range b {
		bits = (bits << 8) | uint(x)
		n += 8
		for n >= 5 {
			n -= 5
			out = append(out, alphabet[(bits>>n)&31])
		}
	}
	if n > 0 {
		out = append(out, alphabet[(bits<<(5-n))&31])
	}
	return string(out)
}

// ExtractDomain returns the normalized domain portion of an email. The
// email itself does NOT need to be pre-normalized; this function lower-
// cases and trims as part of extraction. Returns an error if there is
// no `@` or the domain is empty.
func ExtractDomain(email string) (string, error) {
	at := strings.LastIndex(email, "@")
	if at < 0 {
		return "", errors.New("email missing '@'")
	}
	d := NormalizeDomain(email[at+1:])
	if d == "" {
		return "", errors.New("email has empty domain")
	}
	return d, nil
}
