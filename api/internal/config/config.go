// Package config loads runtime configuration from the environment.
package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds every knob the server reads at startup. Everything comes from
// the environment so the same binary runs unchanged in dev and on Fly.io.
type Config struct {
	Env             string
	Port            string
	AllowedOrigins  []string
	ShutdownTimeout time.Duration

	// Mail delivery. When ResendAPIKey is empty the server falls back to the
	// console mailer, which means `go run ./cmd/server` works with no secrets.
	ResendAPIKey string
	MailFrom     string
	MailTo       string
}

// Load reads configuration from the environment and validates it.
func Load() (Config, error) {
	c := Config{
		Env:             env("APP_ENV", "development"),
		Port:            env("PORT", "8080"),
		ShutdownTimeout: 10 * time.Second,
		ResendAPIKey:    env("RESEND_API_KEY", ""),
		MailFrom:        env("MAIL_FROM", "CleanLeaf Website <onboarding@resend.dev>"),
		MailTo:          env("MAIL_TO", "cleanleaf789@gmail.com"),
	}

	origins := env("ALLOWED_ORIGINS", "http://localhost:3000")
	for _, o := range strings.Split(origins, ",") {
		if o = strings.TrimSpace(o); o != "" {
			c.AllowedOrigins = append(c.AllowedOrigins, o)
		}
	}
	if len(c.AllowedOrigins) == 0 {
		return Config{}, fmt.Errorf("ALLOWED_ORIGINS must list at least one origin")
	}

	// A production deploy that silently logs enquiries instead of mailing them
	// is worse than a failed boot -- the enquiries are gone and nobody notices.
	if c.Env == "production" && c.ResendAPIKey == "" {
		return Config{}, fmt.Errorf("RESEND_API_KEY is required when APP_ENV=production")
	}

	return c, nil
}

// IsProduction reports whether the server is running in the production deploy.
func (c Config) IsProduction() bool { return c.Env == "production" }

func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
