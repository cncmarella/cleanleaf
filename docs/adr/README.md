# Architecture Decision Records

Why the system looks the way it does. See [ADR-0001](0001-record-architecture-decisions.md)
for the process itself.

| #                                                          | Decision                                          | Status   |
| ---------------------------------------------------------- | ------------------------------------------------- | -------- |
| [0001](0001-record-architecture-decisions.md)               | Record architecture decisions                     | Accepted |
| [0002](0002-monorepo-with-separate-frontend-and-backend.md) | Monorepo with a separate frontend and backend     | Accepted |
| [0003](0003-go-standard-library-over-a-web-framework.md)    | Go standard library over a web framework          | Superseded by [0007](0007-gin-router-with-layered-handler-service-architecture.md) |
| [0004](0004-contact-form-posts-to-the-go-api.md)            | Contact form posts to the Go API, not a Server Action | Accepted |
| [0005](0005-no-database-email-only-enquiries.md)            | No database yet — enquiries are emailed           | Accepted |
| [0006](0006-hosting-on-vercel-and-fly-io.md)                | Host the frontend on Vercel and the API on Fly.io | Accepted |
| [0007](0007-gin-router-with-layered-handler-service-architecture.md) | Gin router with a layered handler/service architecture | Accepted |

## Writing a new one

Copy the structure of an existing ADR: **Status**, **Context**, **Decision**,
**Consequences**. Number it sequentially. Add a row above.

Record a decision when reversing it would mean changes across a boundary —
framework, hosting, data store, or the frontend/backend contract. Do not record
routine implementation choices.

ADRs are immutable once accepted. To change course, write a new ADR and mark
the old one `Superseded by ADR-NNNN`.
