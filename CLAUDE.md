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
- **`api/` has zero third-party dependencies** ([ADR-0003](docs/adr/0003-go-standard-library-over-a-web-framework.md)).
  Adding one is an ADR-level decision, not a convenience.
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
