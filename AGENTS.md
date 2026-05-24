# MelodiHub Conventions

Use this file as a project-local map before changing MelodiHub internals. Keep it lightweight: jump to the relevant section below, then verify against the current files before editing.

## Quick Start

1. Identify the touched layer: root domain, provider, contract, feature service, feature handler, cache, or route wiring.
2. Read the matching section in this file:
   - **Architecture** for folder boundaries and data flow.
   - **Errors** for Melodi API errors and provider error mapping.
   - **Naming** for `TargetFromSource` and model naming conventions.
   - **Root Domain** for shared root domain concepts only.
   - **Testing** for focused verification commands and known suite caveats.
3. Inspect the current live files named by the reference before changing code. This file is a map, not a replacement for source.

## Project Rules

- Keep feature-specific value objects discoverable in their feature package. Do not assume names or behavior from this file; inspect files like `internal/search/query.go` or `internal/search/limit.go` when working in that feature.
- Keep root domain objects in `internal/domain` only when they are genuinely shared across features or providers.
- Prefer feature-owned conversion and error mapping over generic helpers when provider semantics need feature context.
- Keep handlers explicit at response-write boundaries; do not hide `c.JSON` behind broad helpers unless the project already moved that boundary.
- Do not fix unrelated feature failures while applying a scoped convention change.

---

## Architecture

### Layers

- `internal/domain`: shared cross-feature types and API-facing error shape.
- `internal/providers/<provider>`: provider DTOs, provider clients, provider sentinel errors, provider response parsing.
- `internal/contracts`: interfaces that feature services depend on when a provider capability is consumed across a boundary.
- `internal/<feature>`: feature-owned parsing, service orchestration, response shaping, cache usage, routes, handlers, and feature-specific conversion.
- `internal/bootstrap`: app wiring, provider construction, and route registration.

### Feature Package Shape

Common feature files:

- `contracts.go`: feature-local service interfaces when needed.
- `feature.go`: feature assembly.
- `routes.go`: route registration.
- `handler.go`: HTTP parsing and response write boundary.
- `service.go`: feature orchestration, provider calls, cache use, response conversion.
- `cache.go`: cache keys and cache serialization for that feature.
- `*_model.go` or focused model files: request/response/domain models owned by the feature.

Not every feature needs every file. Match nearby feature conventions before adding files.

### Data Flow

Typical request flow:

```text
handler -> feature parser/value objects -> service -> contract/provider -> feature conversion -> response DTO
```

Keep provider details out of handlers. Convert provider data and provider errors inside the feature/service boundary where feature context is available.

### Discovery Rule

Root shared concepts are documented in **Root Domain** below. Feature-specific objects are intentionally not listed there; discover them from the feature package currently being changed.

---

## Naming

Use this when adding conversion functions, response DTOs, feature parsers, or model files.

### Conversion Direction

Prefer `TargetFromSource` names:

```go
ResultFromDeezer(...)
ResponseFromResult(...)
SearchErrorFromDeezerError(...)
```

The name should say what is produced first and what it consumes second. Avoid generic names like `Convert`, `Build`, or `Map` when the source and target matter.

### Model Names

- Use response suffixes for HTTP DTOs, for example `ResultResponse`, `TrackSummaryResponse`, or `TopResultResponse`.
- Keep provider DTOs in provider packages.
- Keep feature response shaping inside the feature package.
- Avoid reviving older generic model names if the feature has moved to focused response files.

### Feature Parsers and Value Objects

Feature parsers should live near the feature behavior they protect. Do not document feature-local objects as global conventions in this file; inspect the current feature files.

Search examples:

```text
internal/search/query.go
internal/search/limit.go
```

Use these as local examples only when working on search or a similar feature parser. Do not assume every feature needs equivalent types.

### File Naming

- Put feature error conversion in `internal/<feature>/errors.go`.
- Prefer focused files when models have distinct responsibilities, such as `track_summary.go`, `artist_summary.go`, or `top_result.go`.
- Match sibling files before adding a new naming pattern.

---

## Errors

Use this when changing handler/service/provider error behavior.

### Root Domain Error Model

Canonical files:

- `internal/domain/error_code.go`
- `internal/domain/melodi_error.go`
- `internal/domain/api_error.go`

Public API error response:

```go
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Cause   string    `json:"cause,omitempty"`
}
```

Current code set:

```text
INVALID_PARAM
INVALID_BODY
NOT_FOUND
PROVIDER_UNAVAILABLE
RATE_LIMITED
INTERNAL
```

`Cause` means safe public reason. Constructors like `InvalidParam`, `InvalidBody`, `NotFound`, `ProviderUnavailable`, and `RateLimited` expose safe causes. `Internal` must not expose a cause.

### Handler Boundary

Handlers should stay explicit:

```go
apiErr, status := domain.APIErrorFromMelodiError(err)
return c.JSON(status, apiErr)
```

Parse errors produced by feature parsers should become `domain.InvalidParam(err)` at the handler boundary when the parser message is intended to be public.

### Provider Error Conversion

Use feature-aware conversion names:

```go
func SearchErrorFromDeezerError(err error) error
```

Keep the `TargetFromSource` direction. Do not create broad provider-to-domain helpers until multiple features prove the same mapping is truly generic. Provider errors may need feature context.

For search, Deezer errors are mapped in `internal/search/errors.go`. Inspect that file before adding or changing mapping behavior.

### Service Boundary

Services should wrap provider failures into Melodi errors before returning to handlers:

```go
results, err := s.search.Find(ctx, query.String(), limit.Value())
if err != nil {
	return nil, SearchErrorFromDeezerError(err)
}
```

Keep retry behavior and provider transport behavior in provider/resty code, not in domain errors.

---

## Root Domain

Use this when deciding whether a type belongs in `internal/domain`.

### What Belongs Here

Root domain types should be shared project concepts, not feature conveniences.

Current root domain concepts include:

- `Provider`: shared provider identity and parsing, currently Deezer and LastFM.
- `ID`: shared provider-qualified numeric identifier.
- `Image`: shared image DTO/value used across responses.
- `InfoType`: shared music entity type such as track, artist, or album.
- `ErrorCode`, `MelodiError`, and `APIError`: shared API error identity and response shape.

Before changing these, inspect the current files in `internal/domain`.

### What Does Not Belong Here By Default

Do not move feature-specific objects into root domain just because they are typed or validated.

Examples that should stay discoverable in feature packages unless multiple features need them:

- query parsers
- limit parsers
- feature response aggregations
- feature cache keys
- provider-to-feature conversion helpers
- feature-specific error mappings

If a type is only needed by one feature, keep it in `internal/<feature>`.

### Promotion Rule

Promote a feature type to `internal/domain` only when it is stable, shared by more than one feature or provider boundary, and has meaning independent of one endpoint workflow.

---

## Testing

Use focused package tests for scoped changes.

### Common Commands

For search/domain error work:

```bash
go test ./internal/domain ./internal/search ./internal/artist
```

For provider changes:

```bash
go test ./internal/providers/deezer
go test ./internal/providers/lastfm
```

For full verification:

```bash
go test ./...
```

When reporting verification, separate targeted pass results from unrelated full-suite failures.
