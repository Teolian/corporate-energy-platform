# AGENT_EXECUTION_BRIEF.md
**Project:** Corporate Energy Platform — Japan Dashboard (Backend + Frontend)  
**Audience:** AI coding agent (Claude/Cursor)  
**Purpose:** Single entry-point brief. Follow this doc top‑to‑bottom. For details, see linked specs in /docs.

---

## 0) What to open first
Read these files (in order). Do **not** diverge from their contracts.

1. `AGENT_TECH_SPEC.md` — API contracts (§3) and guardrails (canonical).  
2. `JAPAN_DASHBOARD_SPRINT_PLAN.md` — sequenced backlog aligned to Digital Grid services.  
3. `FRONTEND_PLAN_JAPAN_DASHBOARD.md` — exact UI features/stores/routes.  
4. `FRONTEND_GREENFIELD_ALIGNMENT.md` — instructions for the new greenfield FE (`frontend-live`).

> If anything conflicts, **AGENT_TECH_SPEC.md §3** wins. Update other docs via PR.

---

## 1) Branching & Apps
- Keep the original demo as **`frontend-mock`** (MOCK only).  
- Treat the new simple frontend as **`frontend-live`** (LIVE + flags).  
- Backend lives in `feat/backend-core` (or as /apps/backend if monorepo).  
- `main` stays MOCK by default; preview deployments use LIVE.

---

## 2) P1 Deliverables (1–2 weeks) — *closest to Digital Grid “DGP”*
**Backend (Mode A: static pipeline first)**  
- **OCCTO Reserve adapter** → artefact `public/data/jp/system/reserve-YYYY-MM-DD.json`.  
- **JEPX Spot prices adapter** → `public/data/jp/jepx/spot-AREA-YYYY-MM-DD.json`.  
- **Settlement‑lite** endpoint or CLI → kWh × ¥/kWh (deterministic rounding).

**Frontend (in `frontend-live`)**  
- **ReserveBadge/Gauge** on dashboard (thresholds: >=8 stable, 5–8 watch, <5 tight).  
- **Price overlay** on demand charts + route `/prices`.  
- **CostCard** with totals (+ optional hourly table) and **CSV export**.

**Acceptance for P1**  
- Dashboard shows **demand + reserve badge + price overlay + cost card** in LIVE.  
- All artefacts include `source.name` + `source.url`.  
- Partial/missing data → non‑blocking warnings.  
- Golden tests for adapters; schema validation in CI.

---

## 3) P2 Deliverables (UX context)
- Route `/comparison` (Tokyo vs Kansai overlay).  
- Exports: CSV of day demand, PNG of charts.  
- Dedicated `SpotPriceChart` page with optional demand correlation.

---

## 4) Data Client (greenfield FE)
Implement one adapter layer (`src/services/dataClient.ts`) that exposes:
```ts
getDemand(area: 'tokyo' | 'kansai', date: ISODate): Promise<DemandSeries>
getReserve(date: ISODate): Promise<ReserveByArea>
getJepxSpot(area: 'tokyo' | 'kansai', date: ISODate): Promise<PriceSeries>
runSettlement(input: SettlementInput): Promise<SettlementResult>
```
- Keep component props unchanged; map backend shapes inside this layer.
- Always surface attribution and show `meta.warning` when present.

---

## 5) Flags, Modes, ENV
- `VITE_DATA_MODE=mock|live` (preview=live; main=mock).  
- Feature flags (ON in preview, OFF on main until validated):  
  `VITE_FEAT_RESERVE`, `VITE_FEAT_JEPX`, `VITE_FEAT_SETTLEMENT`, `VITE_FEAT_COMPARISON`.

---

## 6) Definition of Done (per PR)
- Meets contracts in `AGENT_TECH_SPEC.md §3`.  
- Tests: unit + golden (adapters), integration (stores), E2E smoke (render, overlay, export).  
- Artefacts validated by JSON Schema in CI.  
- No breaking UI prop changes; all adaptations in `dataClient`.  
- Source attribution rendered in the UI.

---

## 7) Hand‑off notes
- Prefer **Mode A (static)** for zero runtime cost; later add Mode B (serverless /api).  
- Be polite to sources (timeouts, backoff).  
- If JEPX data is unavailable for a date, fall back to MOCK for that widget only.

---

## 8) Quick links (expected paths in /docs)
- `/docs/AGENT_TECH_SPEC.md`  
- `/docs/JAPAN_DASHBOARD_SPRINT_PLAN.md`  
- `/docs/FRONTEND_PLAN_JAPAN_DASHBOARD.md`  
- `/docs/FRONTEND_GREENFIELD_ALIGNMENT.md`
