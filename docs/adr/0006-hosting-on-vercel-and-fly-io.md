# 6. Host the frontend on Vercel and the API on Fly.io

Date: 2026-07-22

## Status

Accepted

## Context

Two apps (ADR-0002) need hosting. The constraints: near-zero cost while the
site has no traffic, minimal operational work, deploys on `git push`, and a
path to a custom domain later.

Options considered:

- **A single VPS** (Hetzner/DigitalOcean, ~$5/mo) running both behind Caddy.
  Full control and one bill, but we own OS updates, TLS renewal, backups and
  uptime — real recurring work for a site that is mostly static.
- **All on Cloudflare** (Pages + Workers). Cheapest at scale, but Go on
  Workers requires either WASM or the Containers product; both are awkward
  compared to running a normal Go binary.
- **Vercel + Fly.io.**

## Decision

- `web/` deploys to **Vercel**, with the project root set to `web/`. Vercel is
  the reference host for Next.js; the four static pages are served from its
  CDN.
- `api/` deploys to **Fly.io** in the `bom` (Mumbai) region — the closest
  region to Hyderabad — from the Dockerfile in `api/`.

Both deploy from the `main` branch of the same repository.

## Consequences

- Both fit comfortably in free/hobby tiers at current traffic.
- Deploys are `git push` (Vercel) and `fly deploy` (API). No servers to patch.
- The API is configured to scale to zero (`min_machines_running = 0`), so the
  first enquiry after an idle period pays roughly a second of cold start. For a
  contact form this is an acceptable trade for the cost saving; raise it to 1
  if the delay ever becomes visible to customers.
- Two vendors means two dashboards and two places to set environment
  variables. The README documents which variable lives where.
- Vendor lock-in is low in both directions: the frontend is a standard Next.js
  app and the API is a Dockerfile, both movable to a VPS if the free tiers
  change.
- Deploying from a single-region host means visitors far from Mumbai see
  higher API latency. Irrelevant for a form used by customers in Telangana.
