# Universal Code Quality

> NEVER flag or remove commented-out code — it may be intentionally preserved.

## Naming

- Clear, pronounceable, searchable. No single-letter except loop counters.
- Same vocabulary for same concept everywhere.
- No magic numbers/strings inline — named constants.

## Functions

- One thing. If name needs "and", split.
- <= 2-3 params; options object beyond.
- No flag parameters — split into two functions.
- Short. If long, doing too much.

## Side effects

- Prefer pure. Never mutate input arguments.
- Centralize I/O (file, DB, network).
- All I/O handles failure explicitly; no silent swallows.

## Edge cases

- Null/empty handled. Data existence is not assumed.
- Large datasets: paginate, stream, or size-limit.
- Validate + size-limit user input server-side.

## Security

- No raw user input in DOM / SQL / templates.
- PII never in URLs, query params, or logs.
- Sensitive data minimal in API responses.

## Tests

- New code has tests. Bug fixes prove the bug is fixed.
- Cover edge cases: empty, null, zero, negative, huge, wrong type.
- One concept per test. Cover error paths, not only happy.
