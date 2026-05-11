# Architecture

Use this as a map, then inspect current files before editing.

## Layers

- `internal/domain`: shared cross-feature types and API-facing error shape.
- `internal/providers/<provider>`: provider DTOs, provider clients, provider sentinel errors, provider response parsing.
- `internal/contracts`: interfaces that feature services depend on when a provider capability is consumed across a boundary.
- `internal/<feature>`: feature-owned parsing, service orchestration, response shaping, cache usage, routes, handlers, and feature-specific conversion.
- `internal/bootstrap`: app wiring, provider construction, and route registration.

## Feature Package Shape

Common feature files:

- `contracts.go`: feature-local service interfaces when needed.
- `feature.go`: feature assembly.
- `routes.go`: route registration.
- `handler.go`: HTTP parsing and response write boundary.
- `service.go`: feature orchestration, provider calls, cache use, response conversion.
- `cache.go`: cache keys and cache serialization for that feature.
- `*_model.go` or focused model files: request/response/domain models owned by the feature.

Not every feature needs every file. Match nearby feature conventions before adding files.

## Data Flow

Typical request flow:

```text
handler -> feature parser/value objects -> service -> contract/provider -> feature conversion -> response DTO
```

Keep provider details out of handlers. Convert provider data and provider errors inside the feature/service boundary where feature context is available.

## Discovery Rule

Root shared concepts are documented in `root-domain.md`. Feature-specific objects are intentionally not listed there; discover them from the feature package currently being changed.
