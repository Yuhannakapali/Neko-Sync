# Step 3b — Instance ↔ Hub connector

Status: **proposed** (2026-09-30) · Build-order step 3, second half. First half
(scanner → `MediaFile` → read-only API) is in `apps/instance`.

## Goal

Prove the edge-resolves-content loop:

```
Instance scans ─▶ reports files ─▶ Hub matches to Work/Child ─▶ ContentReference
Client asks Hub "play S1E3" ─▶ Hub returns ranked sources + signed stream URL
Client streams bytes directly from the Instance ─▶ Instance verifies Hub token
```

The Hub never sees bytes or file paths. The Instance never sees Hub accounts or
passwords.

## Decisions (agreed 2026-09-30)

1. **Pairing: OAuth 2.0 Device Authorization Grant (RFC 8628).** The Instance
   shows a short code; the owner approves it on the Hub. The Instance then holds
   a rotating refresh token and short-lived access tokens. Rejected: static API
   token (manual secret, never expires), mTLS (cert lifecycle on home servers,
   breaks behind TLS-terminating proxies such as Cloudflare).
2. **Matching: metadata provider first, local-only Work as fallback.** Lookup
   runs on the Hub (API keys and cache in one place). The Instance's own
   planned TMDB matcher (standalone mode) still runs; when it has matched a
   file, it sends the provider IDs and the Hub skips the search.
3. **Streaming: Hub-signed, asymmetric stream tokens.** The Hub signs with an
   Ed25519 private key; Instances verify with the Hub's public key (JWKS).
   Tokens are bound to one Instance and one file. No shared secret exists, so a
   Hub database leak cannot mint stream tokens. The Instance keeps its own local
   auth for standalone use and accepts either.

## Hub

### New feature package `internal/instance`

```go
type Instance struct {
    entity.BaseEntity
    OwnerID    entity.UUID // the user who approved pairing
    Name       string
    PublicURL  string      // how clients reach it, e.g. https://media.example.com
    LastSeenAt *time.Time
    RevokedAt  *time.Time
}
```

Migration `000003_instances`:

- `instances` (above).
- `instance_device_codes`: `device_code_hash` (PK, SHA-256), `user_code`
  (unique while pending), `status` (`pending|approved|denied|expired|used`),
  `instance_id` (set on approval), `expires_at`, `last_polled_at`.
- `instance_refresh_tokens`: `token_hash` (PK, SHA-256), `instance_id`,
  `expires_at`, `replaced_by` (rotation chain; reuse of a replaced token revokes
  the whole Instance — standard refresh-token reuse detection).

Secrets (device codes, refresh tokens) are stored **only as SHA-256 hashes**.
They are 32 random bytes, so an unsalted fast hash is appropriate (unlike
passwords).

### Pairing endpoints (RFC 8628 shapes)

| Method & path                           | Caller            | Does                                                                                                                                                                       |
| --------------------------------------- | ----------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `POST /api/v1/instances/device/code`    | Instance, no auth | Returns `device_code`, `user_code` (8 chars, no ambiguous letters), `verification_uri`, `expires_in` (600), `interval` (5). Rate-limited per IP.                           |
| `POST /api/v1/instances/device/approve` | Owner JWT         | Body `{user_code, name, public_url}`. Creates the `Instance`, marks the code approved.                                                                                     |
| `POST /api/v1/instances/token`          | Instance          | `grant_type=urn:ietf:params:oauth:grant-type:device_code` → `authorization_pending` / `slow_down` / `expired_token` / tokens. `grant_type=refresh_token` → rotated tokens. |
| `GET /api/v1/instances`                 | Owner JWT         | List own Instances.                                                                                                                                                        |
| `DELETE /api/v1/instances/{id}`         | Owner JWT         | Revoke: sets `RevokedAt`, deletes refresh tokens and the Instance's references.                                                                                            |

Instance access tokens: HS256 JWT from the existing `JWTManager` family but a
distinct audience (`aud: "instance"`, `sub: <instanceID>`, 15 min). A separate
`httpx.InstanceAuth` middleware accepts only that audience, and user
`AuthMiddleware` must reject it — tests cover both directions.

### Reporting endpoint

`PUT /api/v1/instance/libraries/{libraryID}/files` (Instance token) — the
**complete** current contents of one Instance library, replacing the previous
report for that library:

```json
{
  "kind": "anime",
  "files": [
    {
      "id": "<MediaFile ID>",
      "provider_ids": { "anilist": "16498" },
      "title": "Shingeki no Kyojin",
      "year": 2013,
      "season": 1,
      "episode": 3,
      "duration": 1440,
      "quality": "1080p h264"
    }
  ]
}
```

No path, size or stream detail is sent — only what matching needs.

For each file the Hub:

1. Resolves the **Work** (see Matching).
2. Resolves the **Child** when `season`/`episode` are present:
   `GetChildByOrdinal(work, season, ChildSeason, n)` then the episode under it,
   creating either if missing. A movie has no child (`ChildID` nil).
3. `Upsert`s a `ContentReference{UserID: owner, Source: instance, Locator}` with
   `Locator = "instance:<instanceID>/<libraryID>/<fileID>"`.

Then it deletes the owner's references whose locator starts with
`instance:<instanceID>/<libraryID>/` and were not in this report (files removed
from disk). The whole report runs in one transaction. Response: counts of
matched, created-local, and failed files, with reasons for failures.

Limit: 10 000 files per request (a larger library is reported in pages via
`?page=N&final=true`; only the final page triggers the stale-reference delete).

### Matching (`internal/metadata`)

```go
type Provider interface {
    Name() string                      // "anilist", "tmdb", ...
    Kinds() []work.Kind                // which kinds it can answer
    Search(ctx context.Context, kind work.Kind, title string, year *int) ([]Candidate, error)
    Fetch(ctx context.Context, id string) (*work.Work, []*work.WorkChild, error)
}
```

Order per file:

1. **Provider IDs sent** → `work.GetByProviderID`; on miss, `Fetch` from that
   provider and store the Work (and its seasons/episodes).
2. **No IDs** → `Search` the providers for the kind. Accept a candidate only if
   its normalized title (or any localized title) equals the reported one and
   the year matches when both are known. Otherwise treat as no match — a wrong
   automatic match is worse than a local-only Work, which a later manual match
   can fix.
3. **No match** → a local-only Work with `ProviderIDs = {"local":
"<instanceID>:<normalized title>:<year>"}`, so re-reports reuse it instead of
   creating duplicates, and it stays private to that Instance's owner.

Providers in this step: **AniList** (anime, manga; GraphQL, no key) and
**TMDB** (movie, series; `TMDB_API_KEY`, provider disabled when unset).
MusicBrainz waits until the Instance scans music. Search results are cached in
memory for 24 h (including misses) so a re-scan does not re-query; fetched Works
are cached permanently in `works` by provider ID. Outbound calls use a 10 s
timeout and per-provider rate limits.

### Playback resolution and stream tokens

`POST /api/v1/playback/resolve` (user JWT), body `{work_id, child_id?}`:

- Loads the caller's references, `reference.Rank`s them.
- For each `instance` reference whose Instance is not revoked, adds
  `stream_url = <PublicURL>/stream/<fileID>?token=<stream token>`.
- Other sources are returned as-is (their locator is already a URL or device
  hint).

Stream token: EdDSA (Ed25519) JWT, claims
`iss=<hub URL>`, `aud=instance:<instanceID>`, `sub=<userID>`, `fid=<fileID>`,
`exp`. Lifetime is **the file's duration + 30 min, capped at 6 h**, because a
player makes range requests for the whole viewing; a client that gets `401`
mid-play re-resolves. Only the owner's own references resolve in this step
(sharing is later).

Keys: `HUB_SIGNING_KEY` (PEM Ed25519 private key; startup fails if missing when
the instance feature is enabled). Public keys served at
`GET /.well-known/jwks.json` with a `kid`, so keys can rotate: publish the new
key, sign with it, retire the old one after the maximum token lifetime.

## Instance (`apps/instance`)

New `internal/application/hub` (connector) and `internal/interfaces/http`
additions:

- **Linking.** `nekosync-instance link --hub https://hub.example --public-url
https://media.example.com` (and `POST /api/hub/link` for a future UI) runs the
  device flow: prints `Go to <uri> and enter WDJB-MJHT`, polls at `interval`,
  and stores the refresh token and Instance ID in the SQLite `settings` table.
  `unlink` deletes them.
- **Token handling.** Refreshes the access token before expiry; on
  `invalid_grant` (revoked) it drops the link and logs clearly.
- **Reporting.** After every completed scan, and hourly, `PUT`s each library.
  Failures back off exponentially; the Instance keeps working unlinked.
- **Stream endpoint.** `GET /stream/{fileID}` serves the file with
  `http.ServeContent` (Range support). Accepts either
  - a Hub stream token: signature via cached JWKS (refetched on unknown `kid`,
    at most once a minute), `aud` = this Instance, `fid` = `{fileID}`, not
    expired; or
  - the Instance's own local credential (the planned standalone API token /
    HMAC URL), so it keeps working with no Hub.
- The stream response never includes the file path.

## Security checks (each gets a test)

- Device code and user code expire; a used code cannot be reused; polling
  faster than `interval` returns `slow_down`.
- Approving requires a user JWT; an Instance token cannot approve or call user
  routes, and a user token cannot call Instance routes.
- Refresh-token reuse revokes the Instance.
- A revoked Instance cannot report, and resolve stops issuing its URLs.
- An Instance can only replace references under its own locator prefix and
  only for its owner.
- Stream token for file A is rejected for file B; for Instance X rejected by Y;
  expired, `alg: none`, HS256-with-public-key and wrong-issuer tokens rejected.
- Reported titles are length-limited; provider responses are size-limited.

## Build order inside this step

1. Hub `instance` package + migration + device flow + instance auth middleware.
2. Hub reporting endpoint with local-only matching only (no providers yet).
3. Instance `link` command + token refresh + reporter. **First end-to-end
   checkpoint:** scan → report → references visible via `ListByUser`.
4. Hub Ed25519 signing, JWKS, `playback/resolve`.
5. Instance `/stream/{id}` with Hub-token verification. **Second checkpoint:**
   resolve on the Hub, stream from the Instance with `curl -r`.
6. `metadata` package with AniList, then TMDB.

Each part lands with unit tests and, where it touches Postgres, `pgtest`
integration tests; `make test-integration` and the Instance's `go test -race`
stay green.

## Out of scope

Relay / NAT traversal (the owner supplies a reachable `PublicURL`), sharing
media with other users, HLS/transcoding, web client UI, MusicBrainz, manual
match UI, multiple owners per Instance.
