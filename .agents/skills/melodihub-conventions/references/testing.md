# Testing

Use focused package tests for scoped changes.

## Common Commands

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

## Known Caveat

At the time this skill was created, `go test ./...` was expected to fail on unrelated `internal/album/service.go` build errors. Do not fix album while performing a scoped domain/search convention change unless the user asks for that.

When reporting verification, separate targeted pass results from unrelated full-suite failures.
