# CleanLeaf Technologies — Website

Marketing website for CleanLeaf Technologies, a crop protection manufacturer in
Ghatkesar, Hyderabad. A Next.js frontend and a Go API, deployed separately.

## Stack

| Part     | Choice                                        | Why                                     |
| -------- | --------------------------------------------- | --------------------------------------- |
| Frontend | Next.js 16 (App Router), React 19, TypeScript, Tailwind v4 | Static pages, good SEO out of the box   |
| Backend  | Go 1.24, standard library only                | Zero dependencies ([ADR-0003](docs/adr/0003-go-standard-library-over-a-web-framework.md)) |
| Email    | Resend API                                    | No database needed yet ([ADR-0005](docs/adr/0005-no-database-email-only-enquiries.md)) |
| Hosting  | Vercel (web) + Fly.io (api)                   | Free tier, deploy on push ([ADR-0006](docs/adr/0006-hosting-on-vercel-and-fly-io.md)) |

## Layout

```
web/    Next.js frontend
  src/app/          Pages: /, /products, /about, /contact
  src/components/   Header, footer, contact form, logo
  src/lib/site.ts   Company details and product data — edit content here
api/    Go backend
  cmd/server/       Entrypoint, graceful shutdown
  internal/config/  Environment configuration
  internal/handler/ Routes, middleware, contact + health handlers
  internal/mailer/  Resend and console mail transports
  internal/ratelimit/ In-memory per-IP limiter
docs/adr/  Architecture Decision Records — read these before changing direction
```

## Running locally

You need **two terminals**: the frontend calls the API over HTTP
([ADR-0004](docs/adr/0004-contact-form-posts-to-the-go-api.md)).

**Terminal 1 — API** (http://localhost:8080)

```bash
cd api
cp .env.example .env      # first time only
go run ./cmd/server
```

Without `RESEND_API_KEY` the server uses the **console mailer**: enquiries are
printed to the log instead of emailed. That is the intended dev setup — no
secrets required.

**Terminal 2 — frontend** (http://localhost:3000)

```bash
cd web
cp .env.example .env.local   # first time only
npm install                  # first time only
npm run dev
```

Submit the contact form, and the enquiry appears in the API terminal as a
`email not sent (console mailer)` log line.

## Checks

```bash
cd api && go test ./... -race && go vet ./... && gofmt -l .
cd web && npx tsc --noEmit && npm run lint && npm run build
```

The Go tests cover contact validation, the honeypot, rate limiting, oversized
bodies, mailer failure and CORS.

## API

Base URL: `http://localhost:8080` in dev.

### `GET /healthz`

```json
{ "status": "ok", "version": "a1b2c3d" }
```

Dependency-free by design, so a mail outage never fails the Fly health check.

### `POST /api/contact`

```json
{
  "name": "Ravi Kumar",
  "email": "ravi@example.com",
  "phone": "8341099962",
  "subject": "Bulk order",
  "message": "Need a quote for 200 litres for cotton.",
  "website": ""
}
```

`phone`, `subject` and `website` are optional. `website` is a **honeypot** —
hidden from users; if filled, the request is silently dropped with a success
response so bots get no signal.

| Status | Meaning                                                    |
| ------ | ---------------------------------------------------------- |
| 202    | Accepted and emailed                                       |
| 400    | Malformed JSON                                             |
| 413    | Body over 16 KiB                                           |
| 422    | Validation failed — see `fields`                           |
| 429    | Rate limited (5 per hour per IP)                           |
| 502    | Mail provider failed; the enquiry was **not** delivered    |

Every error uses one shape:

```json
{ "error": "Please correct the highlighted fields.", "fields": { "email": "..." } }
```

## Environment variables

**API** (`api/.env` locally; `fly.toml` + `fly secrets` in production)

| Variable          | Default                  | Notes                                                     |
| ----------------- | ------------------------ | --------------------------------------------------------- |
| `APP_ENV`         | `development`            | `production` makes `RESEND_API_KEY` mandatory at boot      |
| `PORT`            | `8080`                   |                                                            |
| `ALLOWED_ORIGINS` | `http://localhost:3000`  | Comma-separated. **Must include the live frontend origin** |
| `RESEND_API_KEY`  | _empty_                  | Secret. Empty ⇒ console mailer                             |
| `MAIL_FROM`       | `onboarding@resend.dev`  | Must be a Resend-verified domain to use your own address   |
| `MAIL_TO`         | `cleanleaf789@gmail.com` | Where enquiries land                                       |

**Frontend** (`web/.env.local` locally; Vercel project settings in production)

| Variable              | Notes                                                                  |
| --------------------- | ---------------------------------------------------------------------- |
| `NEXT_PUBLIC_API_URL` | Base URL of the API. Inlined into the client bundle — **public**, and changing it needs a rebuild |

## Deploying

Do these in order. Steps 1 and 2 must both be done before the contact form
works in production.

### 1. Email — Resend

1. Sign up at [resend.com](https://resend.com) and create an API key.
2. Keep the default `onboarding@resend.dev` sender for now. It only delivers to
   the address that owns the Resend account, which is fine for testing.
3. To send from a CleanLeaf address later, verify a domain in Resend and update
   `MAIL_FROM`.

### 2. API — Fly.io

```bash
brew install flyctl
fly auth signup                    # or: fly auth login

cd api
fly launch --no-deploy --copy-config   # keeps the committed fly.toml
fly secrets set RESEND_API_KEY=re_your_key_here
fly deploy --remote-only               # --remote-only builds without local Docker
```

Note the hostname it prints, e.g. `https://cleanleaf-api.fly.dev`. Verify:

```bash
curl https://cleanleaf-api.fly.dev/healthz
```

### 3. Frontend — Vercel

1. Push this repo to GitHub.
2. On [vercel.com](https://vercel.com), **Add New → Project**, import the repo.
3. Set **Root Directory** to `web`. This matters — the repo root is not the app.
4. Add an environment variable: `NEXT_PUBLIC_API_URL` = the Fly URL from step 2.
5. Deploy. Note the URL, e.g. `https://cleanleaf.vercel.app`.

### 4. Close the CORS loop

The API currently only trusts the origin in `api/fly.toml`. Update it to the
real Vercel URL, then redeploy:

```bash
cd api
# edit ALLOWED_ORIGINS in fly.toml to your actual Vercel URL
fly deploy --remote-only
```

Also update `site.url` in `web/src/lib/site.ts` so metadata points at the live
site, and redeploy the frontend.

**Then submit the contact form on the live site and confirm the email arrives.**
A CORS mismatch shows up as "We could not reach our servers" in the form.

### 5. Custom domain (later)

When you buy one (Cloudflare Registrar is cheapest for `.com`; BigRock or
GoDaddy for `.in`):

1. Vercel → Project → Settings → Domains → add the domain, follow its DNS
   instructions.
2. Add the new origin to `ALLOWED_ORIGINS` in `api/fly.toml`, `fly deploy`.
3. Update `site.url` in `web/src/lib/site.ts`, redeploy.
4. Optionally give the API a subdomain (`api.yourdomain.com`) via
   `fly certs add`.

## Editing content

Most copy changes are one file: [`web/src/lib/site.ts`](web/src/lib/site.ts) —
company details, navigation, the product list and the "why us" points. The
product catalogue there is **placeholder text** and should be replaced with the
real range, actives, pack sizes and registration numbers.

## Decisions

Architecture decisions are recorded in [`docs/adr/`](docs/adr/README.md). Read
them before changing direction, and add a new one when you make a call that
would be expensive to reverse.

## Roadmap

- [ ] Replace placeholder product data with the real catalogue
- [ ] Real logo asset and product photography
- [ ] Per-product detail pages
- [ ] Store enquiries (supersedes [ADR-0005](docs/adr/0005-no-database-email-only-enquiries.md))
- [ ] Custom domain
- [ ] Telugu translation
