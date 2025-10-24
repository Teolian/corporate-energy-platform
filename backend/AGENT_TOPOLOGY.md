# AGENT_TOPOLOGY.md
**Project:** Corporate Energy Platform — Monorepo Topology & Deployments  
**Audience:** AI coding agent (Claude/Cursor) + maintainer  
**Purpose:** Define the repository layout with **two frontends** (mock/live) and a **shared backend**, while staying consistent with `AGENT_TECH_SPEC.md`, `AGENT_BIZ_SPEC.md`, and `AGENT_INTEGRATION_GUIDE.md`.

> This document is **additive**: it does not change API contracts or business rules. It only defines how we organize code and deployments.
> If any ambiguity arises, **API contracts in `AGENT_TECH_SPEC.md §3` are canonical.**

---

## 1) Canonical Topology (Recommended)
We use a lightweight **monorepo** with workspaces and two separate frontends.

```
/apps
  /frontend-mock    # current demo; stays mock-only (no backend dependency)
  /frontend-live    # live-data app; connects to static artifacts or /api
  /backend          # Go/FastAPI services and/or serverless handlers
/packages
  /types            # shared DTOs/interfaces used across apps
  /ui               # shared UI components (optional, future)
  /config           # shared config (thresholds, area mappings, i18n stubs)
/pipeline           # GH Actions scripts; normalization; JSON Schemas
/docs               # AGENT_* specs (tech, biz, integration, topology)
```

**Package manager:** `pnpm` or `yarn` workspaces (agent may choose; prefer `pnpm`).  
**TS base config:** a root `tsconfig.base.json` and per-app `tsconfig.json` extending it.  
**Path aliases:** `@types`, `@ui`, `@config` point to `/packages/*`.

---

## 2) Workspace Setup (example with pnpm)
**Root `package.json`** (excerpt):
```json
{
  "private": true,
  "workspaces": [
    "apps/*",
    "packages/*"
  ]
}
```

**Root `pnpm-workspace.yaml`**:
```yaml
packages:
  - "apps/*"
  - "packages/*"
```

Each app keeps its own `package.json` and build scripts.

---

## 3) Frontend Apps
### 3.1 `apps/frontend-mock` (Stable Demo)
- Purpose: keep the current prototype **unchanged**.  
- Data mode: **MOCK only**; hardcode `VITE_DATA_MODE=mock`.  
- Deployment: Vercel Project A → Production URL for the demo.
- Allowed changes: dependency updates, visuals, content. **No assumption of backend.**

### 3.2 `apps/frontend-live` (Live + New Features)
- Purpose: connect to normalized artifacts `/public/data/...` and later to `/api/*`.  
- Data mode: **LIVE by default** with fallback to MOCK per widget.  
- Uses the **data client** and adapters described in `AGENT_INTEGRATION_GUIDE.md §3`.  
- Deployment: Vercel Project B → Production URL for the live app.
- Features such as **forecast overlay, settlement, reserve badges** are wired here (behind flags).

**Shared packages**
- `@types`: shared DTOs that match the API contracts (do not change schema names without updating `AGENT_TECH_SPEC.md`).  
- `@config`: thresholds, area mappings, source attribution strings.  
- `@ui` (optional): design-system components reused by both frontends.

---

## 4) Backend App (`apps/backend`)
- Implements adapters to public sources and exposes `/api` endpoints (or writes static artifacts) per `AGENT_TECH_SPEC.md`.  
- Can be split into Go (core) and FastAPI (ML) services.  
- Deployed separately (serverless functions or dedicated host).  
- Not required for `frontend-mock`; required/optional for `frontend-live` depending on Mode A/B.

---

## 5) Deployments (Vercel)
Create **two Vercel projects** from the same repo:
- **Project A** → Root: `apps/frontend-mock` → ENV: `VITE_DATA_MODE=mock`.
- **Project B** → Root: `apps/frontend-live` → ENV: `VITE_DATA_MODE=live`, `VITE_API_BASE` (when `/api` is available).

**Preview Deployments**
- Every PR/branch generates a preview URL for the affected app.  
- For backend PRs, test via `frontend-live` preview pointing to preview/staging API or static artifacts.  
- Protect `main` branches; require CI green.

---

## 6) Branching Model
- Long-lived **`feat/backend-core`** for backend integration.  
- Short-lived feature branches merge into `feat/backend-core`:  
  - `feat/adapter-tepco`, `feat/adapter-kansai`, `feat/adapter-jepx`, `feat/adapter-occto`,  
  - `feat/settlement-lite`, `feat/weather-solar`, `feat/forecast-ml` (optional).  
- **UI features for live app** in `apps/frontend-live`:  
  - `feat/ui-forecast-overlay`, `feat/ui-settlement`, `feat/ui-reserve-badges`, …  
- `apps/frontend-mock` changes are minimal and independent.

**Promotion to `main` (repo root) happens when:**
1) Contracts unchanged and validated;  
2) `frontend-live` renders preview without regressions;  
3) Fallback to MOCK works;  
4) Source attribution and timestamps present.

---

## 7) CI/CD Matrix
- Jobs by path filters:
  - changes in `/pipeline` → run normalization+schema checks; optionally commit `/public/data/*` artifacts.
  - changes in `apps/frontend-live/**` → build/test only live app.
  - changes in `apps/frontend-mock/**` → build/test only mock app.
  - changes in `apps/backend/**` → build/test backend; run adapter golden tests.
- Nightly pipeline cron to refresh public artifacts.
- Optional: single “status badge” JSON with latest artifact timestamps.

---

## 8) Environment & Config
- `VITE_DATA_MODE`: `"mock"` in Project A; `"live"` in Project B.  
- `VITE_API_BASE`: base URL for `/api` (Project B only; optional in Static Mode).  
- Thresholds and mappings live in `@config` and are shared by both apps.  
- Secrets are never committed; use Vercel env + CI secrets.

---

## 9) Compatibility with Existing Docs
- `AGENT_TECH_SPEC.md`: **API contracts unchanged.**  
- `AGENT_BIZ_SPEC.md`: value/flows unchanged.  
- `AGENT_INTEGRATION_GUIDE.md`: where it says “single app with mock/live modes,” **treat that as instructions for `apps/frontend-live`**; `apps/frontend-mock` simply stays in mock mode.

---

## 10) Cleanup / Migration Notes
- Keep `frontend-mock` as the legacy demo indefinitely (good for portfolio).  
- When live stabilizes, you may deprecate `frontend-mock` or redirect it to a read-only route in `frontend-live`.  
- Avoid duplicating business logic in both frontends; use `@types`/`@config` to share code and keep behavior aligned.

---

## 11) Agent Checklist
- [ ] Initialize workspaces and create folders as in §1.  
- [ ] Move current app to `apps/frontend-mock`.  
- [ ] Create `apps/frontend-live` from the mock app and wire `dataClient` + adapters.  
- [ ] Extract shared DTOs to `/packages/types` and update imports.  
- [ ] Add Vercel Project A (mock) and Project B (live) with proper ENV.  
- [ ] Implement CI path filters and preview deployments for both apps.  
- [ ] Keep contracts in sync with `AGENT_TECH_SPEC.md`; do not change schemas without a spec PR.
