# Project Structure — `web/`

Canonical source. `MODULE_MIGRATION.md` is the checklist that applies this doc.

One line: **`app/` = route shell, `features/{admin,web}/*` = feature modules, `shared/` = common layer. Imports flow one way: `shared → feature → app`.**

## 1. Top level

```
web/
  e2e/                     # Playwright. testDir. Never moves into a module.
    admin/  web/            # one grouping level, no more
  src/
    app/                   # route files ONLY
    features/
      admin/<module>/      # admin panel feature modules
      web/<module>/        # public site feature modules
    shared/                # cross-cutting layer
    middleware.ts
```

No `src/components/`. Legacy. Gets deleted.

Two areas under `features/`, not one flat plane: `admin` and `web` share almost no code and carry opposite constraints (admin is auth-gated and SEO-irrelevant; web is public and SEO-critical). The area folder makes that boundary lint-checkable.

## 2. `src/app/` — thin shell

Allowed files, nothing else:

| Allowed | Forbidden |
|---|---|
| `page.tsx` (re-export only) | `*.module.scss` |
| `layout.tsx`, `template.tsx` | component `.tsx` |
| `error.tsx`, `loading.tsx`, `not-found.tsx` | `types.ts`, `constants.ts` |
| `route.ts`, `sitemap.ts`, `robots.ts`, `manifest.ts` | `lib/`, `hooks/`, `api/` |
| `icon.png`, `apple-icon.png`, `opengraph-image.*` | `__tests__/` |

`page.tsx` is exactly this:

```tsx
export { default } from '@/features/admin/orders/page';
```

Need `metadata` / `generateMetadata` / `generateStaticParams`? Re-export those too:

```tsx
export { default, metadata } from '@/features/web/catalog/[slug]/page';
```

`layout.tsx` and `error.tsx` may hold a body, but pull styles and JSX from the module or `shared/layout`.

## 3. Module anatomy

Same template for `src/features/admin/<module>/` and `src/features/web/<module>/`:

```
<module>/
  page.tsx                       # canonical page
  page.module.scss               # ONLY if page.tsx is the sole consumer
  styles/<name>.module.scss      # 2+ consumers inside the module
  api/<name>.ts + __tests__/
  hooks/use<X>.ts + __tests__/
  helpers/<verbObject>.ts + __tests__/
  store/<x>Store.ts + __tests__/
  constants.ts
  types/<name>.types.ts          # every EXPORTED type lives here
  assets/                        # module-owned images/fonts. NO scss here.
  components/<Comp>/
    <Comp>.tsx
    <Comp>.module.scss
    index.ts                     # export { default } from './<Comp>';
    __tests__/<Comp>.test.tsx
  views/<View>/
    <View>.tsx
    <View>.module.scss
    index.ts
    __tests__/<View>.test.tsx
  [id]/  new/                    # nested route = sub-module, same anatomy
```

Rules:

1. **No loose files at module root.** Only `page.tsx`, `page.module.scss`, `constants.ts`. Everything else goes in a subfolder — types included.
2. **Every component in its own eponymous folder.** `components/BulkUpdateBar/BulkUpdateBar.tsx`. A bare `components/BulkUpdateBar.tsx` is a violation.
3. **Every component folder carries an `index.ts`** that re-exports its own component and nothing else: `export { default } from './<Comp>';`. Import the folder (`@/features/admin/orders/components/OrderStatusBadge`), not the file. The ban is on *aggregate* barrels that re-export several unrelated things or a hook alongside a component — a folder's own single-line self-export is required, lint-enforced via `enforceExistence`.
4. **`views/` holds status-dispatch views only.** The `XLoadingView` / `XErrorView` / `XSuccessView` that feed a `STATE_VIEWS` map. Everything else is `components/`.
5. **No ad-hoc grouping folders.** `orders/marketplace/`, `orders/site/`, `orders/combined/`, `orders/returns/` either drop into `components/` or become real sub-modules. At module root the only folders are `components|views|api|hooks|helpers|store|styles|types|assets|__tests__|[dynamic]|<segment>`. Creating any other folder name at module root is a violation, not a judgement call.
6. **A component folder holds only its own files:** `<Comp>.tsx`, `<Comp>.module.scss`, `index.ts`, nested component folders, `__tests__/`. A hook, an api call, a helper or a types file placed inside a component folder is always misfiled — those belong to the **module**, so they go in `<module>/hooks/`, `<module>/api/`, `<module>/helpers/`, `<module>/types/`. Being used by only one component today does not make it the component's file; the next consumer would have to reach into a component folder to get it.
7. **Nested sub-modules get the full anatomy.** `[id]/hooks/`, `new/api/` are fine — a sub-module is a module. The ban is on subfolders *inside a component folder*, not on depth.

## 3b. Type placement

An **exported** type or interface has consumers outside its file, so it is not a private detail of that file. It lives in a `types/` folder:

| Consumed by | Location |
|---|---|
| One module | `<module>/types/<name>.types.ts` |
| 2+ modules, or both areas | `src/shared/types/<name>.types.ts` |

`api/*.ts` files declare **no** exported types — they import them from `types/`. Same for components, hooks and helpers.

A **non-exported** local type stays where it is used: `type Props = {...}` inside a component file is correct and must not be moved. The line is the `export` keyword, not the size of the type.

Lint: `no-restricted-syntax` flags `export interface` / `export type` outside `**/types/**`.

## 4. SCSS placement

Three lines:

| Case | Location |
|---|---|
| One consuming component | Inside the component folder, same name: `Foo/Foo.module.scss` |
| Only `page.tsx` consumes | `<module>/page.module.scss` |
| 2+ consumers in the module | `<module>/styles/<name>.module.scss` |
| 2+ consumers across modules | `src/shared/styles/` or `src/features/admin/_shared/styles/` |
| Global / tokens / mixins | `src/shared/styles/` |

`assets/` is for binaries **only** (png, svg, woff2). SCSS never goes there — a CSS Module is code, not an asset.

A loose `<module>/orders.module.scss` at module root is a violation. Either it becomes `page.module.scss` (single consumer) or drops into `styles/` (many consumers).

## 5. Layers and import direction

```
shared  →  features/admin/_shared  →  features/admin/<module>  →  app
        →                              features/web/<module>    →  app
```

| Layer | May import |
|---|---|
| `src/shared/**` | `src/shared/**` only |
| `src/features/admin/_shared/**` | `src/shared/**`, `src/features/admin/_shared/**` |
| `src/features/admin/<m>/**` | `src/shared/**`, `src/features/admin/_shared/**`, own module |
| `src/features/web/<m>/**` | `src/shared/**`, own module |
| `src/app/**` | anything, but only as a re-export/shell |

**Sibling modules never import each other.** `features/admin/orders` → `features/admin/products` is forbidden. If two modules need the same thing, that thing gets promoted to `_shared/` or `shared/`. `features/web/*` ↔ `features/admin/*` is forbidden in both directions.

A need shared between `features/web/` and `features/admin/` means `src/shared/` — not a third copy.

## 6. Naming

| Thing | Form | Example |
|---|---|---|
| Module folder | kebab-case, exactly the route segment | `mp-order-deletions`, `home-rails`, `stock-tracking` |
| Component folder + file | PascalCase, eponymous | `OrderStatusBadge/OrderStatusBadge.tsx` |
| Hook | `use<X>.ts` | `useOrderFilters.ts` |
| Store | `<x>Store.ts` | `cartStore.ts` |
| helper file | camelCase, verb + object | `buildCardEntries.ts` |
| api file | camelCase resource name | `orders.ts` |
| SCSS | same name as its consumer | `OrderStatusBadge.module.scss` |
| Test | `__tests__/<same name>.test.tsx` | `__tests__/OrderStatusBadge.test.tsx` |
| Type file | camelCase domain + `.types.ts` | `order.types.ts` |

Module folders are kebab-case so `app/admin/mp-order-deletions/page.tsx` → `@/features/admin/mp-order-deletions/page` is a mechanical mapping. camelCase (`mpOrderDeletions`) turns that mapping into a mental translation step.

## 7. Test placement

| Kind | Location |
|---|---|
| Unit / component | Next to the code, `__tests__/<name>.test.ts(x)` |
| Route shell test | None. Shells are not tested; the module's `page.tsx` is. |
| E2E (Playwright) | `web/e2e/{admin,web}/<flow>.spec.ts` |

E2E stays out of modules: one spec walks several modules, so no single module owns it. Playwright `testDir: './e2e'`, vitest `exclude: ['e2e']` — the two runners stay apart.

## 8. Enforcement

`eslint-plugin-project-structure`:

- `web/projectStructure.mjs` — folder/file name schema (§3, §4, §6)
- `web/independentModules.mjs` — import direction + sibling isolation (§5)

Two sanctioned inversions, both declared explicitly in `independentModules.mjs`:

- `src/shared/ui/Modal/**` — the modal registry maps a modal type to its feature component, so it depends on every feature by definition. It is a composition root, not a shared primitive. Do not copy the pattern.
- `src/shared/test/**` and every `__tests__/**` — test code may import any layer; production code may not.

Extra native rules in `web/eslint.config.mjs`:

- `no-restricted-imports` → `@/components/**` banned (legacy, folder now deleted)
- `no-restricted-syntax` → exported type/interface outside `types/` (`constants/` is exempt: a `typeof X[keyof typeof X]` alias belongs with its const)
- `simple-import-sort/imports` + `/exports` → import order mirrors the layering; auto-fixable with `eslint --fix`

Dead-code sweeps run outside ESLint: `npx fallow dead-code` for TS/JS, `go run golang.org/x/tools/cmd/deadcode@latest -test ./...` for Go. fallow has no Go support.
- `react/no-multi-comp` → one component per file
