package handler

import (
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cncmarella/cleanleaf/api/internal/ratelimit"
)

func requestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"ip", clientIP(c.Request),
		)
	}
}

func recoverPanic(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if v := recover(); v != nil {
				log.Error("panic recovered", "value", v, "path", c.Request.URL.Path)
				c.AbortWithStatusJSON(http.StatusInternalServerError, errorBody{Error: "Something went wrong on our side."})
			}
		}()
		c.Next()
	}
}

func cors(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		origin := c.Request.Header.Get("Origin")
		if origin != "" && slices.Contains(allowed, origin) {
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Content-Type")
			h.Set("Access-Control-Max-Age", "86400")
		}
		// Caches must not serve one origin's response to another.
		h.Add("Vary", "Origin")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// rateLimit rejects a client that has exceeded the per-IP enquiry budget before the request reaches the handler.
func rateLimit(limiter *ratelimit.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(clientIP(c.Request)) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, errorBody{Error: "Too many enquiries from this address. Please try again in a few minutes."})
			return
		}
		c.Next()
	}
}

// clientIP prefers Fly.io's Fly-Client-IP, then the left-most X-Forwarded-For
// entry, then the socket address. Only trust these behind a proxy that sets
// them -- directly exposed, they are attacker-controlled.
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("Fly-Client-IP"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		first, _, _ := strings.Cut(fwd, ",")
		return strings.TrimSpace(first)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
