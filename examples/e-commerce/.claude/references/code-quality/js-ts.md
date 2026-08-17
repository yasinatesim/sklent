# JavaScript / TypeScript Code Quality

> NEVER flag commented-out code.

## Variables

- `const` first, then `let`. Never `var`.
- Default params over `|| fallback`.

## Equality

- Always `===`, never `==` (except `null == undefined` explicitly).
- `Number.isNaN` over `x === NaN` (which is always false).
- `Array.prototype.sort()` requires a comparator — default is lexicographic.

## Async

- Always await or catch. No floating promises.
- No `new Promise(async (resolve) => {...})` antipattern.
- Cancel via AbortController for network fetches.

## Types

- No `any`. Use `unknown` with type guards.
- Catch clauses typed as `unknown`.
- Function return types declared for exported functions.

## React

- No `React.*` namespace — destructure hooks.
- No hooks called conditionally.
- `key` prop stable; never array index for reorderable lists.
- `{count && <C />}` renders `0` when falsy — use `{count > 0 && ...}`.
- `useEffect` cleanup is mandatory when effects subscribe or set timers.

## Style

- `UPPER_SNAKE_CASE` for constants and `as const` objects.
- `PascalCase` components, `camelCase` functions/vars.
- No inline styles — `*.module.scss` only.
- Every exported type/interface in `<module>/types/<name>.types.ts` (or `src/shared/types/` across modules); non-exported local types stay in their file.
