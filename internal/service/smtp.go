package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// SMTPSender delivers transactional email over SMTP. It negotiates STARTTLS and
// AUTH only when the server advertises them, so it works against a bare dev
// relay (Mailpit) and a TLS+auth production relay without configuration churn.
type SMTPSender struct {
	host     string
	port     int
	user     string
	password string
	from     string
	baseURL  string
}

func NewSMTPSender(host string, port int, user, password, from, baseURL string) *SMTPSender {
	if strings.TrimSpace(from) == "" {
		from = user
	}
	return &SMTPSender{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
		baseURL:  strings.TrimRight(baseURL, "/"),
	}
}

func (s *SMTPSender) SendVerification(ctx context.Context, email, token string) error {
	link := s.baseURL + "/?verify=" + token
	body := htmlEmail(
		"Verify your email",
		"Confirm your email address to finish setting up your Waggle account.",
		"Verify email", link,
	)
	return s.send(ctx, email, "Verify your Waggle email", body)
}

func (s *SMTPSender) SendInvite(ctx context.Context, email, orgName, token string) error {
	link := s.baseURL + "/?invite=" + token
	body := htmlEmail(
		"You've been invited to "+orgName,
		fmt.Sprintf("You've been invited to join %q on Waggle. Accept to set a password and get started.", orgName),
		"Accept invitation", link,
	)
	return s.send(ctx, email, "Invitation to "+orgName+" on Waggle", body)
}

// send builds an RFC 5322 HTML message and delivers it, opting into STARTTLS /
// AUTH per the server's advertised capabilities.
func (s *SMTPSender) send(ctx context.Context, to, subject, htmlBody string) error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	d := net.Dialer{Timeout: 10 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}

	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("smtp client: %w", err)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: s.host}); err != nil {
			return fmt.Errorf("smtp starttls: %w", err)
		}
	}
	if ok, _ := c.Extension("AUTH"); ok && s.password != "" {
		if err := c.Auth(smtp.PlainAuth("", s.user, s.password, s.host)); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := c.Mail(s.from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	msg := buildMessage(s.from, to, subject, htmlBody)
	if _, err := w.Write([]byte(msg)); err != nil {
		_ = w.Close()
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}
	return c.Quit()
}

func buildMessage(from, to, subject, htmlBody string) string {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + subject + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlBody)
	return b.String()
}

func htmlEmail(heading, intro, cta, link string) string {
	return fmt.Sprintf(`<!doctype html><html><body style="font-family:system-ui,sans-serif;background:#f4f4f5;padding:24px;">
  <table style="max-width:480px;margin:0 auto;background:#fff;border-radius:12px;padding:32px;">
    <tr><td>
      <h1 style="font-size:18px;margin:0 0 12px;">%s</h1>
      <p style="color:#52525b;font-size:14px;line-height:1.5;margin:0 0 24px;">%s</p>
      <a href="%s" style="display:inline-block;background:#18181b;color:#fff;text-decoration:none;padding:10px 18px;border-radius:8px;font-size:14px;">%s</a>
      <p style="color:#a1a1aa;font-size:12px;margin:24px 0 0;word-break:break-all;">Or paste this link:<br>%s</p>
    </td></tr>
  </table>
</body></html>`, heading, intro, link, cta, link)
}
