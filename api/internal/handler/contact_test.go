package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/cncmarella/cleanleaf/api/internal/config"
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

func newTestServer(t *testing.T, m mailer.Mailer) http.Handler {
	t.Helper()
	cfg := config.Config{
		Env:            "test",
		AllowedOrigins: []string{"http://localhost:3000"},
		MailFrom:       "site@cleanleaf.test",
		MailTo:         "sales@cleanleaf.test",
	}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	_, routes := New(cfg, log, m, "test")
	return routes
}

func postContact(t *testing.T, h http.Handler, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/contact", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

const validEnquiry = `{
	"name": "Ravi Kumar",
	"email": "ravi@example.com",
	"phone": "8341099962",
	"subject": "Bulk order",
	"message": "I would like a quote for 200 litres."
}`

func TestContactAcceptsValidEnquiry(t *testing.T) {
	m := &stubMailer{}
	rec := postContact(t, newTestServer(t, m), validEnquiry)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusAccepted, rec.Body)
	}
	if m.count() != 1 {
		t.Fatalf("sent %d emails, want 1", m.count())
	}

	got := m.sent[0]
	if got.To != "sales@cleanleaf.test" {
		t.Errorf("To = %q, want sales@cleanleaf.test", got.To)
	}
	// Replying to the enquiry must reach the visitor, not our own inbox.
	if got.ReplyTo != "ravi@example.com" {
		t.Errorf("ReplyTo = %q, want ravi@example.com", got.ReplyTo)
	}
	if !strings.Contains(got.Text, "200 litres") {
		t.Errorf("body missing the message text: %q", got.Text)
	}
}

func TestContactValidation(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantField string
	}{
		{"missing name", `{"name":"","email":"a@b.com","message":"a long enough message"}`, "name"},
		{"missing email", `{"name":"Ravi","email":"","message":"a long enough message"}`, "email"},
		{"malformed email", `{"name":"Ravi","email":"not-an-email","message":"a long enough message"}`, "email"},
		{"short message", `{"name":"Ravi","email":"a@b.com","message":"hi"}`, "message"},
		{"implausible phone", `{"name":"Ravi","email":"a@b.com","phone":"12","message":"a long enough message"}`, "phone"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &stubMailer{}
			rec := postContact(t, newTestServer(t, m), tc.body)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
			}
			var body errorBody
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if _, ok := body.Fields[tc.wantField]; !ok {
				t.Errorf("no error for field %q, got %v", tc.wantField, body.Fields)
			}
			if m.count() != 0 {
				t.Errorf("sent %d emails on invalid input, want 0", m.count())
			}
		})
	}
}

func TestContactHoneypotLooksLikeSuccess(t *testing.T) {
	m := &stubMailer{}
	body := `{"name":"Bot","email":"bot@spam.com","message":"buy cheap things now","website":"http://spam.example"}`
	rec := postContact(t, newTestServer(t, m), body)

	// A bot must not be able to tell it was filtered, so the status matches
	// the happy path -- but nothing is actually mailed.
	if rec.Code != http.StatusAccepted {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if m.count() != 0 {
		t.Errorf("sent %d emails for honeypot submission, want 0", m.count())
	}
}

func TestContactReportsMailerFailure(t *testing.T) {
	m := &stubMailer{err: errors.New("resend is down")}
	rec := postContact(t, newTestServer(t, m), validEnquiry)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
	// The upstream error must not leak to the visitor.
	if strings.Contains(rec.Body.String(), "resend is down") {
		t.Errorf("response leaked internal error: %s", rec.Body)
	}
}

func TestContactRateLimitsRepeatSubmissions(t *testing.T) {
	m := &stubMailer{}
	h := newTestServer(t, m)

	// The limiter allows 5 per hour per IP; httptest uses a fixed RemoteAddr.
	for i := range 5 {
		if rec := postContact(t, h, validEnquiry); rec.Code != http.StatusAccepted {
			t.Fatalf("submission %d: status = %d, want %d", i+1, rec.Code, http.StatusAccepted)
		}
	}
	if rec := postContact(t, h, validEnquiry); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("6th submission: status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
	if m.count() != 5 {
		t.Errorf("sent %d emails, want 5", m.count())
	}
}

func TestContactRejectsOversizedBody(t *testing.T) {
	m := &stubMailer{}
	huge := strings.Repeat("a", 32<<10)
	body, err := json.Marshal(map[string]string{"name": "Ravi", "email": "a@b.com", "message": huge})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	rec := postContact(t, newTestServer(t, m), string(body))

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
	if m.count() != 0 {
		t.Errorf("sent %d emails, want 0", m.count())
	}
}

func TestCORSOnlyEchoesAllowedOrigins(t *testing.T) {
	h := newTestServer(t, &stubMailer{})

	tests := []struct {
		origin string
		want   string
	}{
		{"http://localhost:3000", "http://localhost:3000"},
		{"https://evil.example", ""},
	}

	for _, tc := range tests {
		t.Run(tc.origin, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodOptions, "/api/contact", nil)
			req.Header.Set("Origin", tc.origin)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != tc.want {
				t.Errorf("Allow-Origin = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHealthzReportsOK(t *testing.T) {
	h := newTestServer(t, &stubMailer{})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body healthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Status != "ok" {
		t.Errorf("status = %q, want ok", body.Status)
	}
}

func TestUnknownRouteReturnsJSON(t *testing.T) {
	h := newTestServer(t, &stubMailer{})
	req := httptest.NewRequest(http.MethodGet, "/nope", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want JSON", ct)
	}
	if !json.Valid(bytes.TrimSpace(rec.Body.Bytes())) {
		t.Errorf("body is not valid JSON: %s", rec.Body)
	}
}
