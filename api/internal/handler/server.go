// Package handler wires the HTTP routes for the CleanLeaf website API.
package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/cncmarella/cleanleaf/api/internal/config"
	"github.com/cncmarella/cleanleaf/api/internal/mailer"
	"github.com/cncmarella/cleanleaf/api/internal/ratelimit"
)

// Server holds the dependencies shared by every handler.
type Server struct {
	log      *slog.Logger
	mailer   mailer.Mailer
	limiter  *ratelimit.Limiter
	mailFrom string
	mailTo   string
	version  string
}

// New builds a Server and returns it alongside its HTTP handler.
func New(cfg config.Config, log *slog.Logger, m mailer.Mailer, version string) (*Server, http.Handler) {
	s := &Server{
		log:      log,
		mailer:   m,
		limiter:  ratelimit.New(5, time.Hour),
		mailFrom: cfg.MailFrom,
		mailTo:   cfg.MailTo,
		version:  version,
	}
	return s, s.routes(cfg)
}

// Limiter exposes the rate limiter so cmd/server can reap it periodically.
func (s *Server) Limiter() *ratelimit.Limiter { return s.limiter }

func (s *Server) routes(cfg config.Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("POST /api/contact", s.handleContact)

	// ServeMux answers unmatched paths with a plain-text 404; override it so
	// every response from this API is JSON.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, s.log, http.StatusNotFound, "No such endpoint.")
	})

	return chain(mux,
		recoverPanic(s.log),
		requestLogger(s.log),
		cors(cfg.AllowedOrigins),
	)
}
