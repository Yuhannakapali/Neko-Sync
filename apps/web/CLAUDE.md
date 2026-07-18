# CLAUDE.md — Web (`apps/web`)

This file provides guidance to Claude Code (claude.ai/code) when working in the Next.js frontend. For the monorepo big picture, see the root `CLAUDE.md`.

## Scope

Next.js **15** app (App Router, `src/app/`) on React **19**, styled with Tailwind, built and run through **Nx 23**. TypeScript **5.9**. The app is a thin skeleton today: a landing page, root layout, and one API route.

## Commands

Run from the repo root (Nx resolves the `web` project). `node_modules` is **not committed** — run `npm install` first.

```bash
npm install                         # once, before any nx command
npx nx serve web                    # dev server on http://localhost:3000
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
└── project.json            # Nx targets: build/serve/export/test/lint
```

- **Backend proxy:** `next.config.js` rewrites `/api/:path*` to `NEXT_PUBLIC_API_URL` (default `http://localhost:8080/api/v1`). Client code calls same-origin `/api/...`; the rewrite forwards to the Go backend. Set `NEXT_PUBLIC_API_URL` to point elsewhere.
- **Path alias:** import app code as `@/…` (maps to `src/`).
- **Type & lint checks are deferred to Nx, not Next.** `next.config.js` sets `typescript.ignoreBuildErrors` and `eslint.ignoreDuringBuilds` to `true`, so `nx build web` will **not** catch type errors — run `tsc --noEmit` (above) and `nx lint web` explicitly.

## ESLint (flat config)

ESLint **9** with **flat config** only — there are no `.eslintrc.*` files. `apps/web/eslint.config.mjs` extends the root `eslint.config.mjs` conceptually and pulls in `next/core-web-vitals` through a `FlatCompat` shim (kept because that shared config has no native flat entry point). Add rules/overrides as flat config objects; don't reintroduce eslintrc.

## Nx targets

`project.json` defines `build`/`serve`/`export`/`test`/`lint`. `serve` runs on port 3000 and depends on `build`. These still use the classic `@nx/next:*` / `@nx/eslint:lint` executors, which **warn as deprecated** under Nx 23 and are removed in Nx 24 — see root `CLAUDE.md` for the `convert-to-inferred` migration.

## Current known breakages

- **`nx test web` fails.** `project.json`'s `test` target points at `apps/web/jest.config.ts`, which does not exist (no Jest config, setup, or test files). The target is unrunnable until Jest is configured. Typecheck via `tsc --noEmit` in the meantime.
