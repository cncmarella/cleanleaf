# 3. Go standard library over a web framework

Date: 2026-07-22

## Status

Accepted

## Context

The API needs routing, JSON encoding, middleware (logging, recovery, CORS) and
request validation. The Go ecosystem offers Gin, Echo, Fiber and chi for this.

Since Go 1.22, the standard library's `http.ServeMux` supports method-aware
patterns (`POST /api/contact`) and path wildcards, which covers the routing
that previously justified reaching for a router.

The current API surface is two endpoints.

## Decision

We build the API on `net/http` from the standard library, with no third-party
dependencies. Middleware is a plain `func(http.Handler) http.Handler` chain
written by hand in `internal/handler/middleware.go`.

## Consequences

- `go.mod` has zero third-party requirements. There is no dependency-update
  treadmill and no supply-chain surface for a site this small.
- The Docker image is a single static binary on `distroless/static`.
- Anyone who knows Go can read the whole API without first learning a
  framework's conventions.
- We write our own middleware. That is roughly 100 lines today; if it grows
  past what is comfortable to maintain, adopting `chi` is a small, contained
  change because `chi` is `net/http`-compatible. Revisit if that happens.
- Validation is hand-written per handler. If the number of endpoints grows,
  reconsider a validation library rather than repeating the pattern.
