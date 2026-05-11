---
name: melodihub-conventions
description: Use when changing or reviewing MelodiHub Go code under internal/domain, internal/providers, internal/contracts, or feature packages such as search, artist, album, and carousel; especially for architecture boundaries, TargetFromSource conversion naming, API error handling, root domain types, feature folder structure, and targeted Go test commands.
---

# MelodiHub Conventions

Use this skill as a project-local map before changing MelodiHub internals. Keep it lightweight: load only the reference that matches the task, then verify against the current files before editing.

## Quick Start

1. Identify the touched layer: root domain, provider, contract, feature service, feature handler, cache, or route wiring.
2. Read the matching reference:
   - `references/architecture.md` for folder boundaries and data flow.
   - `references/errors.md` for Melodi API errors and provider error mapping.
   - `references/naming.md` for `TargetFromSource` and model naming conventions.
   - `references/root-domain.md` for shared root domain concepts only.
   - `references/testing.md` for focused verification commands and known suite caveats.
3. Inspect the current live files named by the reference before changing code. The reference is a map, not a replacement for source.

## Project Rules

- Keep feature-specific value objects discoverable in their feature package. Do not assume names or behavior from this skill; inspect files like `internal/search/query.go` or `internal/search/limit.go` when working in that feature.
- Keep root domain objects in `internal/domain` only when they are genuinely shared across features or providers.
- Prefer feature-owned conversion and error mapping over generic helpers when provider semantics need feature context.
- Keep handlers explicit at response-write boundaries; do not hide `c.JSON` behind broad helpers unless the project already moved that boundary.
- Do not fix unrelated feature failures while applying a scoped convention change.
