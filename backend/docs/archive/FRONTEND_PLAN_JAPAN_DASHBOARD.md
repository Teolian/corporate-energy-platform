# FRONTEND_PLAN_JAPAN_DASHBOARD.md
**Project:** Corporate Energy Platform — Japanese Energy Dashboard (Vue 3.5 + TS + Tailwind + Chart.js)  
**Goal:** Implement UI features that map directly to *Digital Grid* services while keeping the current demo stable.

> Align with: `JAPAN_DASHBOARD_SPRINT_PLAN.md` (overall plan), `AGENT_TECH_SPEC.md §3` (contracts), and `AGENT_INTEGRATION_GUIDE.md` (modes/flags).

---

## 0) Recap (current)
- Demand charts: **Tokyo** & **Kansai**, date nav, dark mode, peak/avg/forecast-accuracy, dual-mode data client (MOCK/LIVE), responsive ✅

---

## 1) P1 — DGP-aligned UI (System status + Prices + Cost)
### 1.1 Reserve Status (OCCTO) — *closest to DGP*
- **Store**: `src/stores/reserve.ts` → `getReserve(date)` using `dataClient.getReserve(date)`.
- **Components**:
  - `ReserveBadge.vue` — small inline badge on dashboard (“stable / watch / tight”), color-coded.
  - `ReserveGauge.vue` — optional gauge/bar widget under the charts.
- **UX**:
  - Tooltip: short explanation of thresholds.
  - Source attribution (OCCTO) in a compact footer line.
  - Empty/partial day → show neutral badge + `meta.warning`.
- **DoD**:
  - Badge visible for selected date; thresholds correct; handles empty data gracefully.

### 1.2 Spot Prices Overlay (JEPX) — *basis for DGP fixed/spot*
- **Store**: `src/stores/prices.ts` → `getSpotPrice(area, date)` via `dataClient.getJepxSpot(area, date)`.
- **Components**:
  - Extend `DemandChart.vue` to render **price overlay** (second y-axis) with toggle.
  - `SpotPriceChart.vue` (separate route/section) for dedicated price view.
- **UX**:
  - Toggle “Show price overlay”.
  - Legend shows units **¥/kWh**; attribution (JEPX).
  - Gaps → dashed segments + small warning icon with tooltip.
- **DoD**:
  - 24 hourly points aligned to demand; overlay toggle works; gaps handled.

### 1.3 Cost Card (Settlement-lite) — *DGP*
- **Store**: `src/stores/settlement.ts` → `runSettlement(profile, area, date, pvOffsetPct?)`.
- **Components**:
  - `CostCard.vue` — totals: kWh & cost (¥).
  - `CostBreakdownTable.vue` — optional per-hour table when expanded.
- **UX**:
  - Show assumptions: area/date, optional PV% reduction.
  - CSV export button (see §3.3).
- **DoD**:
  - Totals equal backend result for test payload; shows assumptions + attribution.

---

## 2) P2 — Market Context UX
### 2.1 Comparison View (Tokyo vs Kansai)
- **Route**: `/comparison`
- **Component**: `ComparisonChart.vue` (overlay 2 regions, shared x-axis, option shared or dual y-axis)
- **DoD**: renders live data for both areas; date nav synced.

### 2.2 Dedicated Price View
- **Route**: `/prices`
- **Component**: `SpotPriceChart.vue` (line chart of ¥/kWh; optional correlation with demand)
- **DoD**: shows price-only and price+overlay modes.

### 2.3 Reserve Visualization
- Add `ReserveGauge.vue` to dashboard; color palette: green/yellow/red with Tailwind tokens.
- **DoD**: correct color mapping; tooltip text from config.

### 2.4 Exports (CSV/PNG)
- **Composable**: `src/composables/useExport.ts`
- Buttons in dashboard to export current demand (CSV) and chart image (PNG).
- **DoD**: files download with sensible filenames and ISO dates.

---

## 3) Data Client & Types
- Extend `src/services/dataClient.ts` with:
  - `getReserve(date)`, `getJepxSpot(area, date)`, `runSettlement(input)`.
- **Shape adapters (no prop changes in components):**
  - **Demand**: `{ series[].ts, demand_mw, forecast_mw? }` → `{ x, yActual, yForecast? }`
  - **Reserve**: `{ reserve_margin_pct, status }` → badge props `{ status, pct }`
  - **Prices**: `{ price_yen_per_kwh[].{ts, price} }` → `{ x, yPrice }`
  - **Settlement**: `{ totals, by_hour[] }` → card + table DTOs
- Always surface `source.name` + `source.url` in the UI footer for each widget.

---

## 4) Feature Flags & Modes
- Keep **MOCK/LIVE** mode selection from `AGENT_INTEGRATION_GUIDE.md`.
- Add **feature flags** (env or runtime config) per widget:
  - `VITE_FEAT_RESERVE`, `VITE_FEAT_JEPX`, `VITE_FEAT_SETTLEMENT`, `VITE_FEAT_COMPARISON`
- Default: flags **ON** in preview, **OFF** on `main` until validated.

---

## 5) Routes & Navigation
- `/` — Dashboard (Demand + Reserve badge + Price overlay + Cost card)
- `/prices` — Spot prices focus
- `/comparison` — Tokyo vs Kansai overlay
- Future: `/mix`, `/certs`, `/rebridge`, `/battery`

---

## 6) UX States & A11y
- Loading skeletons per widget; empty/partial with neutral messages.
- Keyboard focus states; aria-labels on toggles and export buttons.
- Color contrast: respect Tailwind dark mode and color tokens.

---

## 7) i18n (EN/JA)
- `src/i18n/messages.ts` — minimal dictionary:
  - `reserve.status.stable|watch|tight`
  - `labels.price_yen_per_kwh`, `labels.cost_yen`, `labels.kwh`
  - `tooltips.gaps_detected`, `tooltips.reserve_explain`
- Start with EN; add JA keys stub to show readiness.

---

## 8) Testing
- **Unit**: adapters in `dataClient`, stores (`reserve`, `prices`, `settlement`).
- **Integration**: charts render with live artefacts (mocked fetch).
- **E2E (Playwright)**:
  - Load date → demand renders
  - Toggle price overlay
  - Reserve badge visible
  - Export CSV/PNG works
  - Switch MOCK/LIVE — no crash

---

## 9) Acceptance Criteria (Front P1/P2)
- Dashboard: demand + reserve badge + price overlay + cost card visible in LIVE; warnings on partial data.
- `/prices` and `/comparison` routes functional; exports ok.
- No prop changes to existing components; all new data via `dataClient` and new stores.

---

## 10) Agent Checklist (Front)
- [ ] Implement `stores/reserve.ts`, `stores/prices.ts`, `stores/settlement.ts`
- [ ] Extend `dataClient` with `getReserve`, `getJepxSpot`, `runSettlement`
- [ ] Add `ReserveBadge.vue`, `ReserveGauge.vue`, `SpotPriceChart.vue`, `CostCard.vue`, `CostBreakdownTable.vue`, `ComparisonChart.vue`
- [ ] Wire routes `/prices`, `/comparison`
- [ ] Add `useExport.ts` and buttons on dashboard
- [ ] Add i18n scaffolding (EN/JA), source attribution footers
- [ ] Add feature flags per widget; keep OFF on main until validated
