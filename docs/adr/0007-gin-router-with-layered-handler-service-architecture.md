# 7. Gin router with a layered handler/service architecture

Date: 2026-07-24

## Status

Accepted

Supersedes [ADR-0003](0003-go-standard-library-over-a-web-framework.md).

## Context

[ADR-0003](0003-go-standard-library-over-a-web-framework.md) built the API on
`net/http` alone, on the grounds that two endpoints did not justify a framework
or the dependencies it drags in. That reasoning still holds for the routing
itself.

Two things changed the decision:

- We want a conventional **layered structure** — an HTTP/transport layer over a
  service layer that owns the business rules (with room for a repository layer
  beneath it if persistence ever arrives, per [ADR-0005](0005-no-database-email-only-enquiries.md)).
  In the stdlib version, validation, the honeypot rule, mail orchestration and
  HTTP wiring all lived together in one handler file.
- We want Gin's request binding, middleware model and route grouping as the
  house style for this and future endpoints, rather than growing hand-written
  `func(http.Handler) http.Handler` middleware.

This is an explicit reversal of a project rule ("`api/` has zero third-party
dependencies"), so it is recorded rather than made as a convenience.

## Decision

We build the API on [Gin](https://github.com/gin-gonic/gin) and split it into
layers:

- **Handler layer** (`internal/handler/`) — Gin routes and middleware. Binds
  requests, calls the service, maps results and typed errors to JSON responses.
  Owns transport-only concerns: routing, CORS, rate limiting, request logging,
  panic recovery.
- **Service layer** (`internal/service/`) — the business rules, independent of
  HTTP: normalisation, validation, the honeypot decision and mail delivery
  orchestration. Testable without a server.
- **Repository layer** — not built. There is nothing to persist yet
  ([ADR-0005](0005-no-database-email-only-enquiries.md)); when there is, it sits
  beneath the service.

Behaviour of the two endpoints is unchanged — the existing handler test suite
passes against the new stack.

## Consequences

- Gin is the codebase's first third-party dependency, and it is not a small one:
  it pulls in roughly thirty transitive modules (validator, sonic, etc.). We now
  have a dependency-update and supply-chain surface where ADR-0003 had none.
  Adding a *further* dependency is still an ADR-level decision.
- The Docker image is no longer a near-empty static binary; it grows with Gin
  and its dependencies, though it stays a single `distroless/static` binary.
- Business logic can be unit-tested directly against the service, without
  spinning up HTTP. New endpoints follow the same handler → service split.
- Some stdlib behaviour is re-established explicitly on top of Gin: strict JSON
  decoding via `binding.EnableDecoderDisallowUnknownFields`, our own slog access
  logger and JSON panic recovery (via `gin.New`, not `gin.Default`), and a
  JSON `NoRoute` fallback so every response stays JSON.
- The "small, contained change to `chi`" escape hatch that ADR-0003 kept open is
  now spent; moving off Gin later would be a larger change.
