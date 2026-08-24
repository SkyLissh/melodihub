# Melodihub 🎧

A clean-architecture **Go** API server that powers a music discovery experience — the backend counterpart to [Melodi]. It aggregates, normalizes, and serves music data (albums, artists, carousels) with caching, structured logging, and OpenAPI docs.

> **Stack:** Go 1.25 · Echo v5 · zerolog · swaggo/swagger · Valkey · resty · validator · samber/lo

## Why this project

I wanted to prove I'm not just a UI person. Melodihub is a full-service API layer: a provider-abstracted upstream client (Deezer), a domain that normalizes everything into stable internal types, per-feature caching, and OpenAPI-generated docs. It's designed to be swapped, tested, and observed — not a one-off scraper.

## Highlights

- **Clean architecture** — feature-based internal packages (`artist`, `album`, `carousel`) each with `service`, `handler`, `routes`, `contracts`, and their own `cache`. No leaky global state.
- **Provider abstraction** — an upstream adapter (`internal/providers/deezer`) hides the third-party API shape behind a `contracts` + domain interface. Add another provider without touching the features.
- **Caching** — per-feature Valkey caches (`valkey-go`) with cache-backed services, so hot endpoints don't hammer the upstream.
- **Structured logging** — `zerolog` with a centralized bootstrap logger; correlation-friendly and JSON-ready.
- **Validation** — `go-playground/validator` on inbound DTOs.
- **OpenAPI / Swagger** — swaggo docs generated from typed models, served via `echo-swagger`.
- **Feature tests** — real behavior tests (`feature_test.go`) per feature, plus mocking upstream clients.
- **Domain modeling** — a dedicated `internal/domain` with typed IDs (`Id`), `MelodiError`/`ErrorCode` for consistent error responses, and image/info-type primitives. Errors are *typed*, not stringly-typed.
- **DI/bootstrap** — providers wired in `internal/bootstrap`, environment via `core/env`, CORS + logger middleware composed once.

## Architecture

```
internal/
├── bootstrap/     # app assembly: DI, logger, cors, providers
├── core/          # env
├── domain/        # shared domain types (Id, MelodiError, ErrorCode, image, provider)
├── cache/         # client
├── contracts/     # api contracts
├── <feature>/     # artist, album, carousel — service/handler/routes/cache/feature tests
└── providers/     # upstream adapters (deezer)
```

## Getting started

```bash
# copy env and set DEZER_API_KEY / upstream creds
cp .env.example .env

go run ./cmd/...              # boot the server
go test ./...                 # run feature tests
```

Requires Go 1.25+. Swagger UI served at `/swagger`.

## What it demonstrates

- Idiomatic Go + a strongly-typed domain layer with proper error modeling
- Provider abstraction and feature caching
- Swagger/OpenAPI-first APIs
- Tested features and clean bootstrapping

---

*Built by [Alisson "SkyLissh" Hernandez] — Go backend, API design, and system architecture. This is a personal portfolio project.*
