// @ts-check
import { createIndependentModules } from 'eslint-plugin-project-structure';

// {family_4} = "src/features/<area>/<module>" — fewer parts would let
// src/features/admin count as family and permit cross-module imports.
export const independentModulesConfig = createIndependentModules({
  extensions: ['.ts', '.tsx', '.scss', '.json'],
  pathAliases: { baseUrl: '.', paths: { '@/*': ['./src/*'] } },
  reusableImportPatterns: {
    shared: ['src/shared/**'],
    adminShared: ['src/features/admin/_shared/**'],
  },
  modules: [
    {
      name: 'Route shell (src/app)',
      pattern: 'src/app/**',
      allowImportsFrom: ['src/**'],
      errorMessage:
        '❌ src/app is a route shell only. Feature code lives in src/features/<area>/<module>.',
    },
    {
      name: 'Legacy src/components',
      pattern: 'src/components/**',
      allowImportsFrom: [],
      allowExternalImports: false,
      errorMessage:
        '❌ src/components/ is legacy and being deleted. New code goes in src/features/<area>/<module>/components or src/shared/ui.',
    },
    {
      name: 'Modal registry (composition root)',
      pattern: ['src/shared/ui/Modal/**', 'src/shared/types/modal.types.ts'],
      allowImportsFrom: ['src/**'],
      errorMessage:
        '❌ The modal registry maps a modal type to its feature component, so it depends on every feature by design. This is the only sanctioned inversion — do not copy the pattern elsewhere.',
    },
    {
      name: 'Test infrastructure',
      pattern: ['src/shared/test/**', 'src/**/__tests__/**'],
      allowImportsFrom: ['src/**'],
      errorMessage:
        '❌ Test code may import any layer; production code may not.',
    },
    {
      name: 'App shell (site chrome composes every public feature)',
      pattern: ['src/features/web/layout/**', 'src/shared/layout/**'],
      allowImportsFrom: ['{shared}', 'src/features/web/**'],
    },
    {
      name: 'Shared layer',
      pattern: 'src/shared/**',
      allowImportsFrom: ['{shared}'],
      errorMessage:
        '❌ src/shared/ imports from src/shared/ only. Depending on a feature layer is the wrong direction — move the code into shared/.',
    },
    {
      name: 'Admin shell (nav, guards and shared modals act on admin features)',
      pattern: 'src/features/admin/_shared/**',
      allowImportsFrom: [
        '{shared}',
        '{adminShared}',
        'src/features/admin/*/api/**',
        'src/features/admin/*/types/**',
        'src/features/web/auth/store/**',
      ],
      errorMessage:
        '❌ src/features/admin/_shared/ imports from shared and itself only. Depending on a feature module is the wrong direction.',
    },
    {
      name: 'Admin feature module',
      pattern: 'src/features/admin/*/**',
      // Declared cross-module edges: a coupon targets promotions, an order shows its
      // promotion, a product picks a category. Anything not listed here still fails.
      allowImportsFrom: [
        '{shared}',
        '{adminShared}',
        '{family_4}/**',
        // Declare each real cross-module edge here, one line of justification per entry.
        // An undeclared edge still fails — the rule is visibility, not prohibition.
      ],
      errorMessage:
        '❌ Sibling admin module import. Promote the shared piece to src/features/admin/_shared/ or src/shared/.',
    },
    {
      name: 'Web feature module',
      pattern: 'src/features/web/*/**',
      // Declared cross-module edges: a product card adds to cart, checkout places an
      // order, account edits the signed-in user. Anything not listed here still fails.
      allowImportsFrom: [
        '{shared}',
        '{family_4}/**',
        // Declared cross-module edges — one justification per line. An undeclared edge still fails.
        // Checkout reads the cart it is about to convert and posts the resulting order.
        'src/features/web/cart/(store|types)/**',
        'src/features/web/orders/(api|types)/**',
      ],
      errorMessage:
        '❌ Undeclared sibling web module import. Either promote the shared piece to src/shared/, or add the edge to allowImportsFrom with a one-line reason.',
    },
    {
      name: 'src root files',
      pattern: 'src/*',
      allowImportsFrom: ['src/**'],
    },
    {
      name: 'Outside src',
      pattern: [['**', '!src/**']],
      allowImportsFrom: ['**'],
    },
  ],
});
