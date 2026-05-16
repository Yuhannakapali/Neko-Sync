# NekoSync UI — Implementation Plan

## Context

The spec defines 13 pages for a self-hosted anime/movie/music streaming platform. The Next.js frontend (`apps/web/`) is a bare scaffold — only a landing page and health-check route exist. The Go backend has only user auth endpoints (register, login, profile, follow, devices); no content/media endpoints exist yet. This plan builds a complete UI in 8 phases using mock data for content pages until backend endpoints are ready.

---

## Key Findings

### Frontend state

* `apps/web/src/app/` contains only:

  * `page.tsx`
  * `layout.tsx`
  * `globals.css`
  * `api/health/route.ts`
* No components directory
* No API client or state management
* Tailwind is configured
* Alias: `@/* → ./src/*`
* Strict TypeScript enabled

### API Proxy

* `/api/* → http://localhost:8080/api/v1/*`

---

## Backend Endpoints Available

### Auth

* POST `/api/users/register`
* POST `/api/users/login`
* PUT `/api/users/profile` (auth required)
* POST `/api/users/follow` (auth required)
* POST `/api/users/devices` (auth required)

### Missing (Mock Required)

* Content/media endpoints
* Search
* History tracking
* Playlists
* Player progress

---

## Package Additions

Run from `apps/web/`:

### UI System

```bash
npx shadcn@latest init
```

(Choose: dark style, slate base, CSS variables, `@/components/ui`)

### State & Data

```bash
npm install zustand @tanstack/react-query
```

### Forms

```bash
npm install react-hook-form @hookform/resolvers zod
```

### Media

```bash
npm install video.js hls.js
npm install -D @types/video.js
```

### Drag & Drop

```bash
npm install @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities
```

### Utilities

```bash
npm install date-fns js-cookie
npm install -D @types/js-cookie
```

### Icons

```bash
npm install lucide-react
```

---

## Directory Structure

```
apps/web/src/
├── app/
│   ├── layout.tsx
│   ├── globals.css
│   ├── page.tsx
│   ├── (auth)/
│   │   ├── layout.tsx
│   │   ├── login/page.tsx
│   │   └── register/page.tsx
│   ├── (app)/
│   │   ├── layout.tsx
│   │   ├── dashboard/page.tsx
│   │   ├── library/page.tsx
│   │   ├── media/[id]/page.tsx
│   │   ├── player/[id]/page.tsx
│   │   ├── music/[id]/page.tsx
│   │   ├── playlists/page.tsx
│   │   ├── playlists/[id]/page.tsx
│   │   ├── search/page.tsx
│   │   ├── history/page.tsx
│   │   ├── settings/page.tsx
│   │   └── admin/page.tsx
│   └── api/health/route.ts
│
├── middleware.ts
├── components/
├── lib/
└── types/
```

---

## Architecture Overview

### Auth Flow

1. Login → store JWT in localStorage
2. Store cookie `neko_auth` for middleware
3. Hydrate Zustand store on app load
4. Middleware protects routes
5. API client injects Bearer token

---

## Mock Data Strategy

All content systems use mock mode:

```ts
if (process.env.NEXT_PUBLIC_USE_MOCK_API === 'true') {
  await delay(300);
  return mockData;
}
```

Switch to real backend later without UI changes.

---

## Theme System

* Dark-first UI
* Tailwind + shadcn variables
* Brand: violet (`#8b5cf6`)
* Surfaces: deep slate backgrounds

---

## Phased Implementation

### Phase 0 — Foundation

* Tailwind theme setup
* Types (`lib/types`)
* API client
* Mock layer
* Zustand stores
* shadcn setup
* Providers
* Middleware
* Common UI components

### Phase 1 — Auth + Landing

* Landing page redesign
* Login/Register pages
* AuthGuard

### Phase 2 — App Shell

* Sidebar
* Topbar
* Mobile nav
* App layout wrapper

### Phase 3 — Dashboard + Library

* MediaCard system
* Dashboard widgets
* Library grid + filters

### Phase 4 — Media + Search

* Media detail page
* Episode list
* Search page with debounce

### Phase 5 — Video Player

* video.js integration
* HLS support
* progress tracking

### Phase 6 — Music Player

* Audio player
* Queue system
* persistent NowPlaying bar

### Phase 7 — Playlists + History

* Playlist CRUD
* drag & drop reorder
* history grouping

### Phase 8 — Settings + Admin

* profile settings
* playback settings
* admin panel

---

## Critical Files

* `lib/types/api.ts` → base contracts
* `lib/api/client.ts` → fetch wrapper
* `stores/authStore.ts` → auth core
* `middleware.ts` → route protection
* `components/layout/Providers.tsx` → app bootstrap
* `components/media/MediaCard.tsx` → reused everywhere
* `tailwind.config.js` → design system

---

## Verification Checklist

### Phase 1

* Login returns JWT
* Cookie + localStorage set
* Protected routes redirect correctly

### Phase 2

* App shell renders on all routes

### Phase 3

* Library grid loads mock content

### Phase 4

* Media detail loads episodes
* Search debounced

### Phase 5

* Video plays local MP4
* Progress saved

### Phase 6

* Audio playback works
* NowPlaying persists

### Phase 7

* Playlist drag/drop works
* History grouped

### Phase 8

* Profile update works
* Admin route protected

---

## Build Rule

After each phase:

```bash
npx nx build web
```

Ensure no TypeScript errors before continuing.

---

## End

This plan establishes a full production-grade frontend architecture for NekoSync with mock-first development and seamless backend integration later.
