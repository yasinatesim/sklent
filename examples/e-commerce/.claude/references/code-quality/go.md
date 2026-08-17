# Go Code Quality

> NEVER flag commented-out code.

## Naming

- Package: short lowercase, no underscores.
- Exported `PascalCase`, unexported `camelCase`.
- Acronyms fully uppercase: `HTTPClient`, `URL`, `CORS`, `ID`.
- Project exception: `UPPER_SNAKE_CASE` allowed for package-level constants and immutable vars (the `var-naming` revive rule is disabled for this reason).

## Errors

- Every error handled. Never `_ = err` without comment.
- `errors.Is` / `errors.As` — not string comparison.
- Sentinel errors for well-defined modes: `var ErrX = errors.New("...")`.
- Wrap with `%w`: `fmt.Errorf("context: %w", err)`.
- Panic only for unrecoverable programmer bugs.

## Concurrency

- Share memory by communicating (channels).
- Always provide `context.Context` for cancellation.
- `sync.WaitGroup` or `errgroup` for fan-out.
- Run tests with `-race`.

## Traps

- nil slice vs empty slice: both `len == 0` but `nil != []T{}`. Use `len()` for checks.
- Interface holding typed nil is NOT nil.
- `range` copies values; use index when mutating slice elements.

## Testing

- Table-driven for 3+ similar cases.
- `t.Helper()` in helpers so failure lines point to caller.
- One concept per test case.
