package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/cncmarella/cleanleaf/api/internal/mailer"
)

type ContactInput struct {
	Name    string
	Email   string
	Phone   string
	Subject string
	Message string

	// Website is a honeypot: hidden in the form, so a human leaves it empty and
	// a naive bot fills it in. Populated submissions are dropped silently.
	Website string
}

// ValidationError carries per-field messages for an invalid enquiry, keyed by field name so the handler can return them as structured JSON.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("contact enquiry failed validation on %d field(s)", len(e.Fields))
}

// DeliveryError wraps a failure to hand the enquiry to the mail transport. The cause stays inside for logging and never reaches the visitor.
type DeliveryError struct{ Cause error }

func (e *DeliveryError) Error() string {
	return "contact enquiry delivery failed: " + e.Cause.Error()
}

func (e *DeliveryError) Unwrap() error { return e.Cause }

// Contact validates enquiries and hands them to the mail transport.
type Contact struct {
	mailer   mailer.Mailer
	log      *slog.Logger
	mailFrom string
	mailTo   string
}

// NewContact builds a Contact service.
func NewContact(m mailer.Mailer, log *slog.Logger, mailFrom, mailTo string) *Contact {
	return &Contact{mailer: m, log: log, mailFrom: mailFrom, mailTo: mailTo}
}

// Submit normalises, validates and delivers an enquiry. It returns:
//
//   - *ValidationError when a field is missing or malformed;
//   - *DeliveryError when the mail transport fails;
//   - nil on success, which includes a honeypot hit — silently dropped so a bot
//     gets exactly the outcome a person does.
func (c *Contact) Submit(ctx context.Context, in ContactInput) error {
	// A filled honeypot is spam. Report success without sending anything, so the bot has no signal that it was rejected.
	if strings.TrimSpace(in.Website) != "" {
		c.log.Info("contact honeypot triggered")
		return nil
	}

	in.normalise()
	if fields := in.validate(); len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}

	msg := mailer.Message{
		From:    c.mailFrom,
		To:      c.mailTo,
		ReplyTo: in.Email,
		Subject: fmt.Sprintf("Website enquiry: %s", in.Subject),
		Text:    in.body(),
	}
	if err := c.mailer.Send(ctx, msg); err != nil {
		// The visitor gets a generic message; the detail stays in our logs.
		c.log.Error("send contact email", "error", err)
		return &DeliveryError{Cause: err}
	}

	c.log.Info("contact enquiry accepted", "subject", in.Subject)
	return nil
}

func (in *ContactInput) normalise() {
	in.Name = strings.TrimSpace(in.Name)
	in.Email = strings.TrimSpace(in.Email)
	in.Phone = strings.TrimSpace(in.Phone)
	in.Subject = strings.TrimSpace(in.Subject)
	in.Message = strings.TrimSpace(in.Message)
	if in.Subject == "" {
		in.Subject = "General Enquiry"
	}
}

// validate returns a field-name -> message map, empty when the request is good.
func (in *ContactInput) validate() map[string]string {
	fields := make(map[string]string)

	switch n := utf8.RuneCountInString(in.Name); {
	case n == 0:
		fields["name"] = "Please tell us your name."
	case n > 100:
		fields["name"] = "Please keep your name under 100 characters."
	}

	if in.Email == "" {
		fields["email"] = "Please give us an email address to reply to."
	} else if _, err := mail.ParseAddress(in.Email); err != nil {
		fields["email"] = "That does not look like a valid email address."
	}

	// Phone is optional, but if given it should be plausible.
	if in.Phone != "" {
		digits := 0
		for _, r := range in.Phone {
			if r >= '0' && r <= '9' {
				digits++
			}
		}
		if digits < 7 || digits > 15 {
			fields["phone"] = "Please enter a valid phone number"
		}
	}

	if n := utf8.RuneCountInString(in.Subject); n > 150 {
		fields["subject"] = "Please keep the subject under 150 characters."
	}

	switch n := utf8.RuneCountInString(in.Message); {
	case n < 10:
		fields["message"] = "Please give us a little more detail (at least 10 characters)."
	case n > 4000:
		fields["message"] = "Please keep your message under 4000 characters."
	}

	return fields
}

// body renders the plain-text email sent to the sales inbox.
func (in *ContactInput) body() string {
	phone := in.Phone
	if phone == "" {
		phone = "(not provided)"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "New enquiry from the CleanLeaf website.\n\n")
	fmt.Fprintf(&b, "Name:    %s\n", in.Name)
	fmt.Fprintf(&b, "Email:   %s\n", in.Email)
	fmt.Fprintf(&b, "Phone:   %s\n", phone)
	fmt.Fprintf(&b, "Subject: %s\n\n", in.Subject)
	fmt.Fprintf(&b, "Message:\n%s\n", in.Message)
	return b.String()
}
