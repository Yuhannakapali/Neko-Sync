# CLAUDE.md — Web (`apps/web`)

This file provides guidance to Claude Code (claude.ai/code) when working in the Next.js frontend. For the monorepo big picture, see the root `CLAUDE.md`.

## Scope

Next.js **15** app (App Router, `src/app/`) on React **19**, styled with Tailwind, built and run through **Nx 23.2.1**. TypeScript **5.9**. The app is a thin skeleton today: a landing page, root layout, and one API route.

## Commands

Run from the repo root (Nx resolves the `web` project). `node_modules` is **not committed** — run `npm install` first.

```bash
npm install                         # once, before any nx command
npx nx serve web                    # dev server (next dev) on http://localhost:3000
npx nx run web:start                # serve the production build (next start)
npx nx build web                    # production build
npx nx lint web                     # ESLint 9 (flat config)
npx tsc --noEmit -p apps/web/tsconfig.json   # typecheck (see note below)
npx nx test web                     # currently broken — see "Current known breakages"
```

## Structure & conventions

```
apps/web/
├── src/app/
│   ├── layout.tsx          # root layout + metadata
│   ├── page.tsx            # landing page ("use client"; fetches /api/health)
│   ├── globals.css         # Tailwind entry
│   └── api/health/route.ts # route handler returning JSON via Response.json()
├── next.config.js          # rewrites /api/* → NEXT_PUBLIC_API_URL; Nx-aware
├── tsconfig.json           # paths: "@/*" → "./src/*"; moduleResolution "bundler"
├── tailwind.config.js
├── eslint.config.mjs       # ESLint 9 flat config (see below)
└── project.json            # explicit Nx targets: export/test (+ serve port override)
```

- **Backend proxy:** `next.config.js` rewrites `/api/:path*` to `NEXT_PUBLIC_API_URL` (default `http://localhost:8080/api/v1`). Client code calls same-origin `/api/...`; the rewrite forwards to the Go backend. Set `NEXT_PUBLIC_API_URL` to point elsewhere.
- **Path alias:** import app code as `@/…` (maps to `src/`).
- **Type & lint checks are deferred to Nx, not Next.** `next.config.js` sets `typescript.ignoreBuildErrors` and `eslint.ignoreDuringBuilds` to `true`, so `nx build web` will **not** catch type errors — run `tsc --noEmit` (above) and `nx lint web` explicitly.

## ESLint (flat config)

ESLint **9** with **flat config** only — there are no `.eslintrc.*` files. `apps/web/eslint.config.mjs` extends the root `eslint.config.mjs` conceptually and pulls in `next/core-web-vitals` through a `FlatCompat` shim (kept because that shared config has no native flat entry point). Add rules/overrides as flat config objects; don't reintroduce eslintrc.

## Nx targets

`build`, `serve`, `start`, `serve-static` and `lint` are **inferred** by the `@nx/next/plugin` and `@nx/eslint/plugin` entries in the root `nx.json` — they run `next build`, `next dev`, `next start` and `eslint .` in `apps/web`. Run `npx nx show project web` to see the resolved targets. `project.json` keeps only overrides (`serve` port 3000, extra `build` outputs) plus two executor-based targets: `export` (`@nx/next:export` — deprecated, removed in Nx 24, not handled by `convert-to-inferred`) and `test` (`@nx/jest:jest`, broken — see below). Don't reintroduce `@nx/next:build`/`@nx/next:server`/`@nx/eslint:lint`, and don't add a brace-glob `args` to `lint`: the task shell expands `**/*.{ts,tsx,js,jsx}` into separate patterns and ESLint fails on `**/*.jsx`.

## Current known breakages

- **`nx test web` fails.** `project.json`'s `test` target points at `apps/web/jest.config.ts`, which does not exist (no Jest config, setup, or test files). The target is unrunnable until Jest is configured. Typecheck via `tsc --noEmit` in the meantime.
