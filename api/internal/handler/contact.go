package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/cncmarella/cleanleaf/api/internal/mailer"
)

// maxContactBody caps the request body well above a legitimate enquiry.
const maxContactBody = 16 << 10 // 16 KiB

type contactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
	Message string `json:"message"`

	// Website is a honeypot: hidden in the form, so a human leaves it empty and
	// a naive bot fills it in. Populated submissions are dropped silently.
	Website string `json:"website"`
}

type contactResponse struct {
	Message string `json:"message"`
}

// handleContact validates an enquiry and emails it to the sales inbox.
func (s *Server) handleContact(w http.ResponseWriter, r *http.Request) {
	if !s.limiter.Allow(clientIP(r)) {
		writeError(w, s.log, http.StatusTooManyRequests, "Too many enquiries from this address. Please try again in a few minutes.")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxContactBody)
	var req contactRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, s.log, http.StatusRequestEntityTooLarge, "That message is too long.")
			return
		}
		writeError(w, s.log, http.StatusBadRequest, "We could not read that request.")
		return
	}

	// Answer the honeypot with the same success shape a human gets, so a bot
	// has no signal that it was rejected.
	if strings.TrimSpace(req.Website) != "" {
		s.log.Info("contact honeypot triggered", "ip", clientIP(r))
		writeJSON(w, s.log, http.StatusAccepted, contactResponse{Message: "Thanks for reaching out. We'll be in touch shortly."})
		return
	}

	req.normalise()
	if fields := req.validate(); len(fields) > 0 {
		writeFieldErrors(w, s.log, fields)
		return
	}

	msg := mailer.Message{
		From:    s.mailFrom,
		To:      s.mailTo,
		ReplyTo: req.Email,
		Subject: fmt.Sprintf("Website enquiry: %s", req.Subject),
		Text:    req.body(),
	}
	if err := s.mailer.Send(r.Context(), msg); err != nil {
		// The visitor gets a generic message; the detail stays in our logs.
		s.log.Error("send contact email", "error", err)
		writeError(w, s.log, http.StatusBadGateway, "We could not send your message right now. Please call us on 8341099962.")
		return
	}

	s.log.Info("contact enquiry accepted", "subject", req.Subject)
	writeJSON(w, s.log, http.StatusAccepted, contactResponse{Message: "Thanks for reaching out. We'll be in touch shortly."})
}

func (c *contactRequest) normalise() {
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	c.Phone = strings.TrimSpace(c.Phone)
	c.Subject = strings.TrimSpace(c.Subject)
	c.Message = strings.TrimSpace(c.Message)
	if c.Subject == "" {
		c.Subject = "General enquiry"
	}
}

// validate returns a field-name -> message map, empty when the request is good.
func (c *contactRequest) validate() map[string]string {
	fields := make(map[string]string)

	switch n := utf8.RuneCountInString(c.Name); {
	case n == 0:
		fields["name"] = "Please tell us your name."
	case n > 100:
		fields["name"] = "Please keep your name under 100 characters."
	}

	if c.Email == "" {
		fields["email"] = "Please give us an email address to reply to."
	} else if _, err := mail.ParseAddress(c.Email); err != nil {
		fields["email"] = "That does not look like a valid email address."
	}

	// Phone is optional, but if given it should be plausible.
	if c.Phone != "" {
		digits := 0
		for _, r := range c.Phone {
			if r >= '0' && r <= '9' {
				digits++
			}
		}
		if digits < 7 || digits > 15 {
			fields["phone"] = "Please enter a valid phone number, or leave it blank."
		}
	}

	if n := utf8.RuneCountInString(c.Subject); n > 150 {
		fields["subject"] = "Please keep the subject under 150 characters."
	}

	switch n := utf8.RuneCountInString(c.Message); {
	case n < 10:
		fields["message"] = "Please give us a little more detail (at least 10 characters)."
	case n > 4000:
		fields["message"] = "Please keep your message under 4000 characters."
	}

	return fields
}

// body renders the plain-text email sent to the sales inbox.
func (c *contactRequest) body() string {
	phone := c.Phone
	if phone == "" {
		phone = "(not provided)"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "New enquiry from the CleanLeaf website.\n\n")
	fmt.Fprintf(&b, "Name:    %s\n", c.Name)
	fmt.Fprintf(&b, "Email:   %s\n", c.Email)
	fmt.Fprintf(&b, "Phone:   %s\n", phone)
	fmt.Fprintf(&b, "Subject: %s\n\n", c.Subject)
	fmt.Fprintf(&b, "Message:\n%s\n", c.Message)
	return b.String()
}
