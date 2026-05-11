# Naming

Use this when adding conversion functions, response DTOs, feature parsers, or model files.

## Conversion Direction

Prefer `TargetFromSource` names:

```go
ResultFromDeezer(...)
ResponseFromResult(...)
SearchErrorFromDeezerError(...)
```

The name should say what is produced first and what it consumes second. Avoid generic names like `Convert`, `Build`, or `Map` when the source and target matter.

## Model Names

- Use response suffixes for HTTP DTOs, for example `ResultResponse`, `TrackSummaryResponse`, or `TopResultResponse`.
- Keep provider DTOs in provider packages.
- Keep feature response shaping inside the feature package.
- Avoid reviving older generic model names if the feature has moved to focused response files.

## Feature Parsers and Value Objects

Feature parsers should live near the feature behavior they protect. Do not document feature-local objects as global conventions in this skill; inspect the current feature files.

Search examples:

```text
internal/search/query.go
internal/search/limit.go
```

Use these as local examples only when working on search or a similar feature parser. Do not assume every feature needs equivalent types.

## File Naming

- Put feature error conversion in `internal/<feature>/errors.go`.
- Prefer focused files when models have distinct responsibilities, such as `track_summary.go`, `artist_summary.go`, or `top_result.go`.
- Match sibling files before adding a new naming pattern.
