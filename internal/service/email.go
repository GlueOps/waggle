package service

import (
	"context"
	"log"
)

// EmailSender delivers transactional email. Implementations are pluggable
// (LogSender for dev, an SMTP/SES/SendGrid sender for prod).
type EmailSender interface {
	SendVerification(ctx context.Context, email, token string) error
	SendInvite(ctx context.Context, email, orgName, token string) error
}

// LogSender prints verification/invite tokens to stdout. Use only in dev — the
// token is bearer-strength and should never reach a log destination that
// isn't fully trusted.
type LogSender struct{}

func (LogSender) SendVerification(_ context.Context, email, token string) error {
	log.Printf("email-sender: verification for %s — token=%s", email, token)
	return nil
}

func (LogSender) SendInvite(_ context.Context, email, orgName, token string) error {
	log.Printf("email-sender: invite to %q for %s — token=%s", orgName, email, token)
	return nil
}
