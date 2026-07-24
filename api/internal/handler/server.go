package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/cncmarella/cleanleaf/api/internal/config"
	"github.com/cncmarella/cleanleaf/api/internal/mailer"
	"github.com/cncmarella/cleanleaf/api/internal/ratelimit"
	"github.com/cncmarella/cleanleaf/api/internal/service"
)

type Server struct {
	log     *slog.Logger
	limiter *ratelimit.Limiter
	version string
}

func New(cfg config.Config, log *slog.Logger, m mailer.Mailer, version string) (*Server, http.Handler) {
	gin.SetMode(gin.ReleaseMode)
	// Reject bodies with unknown fields, matching the strict JSON decoding the
	// API had before Gin. This flag is package-level, so it is set once here.
	binding.EnableDecoderDisallowUnknownFields = true

	s := &Server{
		log:     log,
		limiter: ratelimit.New(5, time.Hour),
		version: version,
	}

	contact := &contactHandler{svc: service.NewContact(m, log, cfg.MailFrom, cfg.MailTo)}

	r := gin.New()
	r.Use(recoverPanic(log), requestLogger(log), cors(cfg.AllowedOrigins))

	r.GET("/healthz", s.health)
	r.POST("/api/contact", rateLimit(s.limiter), contact.submit)

	// Every response from this API is JSON, including the not-found fallback.
	r.NoRoute(func(c *gin.Context) {
		writeError(c, http.StatusNotFound, "No such endpoint.")
	})

	return s, r
}

// Limiter exposes the rate limiter so cmd/server can reap it periodically.
func (s *Server) Limiter() *ratelimit.Limiter { return s.limiter }
