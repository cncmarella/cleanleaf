package mailer

import (
	"context"
	"log/slog"
)

// Console writes messages to the log instead of sending them. It is the
// development default so the contact form is testable without any API key.
type Console struct {
	Log *slog.Logger
}

// NewConsole returns a Console mailer.
func NewConsole(log *slog.Logger) *Console { return &Console{Log: log} }

// Send logs the message and always succeeds.
func (c *Console) Send(_ context.Context, m Message) error {
	c.Log.Info("email not sent (console mailer)",
		"to", m.To,
		"reply_to", m.ReplyTo,
		"subject", m.Subject,
		"body", m.Text,
	)
	return nil
}
