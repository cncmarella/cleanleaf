# 1. Record architecture decisions

Date: 2026-07-22

## Status

Accepted

## Context

This project will be built in small increments over many weeks, in sessions
spread far apart. Without a written record, the reasoning behind a choice is
lost by the time we revisit it, and we end up re-litigating settled questions
or silently contradicting an earlier decision.

## Decision

We record every architecturally significant decision as an Architecture
Decision Record (ADR) in `docs/adr/`, numbered sequentially, following
Michael Nygard's format: **Context**, **Decision**, **Consequences**.

A decision is architecturally significant if reversing it would require
changing multiple files across a boundary — the choice of framework, hosting,
data store, or the contract between frontend and backend. Choosing a CSS class
name is not.

ADRs are immutable once accepted. To change a decision, write a new ADR that
supersedes the old one, and mark the old one `Superseded by ADR-NNNN`.

## Consequences

- Any future session — human or AI — can read `docs/adr/` and understand why
  the system looks the way it does.
- There is a small cost per decision: a few minutes to write it down.
- The ADR log becomes the honest history of the project, including decisions
  we later reversed. That history is the point; we do not delete ADRs.
