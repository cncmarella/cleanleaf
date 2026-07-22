# 5. No database yet — enquiries are emailed, not stored

Date: 2026-07-22

## Status

Accepted

## Context

The first version of the site has one piece of dynamic behaviour: a contact
form. The product catalogue is four fixed categories, currently hardcoded in
`web/src/lib/site.ts`.

Adding Postgres now would mean provisioning a database, adding a driver and
migrations, managing a connection string, and handling the operational
questions that follow (backups, connection limits on a scale-to-zero host).
None of that serves a page that shows four static cards.

## Decision

The first version ships with no database. Contact enquiries are delivered by
email to the company inbox via the Resend API and are not persisted anywhere.
Product content lives in a typed TypeScript module in the frontend.

The `mailer.Mailer` interface keeps delivery behind an abstraction so the
transport can change without touching the handler.

## Consequences

- There is no infrastructure to run, pay for or back up. The API is a stateless
  binary that can scale to zero.
- **Enquiries exist only in the company's email inbox.** If a send fails, the
  enquiry is lost — the API returns a 502 telling the visitor to phone instead,
  and logs the failure. This is the main risk we are accepting.
- The in-memory rate limiter is per-instance and resets when the machine
  scales to zero. Adequate for casual spam, not for a determined attacker.
- There is no admin view of past enquiries and no analytics on them.
- When we do need a database — storing enquiries, a real product catalogue with
  registration numbers, or dealer logins — that gets a new ADR superseding
  this one. Expect the catalogue to be the trigger.
