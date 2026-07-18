import js from "@eslint/js";
import globals from "globals";

export default [
    js.configs.recommended,
    { languageOptions: {
            parserOptions: {
                ecmaVersion: "latest",
                sourceType: "module"
            },
            globals: { ...globals.node, ...globals.es6 }
        } },
    {
        rules: {}
    },
    {
        ignores: [
            "dist",
            "build",
            ".next",
            "tmp",
            "coverage"
        ]
    }
];
