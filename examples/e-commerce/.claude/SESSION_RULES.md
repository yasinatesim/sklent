# SESSION RULES — read BEFORE writing code

These apply from the first line, not at code review. Every item caught in review was already
written here.

## Before writing (every time)

1. **Reuse before you write.** Grep for an existing component/hook/util that already covers the
   case and extend it. A near-duplicate is a rejection, not a style nit.
2. **"Make X work for Y" = WIDEN the existing X.** Never build a second X. A parallel `XOther`
   method set, a second route tree, a second migration or a second component is a revert. Before
   writing, note: (a) which endpoint/function/table already does this, (b) the single narrowest
   thing blocking it from serving the new case. Fix only (b).
3. **Check the sibling screen.** Two screens rendering the same real-world concept must not drift.
   Extract a shared component instead of writing it twice.
4. **New route / new migration / new top-level component:** state in one line why the existing one
   cannot be widened. Cannot answer without hand-waving? You are duplicating.
5. **Minimum code is the default.** No speculative abstraction. Do not rewrite a screen that works.
6. **Grep the constants folders before adding any constant.**

## Comments (hook blocks these)

- Default: no comment. Only when WHY is non-obvious. Never WHAT.
- **Max 1 line.** No 2+ line comment blocks anywhere.
- No issue/PR references in comments.

## Where does this file go?

Answer before creating any file. Canonical layout:
`.claude/references/project-structure.md`.

| Writing a… | Goes to | Never |
|---|---|---|
| page | `<module>/page.tsx` (+ a one-line re-export in the route) | logic in the route folder |
| component | `<module>/components/<Comp>/<Comp>.tsx` + `index.ts` | module root, route folder |
| status view | `<module>/views/<View>/<View>.tsx` | `views/<View>.tsx` without a folder |
| **hook** | **`<module>/hooks/use<X>.ts`** | **inside a component folder** |
| fetch/client fn | `<module>/api/<name>.ts` | inside a component folder |
| pure helper | `<module>/helpers/<verbObject>.ts` | module root, `views/` |
| zustand store | `<module>/store/<x>Store.ts` | anywhere else |
| **exported type** | **`<module>/types/<name>.types.ts`** | **inline in `api/*.ts`** |
| local non-exported type | stays in the file that uses it | a `types/` file nobody reads |
| style, one consumer | next to the consumer, same name | module root, `assets/` |
| style, 2+ in module | `<module>/styles/<name>.module.scss` | module root, `assets/` |
| used by 2+ modules | `src/shared/**` | a second copy |

Hard limits:

- **A component folder holds only:** `<Comp>.tsx`, its style, `index.ts`, nested component folders,
  `__tests__/`. A hook, an api call, a helper or a types file in there is always misfiled — it
  belongs to the module.
- **Module root holds only** `page.tsx`, `page.module.scss`, `constants.ts` plus the fixed folder
  set. Any other folder name at module root is a violation, not a judgement call.
- `src/app/**` is a route shell. `page.tsx` is a one-line re-export. No styles, components or types.
- Imports flow one way: `shared → feature → app`. Cross-module edges must be declared in
  `scripts/independentModules.mjs` with a reason.
- **Never invent a directory to make a file fit.** Run `npx eslint <new file>` before continuing.

## Frontend

- Arrow functions. Bottom of file: `const X = () => {...}; export default X;`
- One component per file, each in its own eponymous folder with an `index.ts` re-export.
- Async state: one `useState<State>` over `REQUEST_STATUS`. No parallel `isLoading`/`isError`.
- Render dispatch via an object map at **module scope**, values are **component references**, never
  JSX elements. Never define views inside the parent body.
- No inline styles. No `React.*` namespace. Modals via the modal store.
- Comparison string literals get an `UPPER_SNAKE_CASE` const.
- **No `export interface` / `export type` outside a `types/` folder.** `constants/` is exempt: a
  type derived from a const in the same file stays with its const.
- Import order is lint-enforced and auto-fixable — never hand-sort, run `eslint --fix`.
- **Splitting a fat component: move the state into a hook FIRST.** Extracting JSX while the state
  stays in the parent produces a 30+ prop signature, which is worse than the complexity you started
  with.

## Backend

- No aliased imports. Type definitions in the `models/` subpackage.
- Tests: black-box in `<module>/tests/`, white-box beside `export_test.go`.
- Every handler propagates the request context. Map dispatch over `switch`.
- Runtime values from env — never hardcode merchant/supplier ids or base URLs.
- **Never persist third-party tokens or session state in the DB.** In-process only; a restart
  forcing a fresh login is correct behaviour.

## Before finishing

- Run the full verify lane as the LAST step, after every edit:
  `cd web && npm run lint && npm run stylelint && npm run type-check && npm test`
  `cd api && go build ./... && go vet ./... && go test ./...`
  `npm run lint` does not run stylelint — run it explicitly.
- Dead code: `npx fallow dead-code` (TS/JS) and `go run golang.org/x/tools/cmd/deadcode@latest -test ./...`.
  Check the never-delete list in CLAUDE.md before removing anything.
- Dispatch `wtf-code-reviewer`, then push. Never merge a PR.
