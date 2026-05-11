# Root Domain

Use this when deciding whether a type belongs in `internal/domain`.

## What Belongs Here

Root domain types should be shared project concepts, not feature conveniences.

Current root domain concepts include:

- `Provider`: shared provider identity and parsing, currently Deezer and LastFM.
- `ID`: shared provider-qualified numeric identifier.
- `Image`: shared image DTO/value used across responses.
- `InfoType`: shared music entity type such as track, artist, or album.
- `ErrorCode`, `MelodiError`, and `APIError`: shared API error identity and response shape.

Before changing these, inspect the current files in `internal/domain`.

## What Does Not Belong Here By Default

Do not move feature-specific objects into root domain just because they are typed or validated.

Examples that should stay discoverable in feature packages unless multiple features need them:

- query parsers
- limit parsers
- feature response aggregations
- feature cache keys
- provider-to-feature conversion helpers
- feature-specific error mappings

If a type is only needed by one feature, keep it in `internal/<feature>`.

## Promotion Rule

Promote a feature type to `internal/domain` only when it is stable, shared by more than one feature or provider boundary, and has meaning independent of one endpoint workflow.
