package service

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/cncmarella/cleanleaf/api/internal/mailer"
)

// stubMailer records what it was asked to send and can be made to fail.
type stubMailer struct {
	mu   sync.Mutex
	sent []mailer.Message
	err  error
}

func (s *stubMailer) Send(_ context.Context, m mailer.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err != nil {
		return s.err
	}
	s.sent = append(s.sent, m)
	return nil
}

func (s *stubMailer) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.sent)
}

func newContact(m mailer.Mailer) *Contact {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewContact(m, log, "site@cleanleaf.test", "sales@cleanleaf.test")
}

func validInput() ContactInput {
	return ContactInput{
		Name:    "Ravi Kumar",
		Email:   "ravi@example.com",
		Phone:   "8341099962",
		Subject: "Bulk order",
		Message: "I would like a quote for 200 litres.",
	}
}

func TestSubmitDeliversValidEnquiry(t *testing.T) {
	m := &stubMailer{}
	if err := newContact(m).Submit(context.Background(), validInput()); err != nil {
		t.Fatalf("Submit returned %v, want nil", err)
	}
	if m.count() != 1 {
		t.Fatalf("sent %d emails, want 1", m.count())
	}
	// Replying to the enquiry must reach the visitor, not our own inbox.
	if got := m.sent[0].ReplyTo; got != "ravi@example.com" {
		t.Errorf("ReplyTo = %q, want ravi@example.com", got)
	}
	if !strings.Contains(m.sent[0].Text, "200 litres") {
		t.Errorf("body missing the message text: %q", m.sent[0].Text)
	}
}

func TestSubmitRejectsInvalidInput(t *testing.T) {
	m := &stubMailer{}
	in := validInput()
	in.Email = "not-an-email"

	err := newContact(m).Submit(context.Background(), in)

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("Submit returned %v, want *ValidationError", err)
	}
	if _, ok := ve.Fields["email"]; !ok {
		t.Errorf("no error for field %q, got %v", "email", ve.Fields)
	}
	if m.count() != 0 {
		t.Errorf("sent %d emails on invalid input, want 0", m.count())
	}
}

func TestSubmitDropsHoneypotAsSuccess(t *testing.T) {
	m := &stubMailer{}
	in := validInput()
	in.Website = "http://spam.example"

	// A honeypot hit must look like success to the caller but send nothing.
	if err := newContact(m).Submit(context.Background(), in); err != nil {
		t.Fatalf("Submit returned %v, want nil for honeypot", err)
	}
	if m.count() != 0 {
		t.Errorf("sent %d emails for honeypot submission, want 0", m.count())
	}
}

func TestSubmitWrapsDeliveryFailure(t *testing.T) {
	m := &stubMailer{err: errors.New("resend is down")}

	err := newContact(m).Submit(context.Background(), validInput())

	var de *DeliveryError
	if !errors.As(err, &de) {
		t.Fatalf("Submit returned %v, want *DeliveryError", err)
	}
	// The cause must be preserved for logging via errors.Is/Unwrap.
	if !strings.Contains(de.Error(), "resend is down") {
		t.Errorf("DeliveryError lost its cause: %v", de)
	}
}
