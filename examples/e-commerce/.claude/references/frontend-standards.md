# Frontend Standards (Next.js + React 19)

Rejection criteria for `wtf-js-react`. Read before any web work.

## Components

- One component per file. A `.tsx` exports exactly one default component.
- Arrow functions only. Never `function` declarations.
- Default export at the bottom: `const X = () => {...}; export default X;`.
- Sub-components in their own files under `ComponentName/`.
- Tests in `ComponentName/__tests__/ComponentName.test.tsx`.

## State

- Status-based state for async: one `useState<State>` discriminated over a status enum. No parallel
  `isLoading` / `isError` booleans.
- Render dispatch via object map: `const views = { [STATUS.IDLE]: IdleView, ... }`. No `&&` chains,
  no ternary trees, no `switch` in JSX.
- Sub-views read `state` from closure, never as a prop.
- Global state via Zustand under `web/src/stores/`. Never React Context for state.
- Modals via a central `useModalStore`. Modal components do not call the store themselves.

## Styling / assets

- No inline styles. `*.module.scss` only. No Tailwind.
- SVG icons centralized under `web/src/components/icons/`. No inline SVG.

## React hygiene

- No `React.*` namespace. Destructure hook imports (`import { useState } from "react"`).
- Respect `exhaustive-deps`. Name every effect-internal async function.

## i18n

- All UI strings live in `web/src/i18n/messages/{tr,en}.json`. Never hardcode UI text.
- Admin uses the `admin` key; storefront uses domain keys.
