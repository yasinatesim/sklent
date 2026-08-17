# Module Migration Manifest

Target: **route = thin shell, logic = feature module, every component in its own folder with its
style beside it.**

Canonical rules: `.claude/references/project-structure.md`. This file is the checklist that applies
them, and the worked example of how to drive a structure migration to zero.

## Enforcement

| File | Rule | Checks |
|---|---|---|
| `web/scripts/projectStructure.mjs` | `project-structure/folder-structure` | folder/file naming, style placement, eponymous component folders |
| `web/scripts/independentModules.mjs` | `project-structure/independent-modules` | import direction, cross-module edges |
| `web/eslint.config.mjs` | `no-restricted-syntax` | exported types outside `types/` |

**✅ MIGRATION COMPLETE.** All three rules are `error`; a new violation fails the build.

```
                                        baseline   now
project-structure/folder-structure          11       0
project-structure/independent-modules       50       0
no-restricted-syntax                        14       0
```

## Target layout

```
web/src/
  app/                      route shell only — page.tsx is a one-line re-export
  features/
    admin/<module>/         admin panel features
    web/<module>/           public site features
  shared/                   constants, types, helpers, hooks, stores, ui, layout, styles, test
```

Module anatomy: `page.tsx`, `constants.ts`, plus the fixed folder set
`components | views | api | hooks | helpers | store | styles | types | assets | __tests__`.

## What the migration actually moved

| From | To |
|---|---|
| `src/constants`, `src/i18n`, `src/stores/{modal,toast,theme}` | `src/shared/{constants,i18n,stores}` |
| `src/components/{icons,Toast}`, `src/components/layout/*` | `src/shared/{ui,layout}` |
| `src/lib/api.ts` (a grab-bag) | split into `shared/helpers/api/client.ts`, `shared/helpers/money.ts`, `shared/types/catalog.types.ts`, and four feature `api/` files |
| `src/lib/*-api.ts` | the owning module's `api/` |
| `src/stores/{auth,cart}Store` | `features/web/{auth,cart}/store/` |
| `src/components/*` | `features/{web,admin}/<module>/components/<Comp>/` with an `index.ts` each |
| Components living in `app/[locale]/**` | their feature module; the route keeps a one-line import |
| Route-level `page.module.scss` consumed by a component | the module's `styles/` |
| Exported types in stores/api files | the module's `types/*.types.ts` |
| `MODAL` / `THEME` consts + their derived types | `shared/constants/` (a derived type stays with its const) |

## How to run a phase without breaking the build

1. Move by pattern, not by file — one whole pattern per batch.
2. Repair imports by **resolving** each specifier against the new tree, not by blind find-replace.
3. **Exclude `index.ts` files from bulk rewrites** — a relative→alias rewrite turns a folder's own
   barrel into a self-reference (`Circular definition of import alias 'default'`).
4. Run `npm run type-check` after **every** batch. Never save it for the end.
5. Run the full lane (`lint`, `type-check`, `test`) before starting the next phase.

## Verification

```bash
cd web && npm run lint && npm run type-check && npm test
```
