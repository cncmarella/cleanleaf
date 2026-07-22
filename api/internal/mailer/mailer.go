// Package mailer delivers transactional email for the website.
package mailer

import "context"

// Message is a single outbound plain-text email.
type Message struct {
	From    string
	To      string
	ReplyTo string
	Subject string
	Text    string
}

// Mailer sends a Message. Implementations must be safe for concurrent use.
type Mailer interface {
	Send(ctx context.Context, m Message) error
}
