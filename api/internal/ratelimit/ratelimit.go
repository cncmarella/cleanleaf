// Package ratelimit provides a small in-memory rate limiter.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter allows at most Burst events per key within Window, using a fixed
// window per key. State lives in memory, so it is per-instance: good enough to
// stop casual contact-form spam, not a defence against a distributed flood.
type Limiter struct {
	window time.Duration
	burst  int

	mu      sync.Mutex
	buckets map[string]*bucket
}

type bucket struct {
	count      int
	windowEnds time.Time
}

// New returns a Limiter allowing burst events per key per window.
func New(burst int, window time.Duration) *Limiter {
	return &Limiter{
		window:  window,
		burst:   burst,
		buckets: make(map[string]*bucket),
	}
}

// Allow records an event for key and reports whether it is within the limit.
func (l *Limiter) Allow(key string) bool {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	b, ok := l.buckets[key]
	if !ok || now.After(b.windowEnds) {
		l.buckets[key] = &bucket{count: 1, windowEnds: now.Add(l.window)}
		return true
	}
	if b.count >= l.burst {
		return false
	}
	b.count++
	return true
}

// Reap deletes expired buckets so the map does not grow without bound. Callers
// should run it periodically; see the goroutine started in cmd/server.
func (l *Limiter) Reap() {
	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	for key, b := range l.buckets {
		if now.After(b.windowEnds) {
			delete(l.buckets, key)
		}
	}
}
