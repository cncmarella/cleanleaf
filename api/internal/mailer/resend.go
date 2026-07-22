package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const resendEndpoint = "https://api.resend.com/emails"

// Resend sends mail through the Resend HTTP API.
type Resend struct {
	APIKey string
	Client *http.Client
}

// NewResend returns a Resend mailer with a bounded HTTP client.
func NewResend(apiKey string) *Resend {
	return &Resend{
		APIKey: apiKey,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	ReplyTo string   `json:"reply_to,omitempty"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

// Send delivers the message, returning an error for any non-2xx response.
func (r *Resend) Send(ctx context.Context, m Message) error {
	body, err := json.Marshal(resendRequest{
		From:    m.From,
		To:      []string{m.To},
		ReplyTo: m.ReplyTo,
		Subject: m.Subject,
		Text:    m.Text,
	})
	if err != nil {
		return fmt.Errorf("encode resend request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build resend request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.Client.Do(req)
	if err != nil {
		return fmt.Errorf("call resend: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Cap the body read so a misbehaving upstream cannot blow up the log.
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<10))
		return fmt.Errorf("resend returned %s: %s", resp.Status, bytes.TrimSpace(detail))
	}
	return nil
}
