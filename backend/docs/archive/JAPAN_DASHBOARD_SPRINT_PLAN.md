# JAPAN_DASHBOARD_SPRINT_PLAN.md
**Project:** Corporate Energy Platform — Japan Dashboard (Backend + Frontend)  
**Scope:** Align next work to *Digital Grid* services and maximize demo value.  
**Sources of truth:** `AGENT_TECH_SPEC.md` (§3 contracts), `AGENT_BIZ_SPEC.md`, `AGENT_INTEGRATION_GUIDE.md`, `AGENT_TOPOLOGY.md`.

> This plan assumes the current branches:  
> - `main` — stable demo  
> - `feat/frontend-japan` — active branch (backend adapters + Japan dashboard)  
> If you use two frontends per `AGENT_TOPOLOGY.md`, apply UI steps to `apps/frontend-live` and keep `apps/frontend-mock` unchanged.

---

## 0) Current Status (from agent report)
- **Backend (Go 1.23, Mode A: Static Pipeline):**
  - TEPCO **Tokyo** demand adapter (CSV → JSON), golden tests ✅
  - **Kansai** demand adapter, golden tests ✅
  - Time utilities (Asia/Tokyo, ISO8601 +09:00) ✅
  - HTTP fetch wrapper (timeouts, UA — retries TBD) ✅
  - CLI: emits JSON into `/public/data/jp/{area}/demand-YYYY-MM-DD.json` ✅
- **Frontend (Vue 3.5 + Tailwind + Chart.js):**
  - Japan Dashboard MVP (Tokyo & Kansai charts, dark mode, date nav) ✅
  - Metrics: Peak / Average / Forecast Accuracy ✅
  - **Dual-mode** data client (MOCK/LIVE with graceful fallback) ✅

**Dev URLs:** local dev `http://localhost:5174/`

---

## 1) Priority Plan (value-aligned to *Digital Grid* services)

### P1 — **System status + prices + cost** (1–2 weeks)
1. **OCCTO Reserve Margin adapter** → *closest to* **DGP** (context for risk/price)
   - Backend: `GET /api/jp/system/reserve?date=YYYY-MM-DD` (or static JSON artefact)  
   - Mapping: `>=8% → stable`, `5–8% → watch`, `<5% → tight`
   - Frontend: status badge + tooltip on dashboard
   - **Acceptance:** daily JSON valid by schema; empty days handled
2. **JEPX Spot Prices adapter** → *closest to* **DGP** (fixed/spot mix basis)
   - Backend: `GET /api/jp/jepx/spot?date=&area=` (24 hourly ¥/kWh)  
   - Frontend: price overlay on demand chart
   - **Acceptance:** 24 slots, attribution, graceful gaps
3. **Settlement‑lite endpoint** → *closest to* **DGP**
   - Backend: `POST /api/settlements/run` (profile kWh × JEPX price); deterministic rounding
   - Frontend: “Cost this period” card + CSV export
   - **Acceptance:** golden test equals manual calc

### P2 — **Market context UX** (1 week)
4. **Reserve visualization widget** (gauge/bar) → *DGP*
5. **Spot price chart** (with demand correlation) → *DGP*
6. **Comparison view** Tokyo vs Kansai overlay (context) → *DGP*
7. **Exports** (CSV/PNG) for charts/tables (reporting)

### P3 — **DG services adjacent features** (2–3 weeks, optional)
8. **Portfolio Mix (fixed/spot/RE share)** → *DGP*
   - Compute mix cost; simple scenario compare (Base vs Alt)
9. **Certificates calendar** (`/api/certs/calendar`) → *Econohashi*
   - Read‑only schedule & explainer
10. **PPA digest** (`/api/rebridge/digest`) → *RE Bridge*
    - Aggregate metrics from public press releases (read‑only)
11. **Battery arbitrage sim** (`/api/battery/sim`) → *Battery Aggregation*
    - Simple charge/discharge model vs price

---

## 2) Detailed Tasks & DoD

### 2.1 OCCTO Reserve Margin (Backend P1.1)
- **Tasks**
  - `backend/internal/adapters/occto.go` (HTML/CSV parse → normalized JSON)
  - `public/data/jp/system/reserve-YYYY-MM-DD.json` (Mode A)
  - Golden tests + schema validation
  - CLI flag: `--provider occto --date YYYY-MM-DD`
- **DoD**
  - JSON matches spec; thresholds map to `stable|watch|tight`
  - Missing data → `meta.warning` + 200 status for artefact

### 2.2 JEPX Spot Prices (Backend P1.2)
- **Tasks**
  - `backend/internal/adapters/jepx.go` (history download → normalized 24h ¥/kWh)
  - CLI: `--provider jepx --area tokyo --date YYYY-MM-DD`
  - Artefact: `public/data/jp/jepx/spot-AREA-YYYY-MM-DD.json`
  - Golden tests; handle registration/availability gracefully
- **DoD**
  - 24 timestamps in Asia/Tokyo; attribution; warnings on gaps

### 2.3 Settlement‑lite (Backend P1.3 + FE)
- **Tasks (BE)**
  - `POST /api/settlements/run` (or CLI script that writes artefact for FE MVP)
  - Inputs: profile (kWh), prices (area/date), `pv_offset_pct?`
  - Math: `Σ kWh_slot × price_slot × (1 - pv%)`; rounding: 0.1 JPY, 0.1 kWh
  - Tests: golden + edge (missing slots, misaligned TZ)
- **Tasks (FE)**
  - Card with totals, table by hour; CSV export
- **DoD**
  - Deterministic totals; error messages non‑blocking

### 2.4 Reserve & JEPX visualizations (FE P2)
- **Tasks**
  - `stores/reserve.ts`, `stores/prices.ts`
  - `ReserveGauge.vue`, `SpotPriceChart.vue`
  - Overlay toggle with demand
- **DoD**
  - Works in LIVE; falls back to MOCK per widget

### 2.5 Comparison & Exports (FE P2)
- **Tasks**
  - `ComparisonChart.vue`, route `/comparison`
  - `useExport.ts` (CSV/PNG)
- **DoD**
  - Tokyo & Kansai overlay verified; downloads functional

### 2.6 DG‑adjacent features (P3 optional)
- Portfolio Mix → *DGP*; Certificates Calendar → *Econohashi*; PPA Digest → *RE Bridge*; Battery Sim → *Battery Aggregation*.

---

## 3) Branching & Safety
- Continue work in `feat/frontend-japan` (or `feat/backend-core` + UI branches if using monorepo topology).
- `main` stays MOCK by default. Preview deployments test LIVE.
- Promotion to `main` only if:
  1) contracts unchanged & validated;  
  2) preview renders LIVE without regressions;  
  3) graceful fallback to MOCK works;  
  4) attribution present in responses/artefacts.

---

## 4) Checklists (Agent)

**Backend P1**
- [ ] OCCTO adapter with golden tests and CLI
- [ ] JEPX adapter with golden tests and CLI
- [ ] Real HTTP fetching + polite rate limiting + backoff
- [ ] Structured logs (`source`, `duration_ms`, `status`)
- [ ] Artefacts emitted to `/public/data/...` (Mode A)

**Frontend P2**
- [ ] Reserve badge + widget
- [ ] Spot price overlay + chart
- [ ] Comparison overlay view
- [ ] Exports (CSV/PNG)
- [ ] Empty/partial states covered

**Quality**
- [ ] Schema validation step in CI
- [ ] Unit + integration + golden tests
- [ ] E2E smoke: load date → render charts → export

---

## 5) Notes & Risks
- **JEPX availability** may require registration; implement fallback to MOCK for that date/area.
- **OCCTO site** is JP‑only and can change layout; keep selectors resilient.
- Keep **API contracts** *exactly* per `AGENT_TECH_SPEC.md §3`.
- Always include `source.name` and `source.url` in artefacts.

---

## 6) Success Criteria for P1
- Dashboard shows: Tokyo/Kansai demand + **reserve status** + **price overlay**, with **cost card** for a selected day.  
- Artefacts refreshable by CLI and (optionally) hourly cron.  
- No UI refactors needed; MOCK/LIVE switch safe.
