# 2. Monorepo with a separate frontend and backend

Date: 2026-07-22

## Status

Accepted

## Context

The site needs a public marketing frontend and a backend to handle form
submissions and, later, product data. Two obvious shapes were available:

1. One Next.js application handling both UI and server logic (Route Handlers
   or Server Actions), deployed as a single unit.
2. A separate frontend and a separate backend service.

The owner's stated preference is to write the backend in Go and to review
rather than author the frontend. That makes option 1 a poor fit: it would put
all server logic in TypeScript.

## Decision

We keep one Git repository containing two independently deployable apps:

```
web/    Next.js frontend, deployed to Vercel
api/    Go backend, deployed to Fly.io
docs/   ADRs and project documentation
```

They communicate over HTTP/JSON only. Neither imports code from the other, and
there is no shared build step between them.

## Consequences

- The Go service owns all business logic, which matches where the owner's
  expertise is.
- One repository means one place for issues, one branch per change, and no
  cross-repo version skew in documentation.
- The two apps deploy separately, so a frontend copy change does not redeploy
  the API and vice versa.
- The cost is a real network boundary: CORS must be configured, and local
  development needs both processes running. We accept this as the price of
  keeping the backend in Go.
- The JSON contract between the two is now a real interface. Changing a
  response shape means changing both sides deliberately.
