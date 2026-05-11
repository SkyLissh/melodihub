# Error Handling

Use this when changing handler/service/provider error behavior.

## Root Domain Error Model

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

## Handler Boundary

Handlers should stay explicit:

```go
apiErr, status := domain.APIErrorFromMelodiError(err)
return c.JSON(status, apiErr)
```

Parse errors produced by feature parsers should become `domain.InvalidParam(err)` at the handler boundary when the parser message is intended to be public.

## Provider Error Conversion

Use feature-aware conversion names:

```go
func SearchErrorFromDeezerError(err error) error
```

Keep the `TargetFromSource` direction. Do not create broad provider-to-domain helpers until multiple features prove the same mapping is truly generic. Provider errors may need feature context.

For search, Deezer errors are mapped in `internal/search/errors.go`. Inspect that file before adding or changing mapping behavior.

## Service Boundary

Services should wrap provider failures into Melodi errors before returning to handlers:

```go
results, err := s.search.Find(ctx, query.String(), limit.Value())
if err != nil {
	return nil, SearchErrorFromDeezerError(err)
}
```

Keep retry behavior and provider transport behavior in provider/resty code, not in domain errors.
