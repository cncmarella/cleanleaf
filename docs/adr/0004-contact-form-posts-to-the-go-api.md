# 4. The contact form posts to the Go API, not a Server Action

Date: 2026-07-22

## Status

Accepted

## Context

Next.js offers Server Actions, which let a `<form>` call a server function
directly with no hand-written fetch, no API route and no CORS configuration.
For a contact form this is the path of least resistance, and the Next.js docs
present it as the default.

Using it here would mean the validation rules, spam filtering and email
delivery live in TypeScript on Vercel — which contradicts ADR-0002, where we
decided the Go service owns business logic.

## Decision

The contact form is a client component that `POST`s JSON to
`POST /api/contact` on the Go service. The Go handler owns validation, rate
limiting, honeypot filtering and email delivery. The frontend only collects
input and renders whatever the API returns.

The API returns a single error shape for every failure, so the frontend has
one thing to parse:

```json
{ "error": "human readable message", "fields": { "email": "field message" } }
```

## Consequences

- All enquiry logic is in Go, testable with `go test` and with no browser or
  Next.js runtime involved. The handler has test coverage for validation,
  honeypot, rate limiting, oversized bodies and mailer failure.
- The frontend stays a static site: all four pages prerender, and the only
  client-side JavaScript is the header menu and the form.
- We must configure CORS. `ALLOWED_ORIGINS` on the API must include the
  frontend's origin, and **this is the thing that will break** when the site
  moves to a custom domain. It is called out in the README deploy checklist.
- `NEXT_PUBLIC_API_URL` is baked into the client bundle at build time, so
  changing the API URL requires a frontend rebuild, not just an env var edit.
- Browser validation is duplicated (HTML `required` attributes for immediate
  feedback, Go for the real check). The Go side is authoritative; the HTML
  attributes are a convenience and are never trusted.
