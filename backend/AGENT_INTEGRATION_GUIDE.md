# AGENT_INTEGRATION_GUIDE.md

**Project:** Corporate Energy Platform — UI Integration & Safe Rollout
**Audience:** AI coding agent (Claude/Cursor) + maintainer
**Purpose:** Keep the **existing Vue3/TS demo** working while introducing **live data** with **zero-risk** switches.

> Companion docs: see `PROJECT_ROADMAP.md` (current plan), `AGENT_TECH_SPEC.md` (contracts), and `AGENT_BIZ_SPEC.md` (product logic).
> The agent MUST follow this guide when touching the frontend or deployment pipeline.

**Note:** This project uses a simple structure:
- `frontend/` — old demo (stable, mock-only)
- `frontend-japan/` — new dashboard (dual-mode: MOCK/LIVE)
- Development branch: `feat/japan-dashboard`

---

## 1) Branching & Environments (Zero-Risk)
- Branch strategy:
  - `main` — stable production demo (default **MOCK** mode).
  - `feat/data-adapters` — implementation branch for live data.
  - Use PRs; do **not** push directly to `main`.
- Deployment:
  - Each PR/branch produces a **Preview Deployment** (separate URL).
  - Production URL always points to `main`.
- Protection:
  - Require PR review + CI green before merging to `main`.
  - If a change affects `/public/data`, ensure artifacts validate against schemas.

---

## 2) Data Modes & Feature Flags
We keep two modes side-by-side and switch **without changing UI layout**.

- **MOCK mode (default in `main`)**  
  UI imports or fetches the existing mock datasets (status quo).

- **LIVE mode (behind flag)**  
  UI fetches normalized JSON artifacts or API described in `AGENT_TECH_SPEC.md`.

### 2.1 Mode selection (in order of precedence)
1. URL query: `?mode=live` or `?mode=mock` (session-scoped).  
2. `localStorage.DATA_MODE = "live" | "mock"`.  
3. `.env` at build time: `VITE_DATA_MODE=mock` (default), `live` for preview.  
4. Fallback: **mock**.

### 2.2 Runtime fallback
If LIVE fetch fails (HTTP error, schema mismatch), **gracefully fall back** to MOCK for that widget and surface a tiny non-blocking notice (`meta.warning`).

---

## 3) Data Access Layer (One File)
Create a single adapter/service to decouple UI from sources:

```
/src/services/dataClient.ts
```

Responsibilities:
- Read current **Data Mode**.
- Provide typed functions used by pages/components (no component should know whether data is mock or live):
  - `getTokyoDemand(date): Promise<DemandSeries>`
  - `getKansaiDemand(date): Promise<DemandSeries>`
  - `getReserve(date): Promise<ReserveByArea>`
  - `getJepxSpot(area, date): Promise<PriceSeries>`
  - `getSolar(lat, lon, date): Promise<SolarSeries>`
  - `runSettlement(input): Promise<SettlementResult>`
- Contain **two implementations** under the hood:
  - `mockSource` — current mock loaders.
  - `liveSource` — fetchers of static artifacts `/public/data/...` or `/api/...`.
- Apply **shape adapters** so components always receive the same **UI types**.

---

## 4) File & Path Conventions (Static Artifacts)
If using Static Mode (preferred MVP), normalized JSON files land here:

```
/public/data/jp/tokyo/demand-YYYY-MM-DD.json
/public/data/jp/kansai/demand-YYYY-MM-DD.json
/public/data/jp/system/reserve-YYYY-MM-DD.json
/public/data/jp/jepx/spot-AREA-YYYY-MM-DD.json
/public/data/weather/solar-LAT-LON-YYYY-MM-DD.json
```

Rules:
- All timestamps ISO8601 with `+09:00` offset.  
- Each artifact includes `source.name` and `source.url`.  
- Keep **schema version** via top-level `schema_version` (e.g., `"v1"`).

---

## 5) Mapping (Domain → UI)
> Do **not** change component props. Adapt data in the client.

### 5.1 Demand
- Domain (`/api/jp/{area}/demand`) → UI `DemandSeries` used by charts.  
- Map fields:
  - `series[].ts` → `x`
  - `series[].demand_mw` → `yActual`
  - `series[].forecast_mw` → `yForecast?` (optional line)
- If `forecast_mw` is missing: show only Actual; add non-blocking badge “Forecast unavailable”.

### 5.2 Reserve
- Domain `reserve_margin_pct` + threshold mapping → UI `status` badge (`stable | watch | tight`).
- Tooltip: short explanation from business spec.

### 5.3 JEPX Prices
- Domain `price_yen_per_kwh` → UI series `yPrice`. Ensure 24 rows for day-ahead; if fewer, pad gaps with `null` and annotate warning.

### 5.4 Solar/Weather
- Domain `ghi/dni/dhi/cloudcover` → optional overlay series; default OFF.

### 5.5 Settlement-Lite
- `totals.cost_yen` + `by_hour[]` → cost table and bar chart. If LIVE fails, compute on client using mock price/profile when available.

---

## 6) CI/CD (Safe & Cheap)
- **Static pipeline** (GitHub Actions, cron):
  - Fetch → Normalize → Validate (JSON Schema) → Commit into `/public/data/...`.
  - Avoid full rebuild if not required; allow manual dispatch for backfill.
- **Preview defaults**:
  - In PRs: set `VITE_DATA_MODE=live` to test live artifacts.  
  - On `main`: keep `VITE_DATA_MODE=mock` until the team explicitly flips a flag.
- Add a small `/status.json` with timestamps of latest artifacts for quick health checks.

---

## 7) Testing
- **Unit**: adapters in `dataClient.ts` (happy/edge/empty).  
- **Schema**: validate artifacts against `/pipeline/schemas/*.json`.  
- **E2E smoke**: nightly check that preview URL renders charts/tables in LIVE mode.  
- **Visual** (optional): percy/snapshots for critical charts.

---

## 8) Rollout Playbook (Step-by-Step)
1. Create branch `feat/data-adapters`.
2. Add `/docs/*.md` (tech/biz/this guide).  
3. Implement `dataClient.ts` with **mode switch** (no UI changes).  
4. Add pipeline that writes to `/public/data/...` (cron + manual dispatch).  
5. Open PR → preview deploy defaults to **LIVE** via env (`VITE_DATA_MODE=live`).  
6. QA on preview URL: verify pages, empty states, warnings.  
7. Merge PR; keep `main` in **MOCK** for a few days.  
8. Flip default to **LIVE** by changing env or remote flag.  
9. Monitor; if issues, flip back instantly (no revert needed).

---

## 9) Error Handling & UX
- Never break the page on partial data.  
- Use a compact inline notice for warnings, e.g., “Forecast unavailable today (source)”.  
- Log errors to console with a small prefix `[live-data]` (no secrets).

---

## 10) Agent Tasks Checklist
- [ ] Create `/src/services/dataClient.ts` with dual-mode and adapters.  
- [ ] Wire existing pages to dataClient (no component prop changes).  
- [ ] Implement fetchers for each artifact/API in LIVE mode.  
- [ ] Add JSON Schemas and a validator step to CI.  
- [ ] Implement runtime fallback to mock per widget.  
- [ ] Add `?mode=live|mock` switch and persist to `localStorage`.  
- [ ] Add UI badges for source attribution and optional warnings.  
- [ ] Provide a `README_FOR_AGENT.md` with run instructions (optional).

---

## 11) Non-Goals
- Redesign of charts/tables.  
- Authentication and user management.  
- Persistent DB for UI (handled by backend/pipeline).

---

## 12) Exit Criteria
- Preview deployment renders **LIVE** data with zero UI refactors.  
- Production (`main`) remains stable in **MOCK** until flip.  
- Flip to LIVE can be done by **one env change** and is reversible.


---

## 13) Backend Development Branch (Go/FastAPI) & Safe Frontend Upgrades

**Goal:** develop the backend in isolation, then improve the frontend without risking the current live demo.

### 13.1 Branching Model
- Create a **long‑lived integration branch** for backend work: `feat/backend-core`.
- Open **small PRs** into `feat/backend-core` from focused branches:
  - `feat/adapter-tepco`, `feat/adapter-kansai`, `feat/adapter-jepx`, `feat/adapter-occto`
  - `feat/settlement-lite`, `feat/weather-solar`
  - (optional ML) `feat/forecast-ml` (FastAPI)
- The UI remains unchanged on `main` (MOCK mode). All backend preview testing happens on the **Preview URL** of `feat/backend-core` with `VITE_DATA_MODE=live`.

### 13.2 Promotion to `main`
Merge `feat/backend-core` → `main` only when ALL are true:
1. Static artifacts or API endpoints exist and pass schema validation.
2. `dataClient.ts` supports **dual mode** (MOCK/LIVE) with graceful fallback.
3. LIVE mode is **off by default** on `main` and can be flipped via env/flag.
4. Preview E2E smoke tests are green (charts/tables render; empty states OK).
5. Source attribution + timestamps present in responses/artifacts.
6. No breaking changes to UI props; any data shape differences are adapted inside `dataClient.ts`.

### 13.3 Frontend Improvement Branches (post-backend)
- Create targeted branches for UI features that leverage the backend, e.g.:
  - `feat/ui-forecast-overlay` (solar/weather overlays, confidence bands)
  - `feat/ui-settlement` (cost table + “what‑if” panel)
  - `feat/ui-reserve-badges` (stable/watch/tight)
- Each feature is **behind a flag** and uses the same dual‑mode client.
- UI merges to `main` only after Preview QA; default stays MOCK until flip.

### 13.4 Forecast & Algorithm Changes
- Add a short design note in PRs (inputs, outputs, assumptions).
- Keep **backward‑compatible schemas**; bump `schema_version` when adding fields.
- If LIVE output changes, **update the UI adapter** to preserve component props.
- Include unit tests for math and golden tests for adapter normalization.

### 13.5 Releases & Rollback
- Cut a weekly PR from `feat/backend-core` to `main` (when ready); tag release.
- To rollback, **flip mode to MOCK** (env/flag) — no code revert needed.
- If a hotfix is required, fix on `feat/backend-core` → merge to `main` with the flag still OFF, then re‑enable LIVE after verification.

### 13.6 Environment & Secrets
- Preview deployments (branches/PRs): `VITE_DATA_MODE=live` and any API secrets (if used) scoped to preview only.
- Production (`main`): keep `VITE_DATA_MODE=mock` until an explicit go‑live.
- Never expose secrets in logs or client code; use CI secrets and serverless envs.

> Summary: **Build backend safely in `feat/backend-core`, keep `main` stable**, and gate frontend upgrades behind flags. Promote only when contracts, adapters, and preview checks are solid.