# React / Next.js Code Quality

> Companion to js-ts.md. NEVER flag commented-out code.

## Component design

- 150-300 lines per component; hard limit 500.
- Functional + hooks. No class components.
- No business logic in JSX — extract to custom hooks or a feature-scoped `lib/` (see `frontend-standards.md` § File organization for the tier rules).
- Every component in its own eponymous folder (`components/X/X.tsx`), never flat. No per-component `index.ts` barrel.

## Server vs Client components

- Default server. Mark `'use client'` only when hooks/events are needed.
- Data fetching preferred in server components or route handlers; not in client component effects.

## Props & state

- Destructure props at signature; explicit TypeScript types.
- Derive state; don't duplicate computable values.
- `useMemo` / `useCallback` only where profiling shows benefit.

## Routing

- App Router. `src/app/<segment>/page.tsx`.
- Dynamic segments `[slug]` with `generateStaticParams` where applicable.

## Styles

- `ComponentName.module.scss` + `className={styles.x}`.
- Design tokens via `@use '@/styles/tokens.scss' as t;`.
- No Tailwind.

## Testing (Vitest + RTL)

- `render`, `screen.getByRole`.
- `userEvent` over raw `fireEvent`.
- Test behavior, not DOM shape.
- Mock at network boundary.
