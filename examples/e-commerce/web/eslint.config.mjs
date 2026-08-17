import { dirname } from "node:path";
import { fileURLToPath } from "node:url";

import { FlatCompat } from "@eslint/eslintrc";
import { projectStructureParser,projectStructurePlugin } from "eslint-plugin-project-structure";
import simpleImportSort from "eslint-plugin-simple-import-sort";
import sonarjs from "eslint-plugin-sonarjs";
import unusedImports from "eslint-plugin-unused-imports";

import { independentModulesConfig } from "./scripts/independentModules.mjs";
import { projectStructureConfig } from "./scripts/projectStructure.mjs";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({ baseDirectory: __dirname });

const config = [
  {
    // First: later blocks restore the real parser for ts/tsx.
    files: ["**/*"],
    languageOptions: { parser: projectStructureParser },
    plugins: { "project-structure": projectStructurePlugin },
    rules: { "project-structure/folder-structure": ["error", projectStructureConfig] },
  },
  ...compat.extends("next/core-web-vitals", "next/typescript"),
  {
    rules: {
      "react/no-multi-comp": ["error", { ignoreStateless: false }],
    },
  },
  {
    files: ["src/**/*.{ts,tsx}"],
    plugins: { "project-structure": projectStructurePlugin },
    rules: { "project-structure/independent-modules": ["error", independentModulesConfig] },
  },
  {
    files: ["src/**/*.{ts,tsx}"],
    plugins: { "unused-imports": unusedImports },
    rules: {
      "@typescript-eslint/no-unused-vars": "off",
      "unused-imports/no-unused-imports": "error",
      "unused-imports/no-unused-vars": [
        "error",
        { vars: "all", varsIgnorePattern: "^_", args: "after-used", argsIgnorePattern: "^_" },
      ],
    },
  },
  {
    // An exported type has consumers, so it belongs in the module's types/ folder.
    // constants/ own their derived `typeof X[keyof typeof X]` aliases; splitting those helps nobody.
    files: ["src/**/*.{ts,tsx}"],
    ignores: ["**/types/**", "**/__tests__/**", "**/constants/**", "**/constants.ts"],
    rules: {
      "no-restricted-syntax": [
        "error",
        {
          selector: "ExportNamedDeclaration > TSInterfaceDeclaration",
          message: "Exported interface belongs in <module>/types/<name>.types.ts (or src/shared/types/ across modules).",
        },
        {
          selector: "ExportNamedDeclaration > TSTypeAliasDeclaration",
          message: "Exported type alias belongs in <module>/types/<name>.types.ts (or src/shared/types/ across modules).",
        },
      ],
    },
  },
  {
    ...sonarjs.configs.recommended,
    files: ["src/**/*.{ts,tsx}"],
    ignores: ["**/__tests__/**"],
    rules: {
      ...sonarjs.configs.recommended.rules,
      // `${base}${qs ? `?${qs}` : ''}` is the idiomatic URL builder; nesting is not the problem.
      "sonarjs/no-nested-template-literals": "off",
      // Fires on ERROR_CODE.WEAK_PASSWORD-style enum keys, which hold codes, not secrets.
      "sonarjs/no-hardcoded-passwords": "off",
      // Math.random only feeds toast ids and sampling here, never a token.
      "sonarjs/pseudo-random": "off",
      "sonarjs/cognitive-complexity": "warn",
    },
  },
  {
    // Import order mirrors the layering: framework, packages, shared, feature, relative, styles.
    files: ["src/**/*.{ts,tsx}", "e2e/**/*.ts"],
    plugins: { "simple-import-sort": simpleImportSort },
    rules: {
      "simple-import-sort/exports": "error",
      "simple-import-sort/imports": [
        "error",
        {
          groups: [
            ["^react$", "^react/", "^next$", "^next/"],
            ["^@?\\w"],
            ["^@/shared/constants(/.*|$)"],
            ["^@/shared/types(/.*|$)"],
            ["^@/shared/(helpers|hooks|stores|utils)(/.*|$)"],
            ["^@/shared/(ui|layout|test)(/.*|$)"],
            ["^@/features/[^/]+/[^/]+/types(/.*|$)"],
            ["^@/features/[^/]+/[^/]+/(api|store)(/.*|$)"],
            ["^@/features/[^/]+/[^/]+/(hooks|helpers)(/.*|$)"],
            ["^@/features/[^/]+/[^/]+/(components|views)(/.*|$)"],
            ["^@/"],
            ["^\\u0000"],
            ["^\\.\\.(?!/?$)", "^\\.\\./?$", "^\\./(?=.*/)(?!/?$)", "^\\.(?!/?$)", "^\\./?$"],
            ["^.+\\.s?css$"],
          ],
        },
      ],
    },
  },
  {
    files: ["**/icons/**"],
    rules: { "react/no-multi-comp": "off" },
  },
  { ignores: [".next/**", "node_modules/**", "coverage/**", "projectStructure.cache.json"] },
];

export default config;
