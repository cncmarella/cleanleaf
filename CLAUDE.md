# CleanLeaf — working notes

Marketing website for CleanLeaf Technologies (crop protection manufacturer,
Ghatkesar, Hyderabad). See [README.md](README.md) for setup and deployment.

**This is Nithin's personal project — unrelated to his office work.**

## Shape

Two independently deployed apps in one repo ([ADR-0002](docs/adr/0002-monorepo-with-separate-frontend-and-backend.md)):
`web/` (Next.js → Vercel) and `api/` (Go → Fly.io), talking over HTTP/JSON.

## Rules that are easy to get wrong

- **Business logic goes in Go, not TypeScript.** The frontend collects input and
  renders responses; validation, spam filtering and delivery live in `api/`.
  Do not move contact handling into a Server Action ([ADR-0004](docs/adr/0004-contact-form-posts-to-the-go-api.md)).
- **`api/` is layered: handler → service** ([ADR-0007](docs/adr/0007-gin-router-with-layered-handler-service-architecture.md)).
  Business rules (validation, honeypot, mail orchestration) live in
  `internal/service/`; `internal/handler/` only speaks HTTP (Gin binding,
  middleware, mapping results to JSON). A repository layer would sit beneath the
  service, but there is nothing to persist yet ([ADR-0005](docs/adr/0005-no-database-email-only-enquiries.md)).
- **Gin is the router, and the only heavy dependency** ([ADR-0007](docs/adr/0007-gin-router-with-layered-handler-service-architecture.md)
  superseded the stdlib-only [ADR-0003](docs/adr/0003-go-standard-library-over-a-web-framework.md)).
  Adding a *further* third-party dependency is still an ADR-level decision, not a
  convenience.
- **`web/` runs Next.js 16**, which has breaking changes from older versions.
  `web/AGENTS.md` says to read `web/node_modules/next/dist/docs/` before writing
  Next-specific code. Do that — e.g. `React.FormEvent` is deprecated in React 19
  in favour of `React.SubmitEvent`.
- **Content lives in `web/src/lib/site.ts`**, not scattered through pages. The
  product list there is placeholder text pending the real catalogue.
- **CORS is the thing that breaks on deploy.** `ALLOWED_ORIGINS` in
  `api/fly.toml` must list the exact frontend origin.

## Checks before calling something done

```bash
cd api && go test ./... -race && go vet ./... && gofmt -l .
cd web && npx tsc --noEmit && npm run lint && npm run build
```

## Decisions

Record architecturally significant choices in [`docs/adr/`](docs/adr/README.md)
— anything whose reversal would cross the frontend/backend, hosting or data
boundary. ADRs are immutable; supersede rather than edit.
