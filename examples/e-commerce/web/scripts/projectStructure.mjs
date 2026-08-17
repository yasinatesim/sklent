// @ts-check
import { createFolderStructure } from 'eslint-plugin-project-structure';

// `.` is auto-escaped and `*` becomes a wildcard, so regex-y segments use `*`
// rather than `.+` (which would match a literal dot).
export const projectStructureConfig = createFolderStructure({
  ignorePatterns: [
    'node_modules/**',
    '.next/**',
    'coverage/**',
    'public/**',
    'test-results/**',
    'playwright-report/**',
  ],
  structure: [
    {
      name: 'src',
      children: [
        {
          name: 'app',
          children: [
            { ruleId: 'route_folder' },
            { ruleId: 'route_node' },
            { ruleId: 'route_asset' },
            { ruleId: 'tests_folder' },
          ],
        },
        {
          name: 'features',
          children: [{ name: '(admin|web)', children: [{ ruleId: 'module_folder' }] }],
        },
        { name: 'shared', children: [{ ruleId: 'free_folder' }, { name: '*' }] },
        { ruleId: 'tests_folder' },
        { name: 'middleware.ts' },
      ],
    },
    {
      name: 'e2e',
      children: [
        { name: '(admin|web)', children: [{ name: '{kebab-case}.spec.ts' }] },
        { name: '(helpers|fixtures)', children: [{ name: '{camelCase}.ts' }] },
        { name: '{camelCase}.ts' },
        { name: '{kebab-case}.(spec|setup).ts' },
      ],
    },
    { name: 'scripts', children: [{ name: '{camelCase}.mjs' }] },
    { name: '*' },
  ],
  rules: {
    free_folder: {
      name: '*',
      children: [{ ruleId: 'free_folder' }, { name: '*' }],
    },

    route_node: {
      name: '(page|layout|template|default|error|global-error|loading|not-found|route|sitemap|robots|manifest)(.tsx|.ts)',
    },
    route_asset: {
      name: '(icon|apple-icon|favicon|opengraph-image|twitter-image)*(.png|.ico|.svg|.jpg|.jpeg)',
    },
    route_folder: {
      name: '(\\(*\\)|\\[*\\]|@*|{kebab-case})',
      children: [
        { ruleId: 'route_folder' },
        { ruleId: 'route_node' },
        { ruleId: 'route_asset' },
        { ruleId: 'tests_folder' },
      ],
    },

    module_folder: {
      name: '(_shared|\\[*\\]|{kebab-case}|{camelCase})',
      children: [
        { name: 'page.tsx' },
        { name: 'page.module.scss' },
        { name: 'constants.ts' },
        { ruleId: 'types_folder' },
        { ruleId: 'module_folder' },
        { ruleId: 'components_folder' },
        { ruleId: 'views_folder' },
        { ruleId: 'api_folder' },
        { ruleId: 'hooks_folder' },
        { ruleId: 'helpers_folder' },
        { ruleId: 'store_folder' },
        { ruleId: 'styles_folder' },
        { ruleId: 'assets_folder' },
        { ruleId: 'tests_folder' },
      ],
    },

    components_folder: {
      name: '(components|ui)',
      children: [{ ruleId: 'component_folder' }],
    },
    views_folder: {
      name: 'views',
      children: [{ ruleId: 'component_folder' }, { ruleId: 'tests_folder' }],
    },
    component_folder: {
      name: '{PascalCase}',
      children: [
        { name: '{FolderName}.tsx', enforceExistence: ['index.ts'] },
        { name: '{FolderName}.module.scss' },
        { name: 'index.ts' },
        { ruleId: 'component_folder' },
        { ruleId: 'tests_folder' },
      ],
    },

    api_folder: {
      name: 'api',
      children: [{ name: '{camelCase}.ts' }, { ruleId: 'tests_folder' }],
    },
    hooks_folder: {
      name: 'hooks',
      children: [{ name: 'use{PascalCase}.ts' }, { ruleId: 'tests_folder' }],
    },
    helpers_folder: {
      name: 'helpers',
      children: [{ name: '{camelCase}(.ts|.json)' }, { ruleId: 'tests_folder' }],
    },
    store_folder: {
      name: 'store',
      children: [{ name: '{camelCase}Store.ts' }, { ruleId: 'tests_folder' }],
    },
    types_folder: {
      name: 'types',
      children: [{ name: '{camelCase}.types.ts' }],
    },
    styles_folder: {
      name: 'styles',
      children: [{ name: '{camelCase}.module.scss' }],
    },
    assets_folder: {
      name: 'assets',
      children: [{ name: '*(.png|.jpg|.jpeg|.svg|.webp|.woff2|.json)' }],
    },
    tests_folder: {
      name: '__tests__',
      children: [
        { name: '*(.test.ts|.test.tsx)' },
        { name: '(helpers|fixtures)', children: [{ name: '*(.ts|.tsx)' }] },
      ],
    },
  },
});
