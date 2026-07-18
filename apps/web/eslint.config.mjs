import { FlatCompat } from "@eslint/eslintrc";
import { dirname } from "path";
import { fileURLToPath } from "url";
import js from "@eslint/js";

const compat = new FlatCompat({
  baseDirectory: dirname(fileURLToPath(import.meta.url)),
  recommendedConfig: js.configs.recommended,
});


export default [
    ...compat.extends("next/core-web-vitals"),
    { languageOptions: { parserOptions: {
                ecmaVersion: "latest",
                sourceType: "module"
            } } },
    {
        rules: {
            "@next/next/no-html-link-for-pages": "off",
            "react/jsx-key": "error",
            "react/no-unescaped-entities": "warn"
        }
    },
    {
        ignores: [
            ".next/**/*"
        ]
    }
];
